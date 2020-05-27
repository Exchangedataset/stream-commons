package json

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/exchangedataset/streamcommons/formatter/json/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitmexFormatter formats message from bitmex
type BitmexFormatter struct {
}

func (f *BitmexFormatter) formatOrderBookL2(dataRaw json.RawMessage) (ret [][]byte, err error) {
	orders := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &orders)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}

	ret = make([][]byte, len(orders))
	for i, order := range orders {
		size := float64(order.Size)
		if order.Side == "Sell" {
			// if side is sell, negate size
			size = -size
		}
		marshaled, serr := json.Marshal(jsondef.BitmexOrderBookL2{
			Pair:  order.Symbol,
			Price: order.Price,
			ID:    order.ID,
			Size:  size,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatTrade(dataRaw json.RawMessage) (ret [][]byte, err error) {
	orders := make([]jsonstructs.BitmexTradeDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &orders)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	ret = make([][]byte, len(orders))
	for i, elem := range orders {
		size := float64(elem.Size)
		if elem.Side == "Sell" {
			size = -size
		}
		timestamp, serr := parseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexTrade{
			Pair:            elem.Symbol,
			Price:           elem.Price,
			Size:            size,
			Timestamp:       *timestamp,
			TrdMatchID:      elem.TradeMatchID,
			TickDirection:   elem.TickDirection,
			GrossValue:      elem.GrossValue,
			HomeNotional:    elem.HomeNotional,
			ForeignNotional: elem.ForeignNotional,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatInstrument(dataRaw json.RawMessage) (ret [][]byte, err error) {
	instruments := make([]jsonstructs.BitmexInstrumentDataElem, 0, 10)
	serr := json.Unmarshal(dataRaw, &instruments)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}

	ret = make([][]byte, len(instruments))
	for i, elem := range instruments {
		relistInterval, serr := parseDuration(elem.RelistInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for relistInterval failed: %s", serr.Error())
			return
		}
		calcInterval, serr := parseDuration(elem.CalcInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for calcInterval failed: %s", serr.Error())
			return
		}
		publishInterval, serr := parseDuration(elem.PublishInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for publishInterval failed: %s", serr.Error())
			return
		}
		publishTime, serr := parseDuration(elem.PublishTime)
		if serr != nil {
			err = fmt.Errorf("parse duration for publishTime failed: %s", serr.Error())
			return
		}
		fundingInterval, serr := parseDuration(elem.FundingInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for fundingInterval failed: %s", serr.Error())
			return
		}
		rebalanceInterval, serr := parseDuration(elem.RebalanceInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for rebalanceInterval failed: %s", serr.Error())
			return
		}
		sessionInterval, serr := parseDuration(elem.SessionInterval)
		if serr != nil {
			err = fmt.Errorf("parse duration for sessionInterval failed: %s", serr.Error())
			return
		}
		listing, serr := parseTimestamp(elem.Listing)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for listing failed: %s", serr.Error())
			return
		}
		front, serr := parseTimestamp(elem.Front)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for front failed: %s", serr.Error())
			return
		}
		expiry, serr := parseTimestamp(elem.Expiry)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for expiry failed: %s", serr.Error())
			return
		}
		settle, serr := parseTimestamp(elem.Settle)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for settle failed: %s", serr.Error())
			return
		}
		closingTimestamp, serr := parseTimestamp(elem.ClosingTimestamp)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for closingTimestamp failed: %s", serr.Error())
			return
		}
		fundingTimestamp, serr := parseTimestamp(elem.FundingTimestamp)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for fundingTimestamp failed: %s", serr.Error())
			return
		}
		openingTimestamp, serr := parseTimestamp(elem.OpeningTimestamp)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for openingTimestamp failed: %s", serr.Error())
			return
		}
		rebalanceTimestamp, serr := parseTimestamp(elem.RebalanceTimestamp)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for rebalanceTimestamp failed: %s", serr.Error())
			return
		}
		timestamp, serr := parseTimestamp(&elem.Timestamp)
		if serr != nil {
			err = fmt.Errorf("parse timestamp for timestamp failed: %s", serr.Error())
			return
		}

		marshaled, serr := json.Marshal(jsondef.BitmexInstrument{
			Symbol:                         elem.Symbol,
			RootSymbol:                     elem.RootSymbol,
			State:                          elem.State,
			Typ:                            elem.Typ,
			Listing:                        listing,
			Front:                          front,
			Expiry:                         expiry,
			Settle:                         settle,
			RelistInterval:                 relistInterval,
			InverseLeg:                     elem.InverseLeg,
			SellLeg:                        elem.SellLeg,
			BuyLeg:                         elem.BuyLeg,
			OptionStrikePcnt:               elem.OptionStrikePcnt,
			OptionStrikeRound:              elem.OptionStrikeRound,
			OptionStrikePrice:              elem.OptionStrikePrice,
			OptionMultiplier:               elem.OptionMultiplier,
			PositionCurrency:               elem.PositionCurrency,
			Underlying:                     elem.Underlying,
			QuoteCurrency:                  elem.QuoteCurrency,
			UnderlyingSymbol:               elem.UnderlyingSymbol,
			Reference:                      elem.Reference,
			ReferenceSymbol:                elem.ReferenceSymbol,
			CalcInterval:                   calcInterval,
			PublishInterval:                publishInterval,
			PublishTime:                    publishTime,
			MaxOrderQty:                    elem.MaxOrderQty,
			MaxPrice:                       elem.MaxPrice,
			LotSize:                        elem.LotSize,
			TickSize:                       elem.TickSize,
			Multiplier:                     elem.Multiplier,
			SettlCurrency:                  elem.SettlCurrency,
			UnderlyingToPositionMultiplier: elem.UnderlyingToPositionMultiplier,
			UnderlyingToSettleMultiplier:   elem.UnderlyingToSettleMultiplier,
			QuoteToSettleMultiplier:        elem.QuoteToSettleMultiplier,
			IsQuanto:                       elem.IsQuanto,
			IsInverse:                      elem.IsInverse,
			InitMargin:                     elem.InitMargin,
			MaintMargin:                    elem.MaintMargin,
			RiskLimit:                      elem.RiskLimit,
			RiskStep:                       elem.RiskStep,
			Limit:                          elem.Limit,
			Capped:                         elem.Capped,
			Taxed:                          elem.Taxed,
			Deleverage:                     elem.Deleverage,
			MakerFee:                       elem.MakerFee,
			TakerFee:                       elem.TakerFee,
			SettlementFee:                  elem.SettlementFee,
			InsuranceFee:                   elem.InsuranceFee,
			FundingBaseSymbol:              elem.FundingBaseSymbol,
			FundingQuoteSymbol:             elem.FundingQuoteSymbol,
			FundingPremiumSymbol:           elem.FundingPremiumSymbol,
			FundingTimestamp:               fundingTimestamp,
			FundingInterval:                fundingInterval,
			FundingRate:                    elem.FundingRate,
			IndicativeFundingRate:          elem.IndicativeFundingRate,
			RebalanceTimestamp:             rebalanceTimestamp,
			RebalanceInterval:              rebalanceInterval,
			OpeningTimestamp:               openingTimestamp,
			ClosingTimestamp:               closingTimestamp,
			SessionInterval:                sessionInterval,
			PrevClosePrice:                 elem.PrevClosePrice,
			LimitDownPrice:                 elem.LimitDownPrice,
			LimitUpPrice:                   elem.LimitUpPrice,
			BankruptLimitDownPrice:         elem.BankruptLimitDownPrice,
			BankruptLimitUpPrice:           elem.BankruptLimitUpPrice,
			PrevTotalVolume:                elem.PrevTotalVolume,
			TotalVolume:                    elem.TotalVolume,
			Volume:                         elem.Volume,
			Volume24h:                      elem.Volume24h,
			PrevTotalTurnover:              elem.PrevTotalTurnover,
			TotalTurnover:                  elem.TotalTurnover,
			Turnover:                       elem.Turnover,
			Turnover24h:                    elem.Turnover24h,
			HomeNotional24h:                elem.HomeNotional24h,
			ForeignNotional24h:             elem.ForeignNotional24h,
			PrevPrice24h:                   elem.PrevPrice24h,
			Vwap:                           elem.Vwap,
			HighPrice:                      elem.HighPrice,
			LowPrice:                       elem.LowPrice,
			LastPrice:                      elem.LastPrice,
			LastPriceProtected:             elem.LastPriceProtected,
			LastTickDirection:              elem.LastTickDirection,
			LastChangePcnt:                 elem.LastChangePcnt,
			BidPrice:                       elem.BidPrice,
			MidPrice:                       elem.MidPrice,
			AskPrice:                       elem.AskPrice,
			ImpactBidPrice:                 elem.ImpactBidPrice,
			ImpactMidPrice:                 elem.ImpactMidPrice,
			ImpactAskPrice:                 elem.ImpactAskPrice,
			HasLiquidity:                   elem.HasLiquidity,
			OpenInterest:                   elem.OpenInterest,
			OpenValue:                      elem.OpenValue,
			FairMethod:                     elem.FairMethod,
			FairBasisRate:                  elem.FairBasisRate,
			FairBasis:                      elem.FairBasis,
			FairPrice:                      elem.FairPrice,
			MarkMethod:                     elem.MarkMethod,
			MarkPrice:                      elem.MarkPrice,
			IndicativeTaxRate:              elem.IndicativeTaxRate,
			IndicativeSettlePrice:          elem.IndicativeSettlePrice,
			OptionUnderlyingPrice:          elem.OptionUnderlyingPrice,
			SettledPrice:                   elem.SettledPrice,
			Timestamp:                      *timestamp,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}

		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatLiquidation(dataRaw json.RawMessage) (ret [][]byte, err error) {
	liquidations := make([]jsonstructs.BitmexLiquidationDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &liquidations)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	ret = make([][]byte, len(liquidations))
	for i, elem := range liquidations {
		// Note: this structure is completely the same as what bitmex sends us...
		marshaled, serr := json.Marshal(jsondef.BitmexLiquidation{
			OrderID:   elem.OrderID,
			Symbol:    elem.Symbol,
			Side:      elem.Side,
			Price:     elem.Price,
			LeavesQty: elem.LeavesQty,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatSettlement(dataRaw json.RawMessage) (ret [][]byte, err error) {
	settlements := make([]jsonstructs.BitmexSettlementDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &settlements)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	ret = make([][]byte, len(settlements))
	for i, elem := range settlements {
		timestamp, serr := parseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexSettlement{
			Timestamp:             *timestamp,
			Symbol:                elem.Symbol,
			SettlementType:        elem.SettlementType,
			SettledPrice:          elem.SettledPrice,
			OptionStrikePrice:     elem.OptionStrikePrice,
			OptionUnderlyingPrice: elem.OptionUnderlyingPrice,
			Bankrupt:              elem.Bankrupt,
			TaxBase:               elem.TaxBase,
			TaxRate:               elem.TaxRate,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatInsurance(dataRaw json.RawMessage) (ret [][]byte, err error) {
	insurances := make([]jsonstructs.BitmexInsuranceDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &insurances)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	ret = make([][]byte, len(insurances))
	for i, elem := range insurances {
		timestamp, serr := parseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexInsurance{
			Currency:      elem.Currency,
			Timestamp:     *timestamp,
			WalletBalance: elem.WalletBalance,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

func (f *BitmexFormatter) formatFunding(dataRaw json.RawMessage) (ret [][]byte, err error) {
	fundings := make([]jsonstructs.BitmexFundingDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &fundings)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	ret = make([][]byte, len(fundings))
	for i, elem := range fundings {
		timestamp, serr := parseTimestamp(&elem.Timestamp)
		if serr != nil {
			err = fmt.Errorf("parsing timestamp failed for timestamp from funding: %v", serr)
		}
		fundingInterval, serr := parseDuration(&elem.FundingInterval)
		if serr != nil {
			err = fmt.Errorf("parsing timestamp failed for fundingInterval from funding: %v", serr)
		}
		marshaled, serr := json.Marshal(jsondef.BitmexFunding{
			Timestamp:        *timestamp,
			Symbol:           elem.Symbol,
			FundingInterval:  *fundingInterval,
			FundingRate:      elem.FundingRate,
			FundingRateDaily: elem.FundingRateDaily,
		})
		if serr != nil {
			err = fmt.Errorf("marshal failed: %s", serr.Error())
			return
		}
		ret[i] = marshaled
	}
	return
}

// Format formats incoming message given and returns formatted strings
func (f *BitmexFormatter) Format(channel string, line []byte) (ret [][]byte, err error) {
	subscribed := jsonstructs.BitmexSubscribe{}
	serr := json.Unmarshal(line, &subscribed)
	if serr != nil {
		err = fmt.Errorf("unmarshal failed: %s", serr.Error())
		return
	}
	if subscribed.Success {
		// this is a response to subscription
		// return header row
		if channel == "orderBookL2" {
			ret = [][]byte{jsondef.TypeDefBitmexOrderBookL2}
		} else if channel == "trade" {
			ret = [][]byte{jsondef.TypeDefBitmexTrade}
		} else if channel == "instrument" {
			ret = [][]byte{jsondef.TypeDefBitmexInstrument}
		} else if channel == "liquidation" {
			ret = [][]byte{jsondef.TypeDefBitmexLiquidation}
		} else if channel == "settlement" {
			ret = [][]byte{jsondef.TypeDefBitmexSettlement}
		} else if channel == "insurance" {
			ret = [][]byte{jsondef.TypeDefBitmexInsurance}
		} else if channel == "funding" {
			ret = [][]byte{jsondef.TypeDefBitmexFunding}
		} else {
			err = fmt.Errorf("unsupported channel '%s' for csvlike formatting", channel)
		}
		return
	}

	root := new(jsonstructs.BitmexRoot)
	serr = json.Unmarshal(line, root)
	if serr != nil {
		err = fmt.Errorf("root unmarshal failed: %s", serr.Error())
		return
	}

	if channel == "orderBookL2" {
		return f.formatOrderBookL2(root.Data)
	} else if channel == "trade" {
		return f.formatTrade(root.Data)
	} else if channel == "instrument" {
		return f.formatInstrument(root.Data)
	} else if channel == "liquidation" {
		return f.formatLiquidation(root.Data)
	} else if channel == "settlement" {
		return f.formatSettlement(root.Data)
	} else if channel == "insurance" {
		return f.formatInsurance(root.Data)
	} else if channel == "funding" {
		return f.formatFunding(root.Data)
	}

	err = fmt.Errorf("channel '%s' if not supported for csvlike formatting", channel)
	return
}

// IsSupported returns true if given channel is supported to be formatted using this formatter
func (f *BitmexFormatter) IsSupported(channel string) bool {
	return channel == "orderBookL2" || channel == "trade" || channel == "instrument" ||
		channel == "liquidation" || channel == "settlement" || channel == "insurance" ||
		channel == "funding"
}

func init() {
	var serr error
	durationBaseTime, serr = time.Parse(time.RFC3339Nano, "2000-01-01T00:00:00.000Z")
	if serr != nil {
		panic(fmt.Sprintf("initialization for durationBaseTime failed: %s", serr.Error()))
	}

}
