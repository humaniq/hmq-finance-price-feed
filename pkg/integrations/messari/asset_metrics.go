package messari

import "time"

type assetMetricsResponse struct {
	Id                  string `json:"id"`
	SerialId            int    `json:"serial_id"`
	Symbol              string `json:"symbol"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	InternalTempAgoraId string `json:"_internal_temp_agora_id"`
	ContractAddresses   []struct {
		Platform        string `json:"platform"`
		ContractAddress string `json:"contract_address"`
	} `json:"contract_addresses"`
	MarketData struct {
		PriceUsd                               float64              `json:"price_usd"`
		PriceBtc                               float64              `json:"price_btc"`
		PriceEth                               float64              `json:"price_eth"`
		VolumeLast24Hours                      float64              `json:"volume_last_24_hours"`
		RealVolumeLast24Hours                  float64              `json:"real_volume_last_24_hours"`
		VolumeLast24HoursOverstatementMultiple float64              `json:"volume_last_24_hours_overstatement_multiple"`
		PercentChangeUsdLast1Hour              float64              `json:"percent_change_usd_last_1_hour"`
		PercentChangeBtcLast1Hour              float64              `json:"percent_change_btc_last_1_hour"`
		PercentChangeEthLast1Hour              float64              `json:"percent_change_eth_last_1_hour"`
		PercentChangeUsdLast24Hours            float64              `json:"percent_change_usd_last_24_hours"`
		PercentChangeBtcLast24Hours            float64              `json:"percent_change_btc_last_24_hours"`
		PercentChangeEthLast24Hours            float64              `json:"percent_change_eth_last_24_hours"`
		OhlcvLast1Hour                         *ohlcvResponseRecord `json:"ohlcv_last_1_hour"`
		OhlcvLast24Hour                        *ohlcvResponseRecord `json:"ohlcv_last_24_hour"`
		LastTradeAt                            time.Time
	} `json:"market_data"`
	Marketcap struct {
		Rank                             int     `json:"rank"`
		MarketcapDominancePercent        float64 `json:"marketcap_dominance_percent"`
		CurrentMarketcapUsd              float64 `json:"current_marketcap_usd"`
		Y2050MarketcapUsd                float64 `json:"y_2050_marketcap_usd"`
		YPlus10MarketcapUsd              float64 `json:"y_plus10_marketcap_usd"`
		LiquidMarketcapUsd               float64 `json:"liquid_marketcap_usd"`
		VolumeTurnoverLast24HoursPercent float64 `json:"volume_turnover_last_24_hours_percent"`
		RealizedMarketcapUsd             float64 `json:"realized_marketcap_usd"`
		OutstandingMarketcapUsd          float64 `json:"outstanding_marketcap_usd"`
	} `json:"marketcap"`
	Supply struct {
		Y2050                  float64 `json:"y_2050"`
		YPlus10                float64 `json:"y_plus10"`
		Liquid                 float64 `json:"liquid"`
		Circulating            float64 `json:"circulating"`
		Y250IssuedPercent      float64 `json:"y_250_issued_percent"`
		AnnualInflationPercent float64 `json:"annual_inflation_percent"`
		StockToFlow            float64 `json:"stock_to_flow"`
		YPlus10IssuedPercent   float64 `json:"y_plus10_issued_percent"`
		SupplyRevived90D       float64 `json:"supply_revived_90d"`
	} `json:"supply"`
}
type ohlcvResponseRecord struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	close  float64 `json:"close"`
	volume float64 `json:"volume"`
}
