package formatter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/formatter/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// binanceFormatter is json formatter for Binance.
type binanceFormatter struct{}

// FormatStart formats start line (URL) and returns the array of known subscribed channel in case the server won't
// tell the client what channels are successfully subscribed.
func (f *binanceFormatter) FormatStart(urlStr string) (formatted []Result, err error) {
	u, serr := url.Parse(string(urlStr))
	if serr != nil {
		return nil, fmt.Errorf("FormatStart: %v", serr)
	}
	q := u.Query()
	streams := q.Get("streams")
	channels := strings.Split(streams, "/")
	formatted = make([]Result, len(channels))
	for i, ch := range channels {
		_, stream, serr := streamcommons.BinanceDecomposeChannel(ch)
		if serr != nil {
			err = fmt.Errorf("FormatStart: %v", serr)
			return
		}
		switch stream {
		case streamcommons.BinanceStreamDepth:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceDepth,
			}
		case streamcommons.BinanceStreamTrade:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceTrade,
			}
		case streamcommons.BinanceStreamTicker:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceTicker,
			}
		case streamcommons.BinanceStreamRESTDepth:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceRestDepth,
			}
		default:
			err = fmt.Errorf("FormatStart: channel not supported: %s", ch)
			return
		}
	}
	return formatted, nil
}

func (f *binanceFormatter) formatTicket(channel string, line []byte) (formatted []Result, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("formatTicket: %v", err)
		}
	}()
	root := new(jsonstructs.BinanceReponseRoot)
	serr := json.Unmarshal(line, root)
	if serr != nil {
		err = fmt.Errorf("BinanceReponseRoot: %v", serr)
		return
	}
	ticker := new(jsonstructs.BinanceTickerStream)
	serr = json.Unmarshal(root.Data, ticker)
	if serr != nil {
		err = fmt.Errorf("BinanceTickerStream: %v", serr)
		return
	}
	ft := new(jsondef.BinanceTicker)
	ft.EventTime = strconv.FormatInt(ticker.EventTime*int64(time.Millisecond), 10)
	ft.Symbol = ticker.Symbol
	priceChange, serr := strconv.ParseFloat(ticker.PriceChange, 64)
	if serr != nil {
		err = fmt.Errorf("priceChange: %v", serr)
		return
	}
	ft.PriceChange = priceChange
	priceChangePercent, serr := strconv.ParseFloat(ticker.PriceChanePercent, 64)
	if serr != nil {
		err = fmt.Errorf("priceChangePercent: %v", serr)
		return
	}
	ft.PriceChanePercent = priceChangePercent
	weightedAveragePrice, serr := strconv.ParseFloat(ticker.WeightedAveragePrice, 64)
	if serr != nil {
		err = fmt.Errorf("weightedAveragePrice: %v", serr)
		return
	}
	ft.WeightedAveragePrice = weightedAveragePrice
	firstTradePrice, serr := strconv.ParseFloat(ticker.FirstTradePrice, 64)
	if serr != nil {
		err = fmt.Errorf("firstTradePrice: %v", serr)
		return
	}
	ft.FirstTradePrice = firstTradePrice
	lastPrice, serr := strconv.ParseFloat(ticker.LastPrice, 64)
	if serr != nil {
		err = fmt.Errorf("lastPrice: %v", serr)
		return
	}
	ft.LastPrice = lastPrice
	lastQuantity, serr := strconv.ParseFloat(ticker.LastQuantity, 64)
	if serr != nil {
		err = fmt.Errorf("lastQuantity: %v", serr)
		return
	}
	ft.LastQuantity = lastQuantity
	bestBidPrice, serr := strconv.ParseFloat(ticker.BestBidPrice, 64)
	if serr != nil {
		err = fmt.Errorf("bestBidPrice: %v", serr)
		return
	}
	ft.BestBidPrice = bestBidPrice
	bestBidQuantity, serr := strconv.ParseFloat(ticker.BestBidQuantity, 64)
	if serr != nil {
		err = fmt.Errorf("bestBidQuantity: %v", serr)
		return
	}
	ft.BestBidQuantity = bestBidQuantity
	bestAskPrice, serr := strconv.ParseFloat(ticker.BestAskPrice, 64)
	if serr != nil {
		err = fmt.Errorf("bestAskPrice: %v", serr)
		return
	}
	ft.BestAskPrice = bestAskPrice
	bestAskQuantity, serr := strconv.ParseFloat(ticker.BestAskQuantity, 64)
	if serr != nil {
		err = fmt.Errorf("bestAskQuantity: %v", serr)
		return
	}
	ft.BestAskQuantity = bestAskQuantity
	openPrice, serr := strconv.ParseFloat(ticker.OpenPrice, 64)
	if serr != nil {
		err = fmt.Errorf("openPrice: %v", serr)
		return
	}
	ft.OpenPrice = openPrice
	highPrice, serr := strconv.ParseFloat(ticker.HighPrice, 64)
	if serr != nil {
		err = fmt.Errorf("highPrice: %v", serr)
		return
	}
	ft.HighPrice = highPrice
	lowPrice, serr := strconv.ParseFloat(ticker.LowPrice, 64)
	if serr != nil {
		err = fmt.Errorf("lowPrice: %v", serr)
		return
	}
	ft.LowPrice = lowPrice
	totalTradedBaseAssetVolume, serr := strconv.ParseFloat(ticker.TotalTradedBaseAssetVolume, 64)
	if serr != nil {
		err = fmt.Errorf("totalTradedBaseAssetVolume: %v", serr)
		return
	}
	ft.TotalTradedBaseAssetVolume = totalTradedBaseAssetVolume
	totalTradedQuoteAssetVolume, serr := strconv.ParseFloat(ticker.TotalTradedQuoteAssetVolume, 64)
	if serr != nil {
		err = fmt.Errorf("totalTradedQuoteAssetVolume: %v", serr)
		return
	}
	ft.TotalTradedQuoteAssetVolume = totalTradedQuoteAssetVolume
	ft.StatisticsOpenTime = strconv.FormatInt(ticker.StatisticsOpenTime*int64(time.Millisecond), 10)
	ft.StatisticsCloseTime = strconv.FormatInt(ticker.StatisticsCloseTime*int64(time.Millisecond), 10)
	ft.FirstTradeID = ticker.FirstTradeID
	ft.LastTradeID = ticker.LastTradeID
	ft.TotalNumberOfTrades = ticker.TotalNumberOfTrades
	formatted = make([]Result, 1)
	mft, serr := json.Marshal(ft)
	if serr != nil {
		err = fmt.Errorf("BinanceTicker: %v", serr)
		return
	}
	formatted[0] = Result{
		Channel: channel,
		Message: mft,
	}
	return
}

func (f *binanceFormatter) formatTrade(channel string, line []byte) (formatted []Result, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("formatTrade: %v", err)
		}
	}()
	root := new(jsonstructs.BinanceReponseRoot)
	serr := json.Unmarshal(line, root)
	if serr != nil {
		err = fmt.Errorf("BinanceReponseRoot: %v", serr)
		return
	}
	trade := new(jsonstructs.BinanceTrade)
	serr = json.Unmarshal(root.Data, trade)
	if serr != nil {
		err = fmt.Errorf("BinanceTrade: %v", serr)
		return
	}
	ft := new(jsondef.BinanceTrade)
	ft.EventTime = strconv.FormatInt(trade.EventTime*int64(time.Millisecond), 10)
	ft.Timestamp = strconv.FormatInt(trade.TradeTime*int64(time.Millisecond), 10)
	ft.Symbol = trade.Symbol
	ft.Price, serr = strconv.ParseFloat(trade.Price, 64)
	if serr != nil {
		err = fmt.Errorf("price: %v", serr)
		return
	}
	ft.Size, serr = strconv.ParseFloat(trade.Quantity, 64)
	if serr != nil {
		err = fmt.Errorf("quantity: %v", serr)
		return
	}
	if trade.IsBuyerMarketMaker {
		// Buyer is the maker = seller is the taker
		ft.Side = streamcommons.CommonFormatSell
	} else {
		ft.Side = streamcommons.CommonFormatBuy
	}
	ft.SellterOrderID = trade.SellerOrderID
	ft.BuyerOrderID = trade.BuyerOrderID
	ft.TradeID = trade.TradeID
	formatted = make([]Result, 1)
	mft, serr := json.Marshal(ft)
	if serr != nil {
		err = fmt.Errorf("BinanceTrade: %v", serr)
		return
	}
	formatted[0] = Result{
		Channel: channel,
		Message: mft,
	}
	return
}

func (f *binanceFormatter) formatRESTDepth(channel string, line []byte, symbol string) (formatted []Result, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("formatRESTDepth: %v", err)
		}
	}()
	symbolCap := strings.ToUpper(symbol)
	depth := new(jsonstructs.BinanceDepthREST)
	serr := json.Unmarshal(line, depth)
	if serr != nil {
		err = fmt.Errorf("BinanceDepthREST: %v", serr)
		return
	}
	formatted = make([]Result, len(depth.Asks)+len(depth.Bids))
	i := 0
	for _, order := range depth.Asks {
		fo := new(jsondef.BinanceRestDepth)
		fo.Symbol = symbolCap
		fo.Price, serr = strconv.ParseFloat(order[0], 64)
		if serr != nil {
			err = fmt.Errorf("ask price: %v", serr)
			return
		}
		fo.Side = streamcommons.CommonFormatSell
		size, serr := strconv.ParseFloat(order[1], 64)
		if serr != nil {
			err = fmt.Errorf("ask size: %v", serr)
			return
		}
		fo.Size = size
		mfo, serr := json.Marshal(fo)
		if serr != nil {
			err = fmt.Errorf("ask BinanceDepth: %v", serr)
			return
		}
		formatted[i] = Result{
			Channel: channel,
			Message: mfo,
		}
		i++
	}
	for _, order := range depth.Bids {
		fo := new(jsondef.BinanceRestDepth)
		fo.Symbol = symbolCap
		fo.Price, serr = strconv.ParseFloat(order[0], 64)
		if serr != nil {
			err = fmt.Errorf("bid price: %v", serr)
			return
		}
		fo.Side = streamcommons.CommonFormatBuy
		fo.Size, serr = strconv.ParseFloat(order[1], 64)
		if serr != nil {
			err = fmt.Errorf("bid size: %v", serr)
			return
		}
		mfo, serr := json.Marshal(fo)
		if serr != nil {
			err = fmt.Errorf("bid BinanceDepth: %v", serr)
			return
		}
		formatted[i] = Result{
			Channel: channel,
			Message: mfo,
		}
		i++
	}
	return
}

func (f *binanceFormatter) formatDepth(channel string, line []byte) (formatted []Result, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("formatDepth: %v", err)
		}
	}()
	root := new(jsonstructs.BinanceReponseRoot)
	serr := json.Unmarshal(line, root)
	if serr != nil {
		err = fmt.Errorf("BinanceReponseRoot: %v", serr)
		return
	}
	depth := new(jsonstructs.BinanceDepthStream)
	serr = json.Unmarshal(root.Data, depth)
	if serr != nil {
		err = fmt.Errorf("BinanceDepthStream: %v", serr)
		return
	}
	eventTime := strconv.FormatInt(depth.EventTime*int64(time.Millisecond), 10)
	formatted = make([]Result, len(depth.Asks)+len(depth.Bids))
	i := 0
	for _, order := range depth.Asks {
		fo := new(jsondef.BinanceDepth)
		fo.Symbol = depth.Symbol
		fo.EventTime = eventTime
		fo.Price, serr = strconv.ParseFloat(order[0], 64)
		if serr != nil {
			err = fmt.Errorf("ask price: %v", serr)
			return
		}
		fo.Side = streamcommons.CommonFormatSell
		fo.Size, serr = strconv.ParseFloat(order[1], 64)
		if serr != nil {
			err = fmt.Errorf("ask size: %v", serr)
			return
		}
		mfo, serr := json.Marshal(fo)
		if serr != nil {
			err = fmt.Errorf("ask BinanceDepth: %v", serr)
			return
		}
		formatted[i] = Result{
			Channel: channel,
			Message: mfo,
		}
		i++
	}
	for _, order := range depth.Bids {
		fo := new(jsondef.BinanceDepth)
		fo.Symbol = depth.Symbol
		fo.EventTime = eventTime
		fo.Price, serr = strconv.ParseFloat(order[0], 64)
		if serr != nil {
			err = fmt.Errorf("bid price: %v", serr)
			return
		}
		fo.Side = streamcommons.CommonFormatBuy
		fo.Size, serr = strconv.ParseFloat(order[1], 64)
		if serr != nil {
			err = fmt.Errorf("bid size: %v", serr)
			return
		}
		mfo, serr := json.Marshal(fo)
		if serr != nil {
			err = fmt.Errorf("bid BinanceDepth: %v", serr)
			return
		}
		formatted[i] = Result{
			Channel: channel,
			Message: mfo,
		}
		i++
	}
	return
}

// FormatMessage formats messages from server.
func (f *binanceFormatter) FormatMessage(channel string, line []byte) (formatted []Result, err error) {
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: %v", serr)
		return
	}
	subscribed := new(jsonstructs.BinanceSubscribeResponse)
	serr = json.Unmarshal(line, subscribed)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: line: %v", serr)
		return
	}
	if subscribed.ID != 0 {
		// Subscribe message
		formatted = make([]Result, 1)
		switch stream {
		case streamcommons.BinanceStreamDepth:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceDepth,
			}
		case streamcommons.BinanceStreamTrade:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceTrade,
			}
		case streamcommons.BinanceStreamRESTDepth:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceRestDepth,
			}
		case streamcommons.BinanceStreamTicker:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceTicker,
			}
		default:
			err = fmt.Errorf("FormatMessage: channel not supported: %s", channel)
		}
		return
	}
	switch stream {
	case streamcommons.BinanceStreamDepth:
		return f.formatDepth(channel, line)
	case streamcommons.BinanceStreamRESTDepth:
		return f.formatRESTDepth(channel, line, symbol)
	case streamcommons.BinanceStreamTrade:
		return f.formatTrade(channel, line)
	case streamcommons.BinanceStreamTicker:
		return f.formatTicket(channel, line)
	default:
		err = fmt.Errorf("FormatMessage: unsupported: %v", channel)
		return
	}
}

// IsSupported returns true if the given channel is supported by this formatter.
func (f *binanceFormatter) IsSupported(channel string) bool {
	_, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return false
	}
	return stream == streamcommons.BinanceStreamDepth ||
		stream == streamcommons.BinanceStreamTrade ||
		stream == streamcommons.BinanceStreamTicker ||
		stream == streamcommons.BinanceStreamRESTDepth
}

func newBinanceFormatter() *binanceFormatter {
	return new(binanceFormatter)
}
