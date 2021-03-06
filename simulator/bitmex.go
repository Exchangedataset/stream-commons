package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitmexChannelInfo is the channel name for info channel on Bitmex
const BitmexChannelInfo = "info"

// BitmexChannelError is the channel name for error channel on Bitmex
const BitmexChannelError = "error"

// BitmexSideSell is the string name for sell side
const BitmexSideSell = "Sell"

// BitmexSideBuy is the string name for buy side
const BitmexSideBuy = "Buy"

type bitmexOrderBookL2Element struct {
	price float64
	size  uint64
}

type bitmexSimulator struct {
	filterChannel map[string]bool
	subscribed    map[string]bool
	// map[symbol]map[side]map[id]
	orderBooks map[string]map[string]map[int64]bitmexOrderBookL2Element
}

func (s *bitmexSimulator) ProcessStart(line []byte) error {
	return nil
}

func (s *bitmexSimulator) ProcessSend(line []byte) (channel string, err error) {
	// this should not be called
	return streamcommons.ChannelUnknown, nil
}

func (s *bitmexSimulator) processData(action string, dataSlice []jsonstructs.BitmexOrderBookL2DataElement) error {
	for _, data := range dataSlice {
		if action == "partial" || action == "insert" {
			sides, ok := s.orderBooks[data.Symbol]
			if !ok {
				// symbol is not yet pushed, create new map
				sides = make(map[string]map[int64]bitmexOrderBookL2Element)
				s.orderBooks[data.Symbol] = sides
			}
			ids, ok := sides[data.Side]
			if !ok {
				// map for this side is not prepared, create new one
				ids = make(map[int64]bitmexOrderBookL2Element)
				sides[data.Side] = ids
			}
			// set new element
			ids[data.ID] = bitmexOrderBookL2Element{price: data.Price, size: data.Size}

			// check logical error
			if data.Side == BitmexSideBuy {
				sellIDs, ok := sides[BitmexSideSell]
				// actually, it does not have to check if map exists, because range of nil map is no-op
				if ok {
					dueToRemove := make([]int64, 0, 5)
					for anoID, anoElem := range sellIDs {
						if anoElem.price < data.Price {
							// original order is buy, this order is sell but has lower price than original, weird
							dueToRemove = append(dueToRemove, anoID)
							fmt.Println("sell logical error:", data.Price, data.Size, anoElem.price, anoElem.size)
						}
					}
					for _, anoID := range dueToRemove {
						delete(sellIDs, anoID)
					}
				}
			} else {
				buyIDs, ok := sides[BitmexSideBuy]
				if ok {
					dueToRemove := make([]int64, 0, 5)
					for anoID, anoElem := range buyIDs {
						if anoElem.price > data.Price {
							// original order is sell, this order is buy but has higher price than original, weird
							dueToRemove = append(dueToRemove, anoID)
							fmt.Println("buy logical error:", data.Price, data.Size, anoElem.price, anoElem.size)
						}
					}
					for _, anoID := range dueToRemove {
						delete(buyIDs, anoID)
					}
				}
			}
		} else if action == "update" {
			// update for element, it can expect element to be there already,
			// so map is already prepared
			// map returns-by-value it needs to replace value after you updated it
			elem, ok := s.orderBooks[data.Symbol][data.Side][data.ID]
			if ok {
				elem.size = data.Size
				s.orderBooks[data.Symbol][data.Side][data.ID] = elem
			} else {
				fmt.Println("order id not found")
			}
		} else if action == "delete" {
			// delete element
			delete(s.orderBooks[data.Symbol][data.Side], data.ID)
		} else {
			return fmt.Errorf("unknown action type '%s'", action)
		}
	}
	return nil
}

func (s *bitmexSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown

	// check if this message is a response to subscribe
	subscribe := jsonstructs.BitmexSubscribe{}
	err = json.Unmarshal(line, &subscribe)
	if err != nil {
		return
	}
	if subscribe.Success {
		// this is subscribe message
		channel = subscribe.Subscribe
		// check if this channel should be tracked
		if s.filterChannel == nil {
			s.subscribed[channel] = true
		} else {
			// filtering is enabled
			_, ok := s.filterChannel[channel]
			if ok {
				s.subscribed[channel] = true
			}
		}

		return
	}

	decoded := new(jsonstructs.BitmexRoot)
	err = json.Unmarshal(line, decoded)
	if err != nil {
		return
	}
	if decoded.Info != nil {
		channel = BitmexChannelInfo
		return
	}
	if decoded.Error != nil {
		channel = BitmexChannelError
		return
	}
	channel = decoded.Table

	if channel == "orderBookL2" {
		dataSlice := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
		err = json.Unmarshal(decoded.Data, &dataSlice)
		if err != nil {
			return
		}

		err = s.processData(decoded.Action, dataSlice)

		return
	}
	// ignore other channels
	return
}

func (s *bitmexSimulator) ProcessMessageChannelKnown(channel string, line []byte) error {
	wsChannel, serr := s.ProcessMessageWebSocket(line)
	if serr != nil {
		return serr
	}
	if wsChannel != channel {
		return fmt.Errorf("channel differs: %v, expected: %v", wsChannel, channel)
	}
	return nil
}

func (s *bitmexSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == streamcommons.StateChannelSubscribed {
		// add to subscribed
		subscribed := jsonstructs.BitmexStateSubscribed{}
		err = json.Unmarshal(line, &subscribed)
		if err != nil {
			return
		}
		for _, subscrCh := range subscribed {
			// record subscribed channel only if it is in target channel
			if _, ok := s.filterChannel[subscrCh]; ok {
				s.subscribed[subscrCh] = true
			}
		}
		return
	}

	if s.filterChannel != nil {
		if _, ok := s.filterChannel[channel]; !ok {
			return
		}
	}

	if channel == "orderBookL2" {
		// process orderbook
		decoded := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
		err = json.Unmarshal(line, &decoded)
		if err != nil {
			return
		}
		return s.processData("partial", decoded)
	}

	return
}

func sortBitmexSubscribe(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexOrderbooks(m map[string]map[string]map[int64]bitmexOrderBookL2Element) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexSides(m map[string]map[int64]bitmexOrderBookL2Element) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexID(m map[int64]bitmexOrderBookL2Element) []int64 {
	keys := make([]int64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (s *bitmexSimulator) orderBookL2DataElements() []jsonstructs.BitmexOrderBookL2DataElement {
	// reconstruct raw-like json format bitmex sends
	data := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
	for _, symbol := range sortBitmexOrderbooks(s.orderBooks) {
		sides := s.orderBooks[symbol]
		for _, side := range sortBitmexSides(sides) {
			ids := sides[side]
			for _, id := range sortBitmexID(ids) {
				elem := ids[id]
				data = append(data, jsonstructs.BitmexOrderBookL2DataElement{
					ID:     id,
					Price:  elem.price,
					Side:   side,
					Size:   elem.size,
					Symbol: symbol,
				})
			}
		}
	}
	return data
}

func (s *bitmexSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	if s.filterChannel != nil {
		// If channel filtering is enabled, this should not be called
		err = errors.New("channel filter is enabled")
		return
	}
	snapshots = make([]Snapshot, 0, 5)

	// list subscribed channels
	subList := make([]string, len(s.subscribed))
	for i, channel := range sortBitmexSubscribe(s.subscribed) {
		subList[i] = channel
	}
	var subListMarshaled []byte
	subListMarshaled, err = json.Marshal(subList)
	if err != nil {
		return nil, fmt.Errorf("error on json marshal: %s", err.Error())
	}
	snapshots = append(snapshots, Snapshot{Channel: streamcommons.StateChannelSubscribed, Snapshot: subListMarshaled})

	data := s.orderBookL2DataElements()
	var orderBookL2ElementsMarshaled []byte
	orderBookL2ElementsMarshaled, err = json.Marshal(data)
	if err != nil {
		return
	}
	snapshots = append(snapshots, Snapshot{Channel: "orderBookL2", Snapshot: orderBookL2ElementsMarshaled})

	return
}

func (s *bitmexSimulator) TakeSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// subscribe message
	for _, channel := range sortBitmexSubscribe(s.subscribed) {
		subscr := jsonstructs.BitmexSubscribe{}
		subscr.Initialize()
		subscr.Subscribe = channel

		var subscribeMarshaled []byte
		subscribeMarshaled, err = json.Marshal(subscr)
		if err != nil {
			return
		}

		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: subscribeMarshaled})
	}

	_, ok := s.subscribed["orderBookL2"]
	if ok {
		// reconstruct raw-like json format bitmex sends
		data := s.orderBookL2DataElements()
		var dataMarshaled []byte
		dataMarshaled, err = json.Marshal(data)
		if err != nil {
			return
		}
		root := new(jsonstructs.BitmexRoot)
		root.Table = "orderBookL2"
		// partial means full orderbook snapshot
		root.Action = "partial"
		root.Data = json.RawMessage(dataMarshaled)
		var rootMarshaled []byte
		rootMarshaled, err = json.Marshal(root)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: "orderBookL2", Snapshot: rootMarshaled})
	}

	return
}

func newBitmexSimulator(filterChannels []string) Simulator {
	gen := bitmexSimulator{}
	if filterChannels != nil {
		gen.filterChannel = make(map[string]bool)
		for _, ch := range filterChannels {
			// this value will be ignored, filter will be applied to channel that has value in this map
			// value itself does not matter
			gen.filterChannel[ch] = true
		}
	}
	gen.subscribed = make(map[string]bool, 0)
	gen.orderBooks = make(map[string]map[string]map[int64]bitmexOrderBookL2Element)
	return &gen
}
