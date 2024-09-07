/*
 * BitMEX API
 *
 * ## REST API for the BitMEX Trading Platform  _If you are building automated tools, please subscribe to the_ _[BitMEX API RSS Feed](https://blog.bitmex.com/api_announcement/feed/) for changes. The feed will be updated_ _regularly and is the most reliable way to get downtime and update announcements._  [View Changelog](/app/apiChangelog)  ---  #### Getting Started  Base URI: [https://www.bitmex.com/api/v1](/api/v1)  ##### Fetching Data  All REST endpoints are documented below. You can try out any query right from this interface.  Most table queries accept `count`, `start`, and `reverse` params. Set `reverse=true` to get rows newest-first.  Additional documentation regarding filters, timestamps, and authentication is available in [the main API documentation](/app/restAPI).  _All_ table data is available via the [Websocket](/app/wsAPI). We highly recommend using the socket if you want to have the quickest possible data without being subject to ratelimits.  ##### Return Types  By default, all data is returned as JSON. Send `?_format=csv` to get CSV data or `?_format=xml` to get XML data.  ##### Trade Data Queries  _This is only a small subset of what is available, to get you started._  Fill in the parameters and click the `Try it out!` button to try any of these queries.  - [Pricing Data](#!/Quote/Quote_get)  - [Trade Data](#!/Trade/Trade_get)  - [OrderBook Data](#!/OrderBook/OrderBook_getL2)  - [Settlement Data](#!/Settlement/Settlement_get)  - [Exchange Statistics](#!/Stats/Stats_history)  Every function of the BitMEX.com platform is exposed here and documented. Many more functions are available.  ##### Swagger Specification  [⇩ Download Swagger JSON](swagger.json)  ---  ## All API Endpoints  Click to expand a section. 
 *
 * API version: 1.2.0
 * Contact: support@bitmex.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"time"
)

type UserPreferences struct {
	AlertOnLiquidations bool `json:"alertOnLiquidations,omitempty"`
	AnimationsEnabled bool `json:"animationsEnabled,omitempty"`
	AnnouncementsLastSeen time.Time `json:"announcementsLastSeen,omitempty"`
	ChatChannelID float64 `json:"chatChannelID,omitempty"`
	ColorTheme string `json:"colorTheme,omitempty"`
	Currency string `json:"currency,omitempty"`
	Debug bool `json:"debug,omitempty"`
	DisableEmails []string `json:"disableEmails,omitempty"`
	DisablePush []string `json:"disablePush,omitempty"`
	DisplayCorpEnrollUpsell bool `json:"displayCorpEnrollUpsell,omitempty"`
	EquivalentCurrency string `json:"equivalentCurrency,omitempty"`
	Features []string `json:"features,omitempty"`
	Favourites []string `json:"favourites,omitempty"`
	FavouritesAssets []string `json:"favouritesAssets,omitempty"`
	FavouritesOrdered []string `json:"favouritesOrdered,omitempty"`
	FavouriteBots []string `json:"favouriteBots,omitempty"`
	FavouriteContracts []string `json:"favouriteContracts,omitempty"`
	HasSetTradingCurrencies bool `json:"hasSetTradingCurrencies,omitempty"`
	HideConfirmDialogs []string `json:"hideConfirmDialogs,omitempty"`
	HideConnectionModal bool `json:"hideConnectionModal,omitempty"`
	HideFromLeaderboard bool `json:"hideFromLeaderboard,omitempty"`
	HideNameFromLeaderboard bool `json:"hideNameFromLeaderboard,omitempty"`
	HidePnlInGuilds bool `json:"hidePnlInGuilds,omitempty"`
	HideRoiInGuilds bool `json:"hideRoiInGuilds,omitempty"`
	HideNotifications []string `json:"hideNotifications,omitempty"`
	HidePhoneConfirm bool `json:"hidePhoneConfirm,omitempty"`
	GuidesShownVersion float32 `json:"guidesShownVersion,omitempty"`
	IsSensitiveInfoVisible bool `json:"isSensitiveInfoVisible,omitempty"`
	IsWalletZeroBalanceHidden bool `json:"isWalletZeroBalanceHidden,omitempty"`
	Locale string `json:"locale,omitempty"`
	LocaleSetTime float64 `json:"localeSetTime,omitempty"`
	MarginPnlRow string `json:"marginPnlRow,omitempty"`
	MarginPnlRowKind string `json:"marginPnlRowKind,omitempty"`
	MobileLocale string `json:"mobileLocale,omitempty"`
	MsgsSeen []string `json:"msgsSeen,omitempty"`
	Notifications interface{} `json:"notifications,omitempty"`
	OptionsBeta bool `json:"optionsBeta,omitempty"`
	OrderBookBinning interface{} `json:"orderBookBinning,omitempty"`
	OrderBookType string `json:"orderBookType,omitempty"`
	OrderClearImmediate bool `json:"orderClearImmediate,omitempty"`
	OrderControlsPlusMinus bool `json:"orderControlsPlusMinus,omitempty"`
	PlatformLayout string `json:"platformLayout,omitempty"`
	SelectedFiatCurrency string `json:"selectedFiatCurrency,omitempty"`
	ShowChartBottomToolbar bool `json:"showChartBottomToolbar,omitempty"`
	ShowLocaleNumbers bool `json:"showLocaleNumbers,omitempty"`
	Sounds []string `json:"sounds,omitempty"`
	SpacingPreference string `json:"spacingPreference,omitempty"`
	StrictIPCheck bool `json:"strictIPCheck,omitempty"`
	StrictTimeout bool `json:"strictTimeout,omitempty"`
	TickerGroup string `json:"tickerGroup,omitempty"`
	TickerPinned bool `json:"tickerPinned,omitempty"`
	TradeLayout string `json:"tradeLayout,omitempty"`
	UserColor string `json:"userColor,omitempty"`
}
