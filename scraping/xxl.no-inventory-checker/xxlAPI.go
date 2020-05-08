package main

// StockReply is the reply to requests like these:
// curl 'https://www.xxl.no/p/1159507_1_Size_3/stores/stock'
type StockReply struct {
	SelectedSize        string `json:"selectedSize"`
	StoreAvailabilities []struct {
		StockStatus string `json:"stockStatus"`
		StoreData   struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			GeoPoint    struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"geoPoint"`
			Address struct {
				Line1                string `json:"line1"`
				Town                 string `json:"town"`
				PostalCode           string `json:"postalCode"`
				Phone                string `json:"phone"`
				Email                string `json:"email"`
				ShippingAddress      bool   `json:"shippingAddress"`
				DefaultAddress       bool   `json:"defaultAddress"`
				VisibleInAddressBook bool   `json:"visibleInAddressBook"`
				FormattedAddress     bool   `json:"formattedAddress"`
			} `json:"address"`
			OpeningHours struct {
				SpecialDayOpeningList []interface{} `json:"specialDayOpeningList"`
				WeekDayOpeningList    []struct {
					Closed      bool   `json:"closed"`
					ClosingTime string `json:"closingTime,omitempty"`
					OpeningTime string `json:"openingTime,omitempty"`
					WeekDay     string `json:"weekDay"`
				} `json:"weekDayOpeningList"`
			} `json:"openingHours"`
			GooglePlaceID string `json:"googlePlaceId"`
			LinkButton    struct {
				DisplayName string `json:"displayName"`
				URL         string `json:"url"`
			} `json:"linkButton"`
		} `json:"storeData"`
	} `json:"storeAvailabilities"`
}
