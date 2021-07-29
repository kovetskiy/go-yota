package yota

import "time"

type DevicesInfo struct {
	HasSliderDevice         bool              `json:"hasSliderDevice"`
	DeviceOperationDisabled bool              `json:"deviceOperationDisabled"`
	AuthLevelHigh           bool              `json:"authLevelHigh"`
	CurrentRegion           string            `json:"currentRegion"`
	Devices                 []Devices         `json:"devices"`
	CurrentDeviceNotNull    bool              `json:"currentDeviceNotNull"`
	ShowActivationButton    bool              `json:"showActivationButton"`
	CurrentResourceID       CurrentResourceID `json:"currentResourceId"`
}
type Price struct {
	Amount       int    `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}
type Product struct {
	ProductID           int       `json:"productId"`
	Status              string    `json:"status"`
	BeginDate           time.Time `json:"beginDate"`
	EndDate             time.Time `json:"endDate"`
	ExternalID          int       `json:"externalId"`
	ProductOfferingCode string    `json:"productOfferingCode"`
	TarifficationType   string    `json:"tarifficationType"`
	Price               Price     `json:"price"`
	AllowAutoProlong    bool      `json:"allowAutoProlong"`
}
type ResourceID struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}
type RegistrState struct {
	Name                string    `json:"name"`
	RegisterDate        time.Time `json:"registerDate"`
	AnonymousManagement bool      `json:"anonymousManagement"`
	RegionCode          string    `json:"regionCode"`
}
type PhysicalResource struct {
	PhysicalResourceID int          `json:"physicalResourceId"`
	ResourceSpecCode   string       `json:"resourceSpecCode"`
	Created            time.Time    `json:"created"`
	Updated            time.Time    `json:"updated"`
	ResourceID         ResourceID   `json:"resourceID"`
	RegistrState       RegistrState `json:"registrState"`
}
type OfferingSpeed struct {
	SpeedValue    string `json:"speedValue"`
	UnitOfMeasure string `json:"unitOfMeasure"`
}
type SpecialOffers struct {
	DefaultSelectedPoint bool   `json:"defaultSelectedPoint"`
	Position             int    `json:"position"`
	Speed                string `json:"speed"`
	SpeedType            string `json:"speedType"`
	Code                 string `json:"code"`
	SpecialOffer         bool   `json:"specialOffer"`
	SpeedMaximum         bool   `json:"speedMaximum"`
	TheShortestOffer     bool   `json:"theShortestOffer"`
	Light                bool   `json:"light"`
	OfferDescription     string `json:"offerDescription"`
	Amount               int    `json:"amount"`
	Period               int    `json:"period"`
	PeriodString         string `json:"periodString"`
	MoneyEnough          bool   `json:"moneyEnough"`
	PayFromBalance       int    `json:"payFromBalance"`
	PayFromCard          int    `json:"payFromCard"`
	PayFromAll           int    `json:"payFromAll"`
	Remain               int    `json:"remain"`
	RemainString         string `json:"remainString"`
	ReturnAmount         int    `json:"returnAmount"`
	TestDrive            bool   `json:"testDrive"`
}
type Steps struct {
	DefaultSelectedPoint bool    `json:"defaultSelectedPoint"`
	Position             float64 `json:"position"`
	Speed                string  `json:"speed"`
	SpeedType            string  `json:"speedType"`
	Code                 string  `json:"code"`
	SpecialOffer         bool    `json:"specialOffer"`
	SpeedMaximum         bool    `json:"speedMaximum"`
	TheShortestOffer     bool    `json:"theShortestOffer"`
	Light                bool    `json:"light"`
	OfferDescription     string  `json:"offerDescription"`
	Amount               int     `json:"amount"`
	Period               int     `json:"period"`
	PeriodString         string  `json:"periodString"`
	MoneyEnough          bool    `json:"moneyEnough"`
	PayFromBalance       int     `json:"payFromBalance"`
	PayFromCard          int     `json:"payFromCard"`
	PayFromAll           int     `json:"payFromAll"`
	Remain               int     `json:"remain"`
	RemainString         string  `json:"remainString"`
	ReturnAmount         float64 `json:"returnAmount"`
	DisablingAutoprolong bool    `json:"disablingAutoprolong,omitempty"`
	TestDrive            bool    `json:"testDrive"`
	OfferDisabled        bool    `json:"offerDisabled,omitempty"`
}
type CurrentProduct struct {
	DefaultSelectedPoint bool    `json:"defaultSelectedPoint"`
	Position             float64 `json:"position"`
	Speed                string  `json:"speed"`
	SpeedType            string  `json:"speedType"`
	Code                 string  `json:"code"`
	SpecialOffer         bool    `json:"specialOffer"`
	SpeedMaximum         bool    `json:"speedMaximum"`
	TheShortestOffer     bool    `json:"theShortestOffer"`
	Light                bool    `json:"light"`
	OfferDescription     string  `json:"offerDescription"`
	Amount               int     `json:"amount"`
	Period               int     `json:"period"`
	PeriodString         string  `json:"periodString"`
	MoneyEnough          bool    `json:"moneyEnough"`
	PayFromBalance       int     `json:"payFromBalance"`
	PayFromCard          int     `json:"payFromCard"`
	PayFromAll           int     `json:"payFromAll"`
	Remain               int     `json:"remain"`
	RemainString         string  `json:"remainString"`
	OfferDisabled        bool    `json:"offerDisabled"`
	ReturnAmount         int     `json:"returnAmount"`
	TestDrive            bool    `json:"testDrive"`
}
type Slider struct {
	ProductID      int             `json:"productId"`
	TestDrive      bool            `json:"testDrive"`
	SpecialOffers  []SpecialOffers `json:"specialOffers"`
	Steps          []Steps         `json:"steps"`
	CurrentProduct CurrentProduct  `json:"currentProduct"`
}
type Devices struct {
	ActiveDevice          bool             `json:"activeDevice"`
	ProductStatus         string           `json:"productStatus"`
	Product               Product          `json:"product"`
	HomeRegion            bool             `json:"homeRegion"`
	SliderProduct         bool             `json:"sliderProduct"`
	PhysicalResource      PhysicalResource `json:"physicalResource"`
	OfferingSpeedMaximum  bool             `json:"offeringSpeedMaximum"`
	AutoprolongFlag       bool             `json:"autoprolongFlag"`
	SpecialOffersExpanded bool             `json:"specialOffersExpanded"`
	TurboEnabled          bool             `json:"turboEnabled"`
	OfferingSpeed         OfferingSpeed    `json:"offeringSpeed"`
	Light                 bool             `json:"light"`
	SliderEnabled         bool             `json:"sliderEnabled"`
	TestDrive             bool             `json:"testDrive"`
	Slot                  bool             `json:"slot"`
	Slider                Slider           `json:"slider"`
}
type CurrentResourceID struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type Balance struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type UserInfo struct {
	AccountType string    `json:"accountType"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Registered  time.Time `json:"registered"`
	Status      string    `json:"status"`
	UserID      int       `json:"userId"`
}

type CurrentInfo struct {
	BeginDate time.Time
	EndDate   time.Time
	Price     Price
	Speed     OfferingSpeed
}

type ChangeTariff struct {
	CurrentProductID     int    `json:"currentProductId"`
	DisablingAutoprolong bool   `json:"disablingAutoprolong"`
	OfferCode            string `json:"offerCode"`
	ResourceID           struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	} `json:"resourceID"`
}

type PaymentsInfo struct {
	Payments []Payment `json:"payments"`
}
type Payment struct {
	DatePayment time.Time `json:"datePayment"`
	Amount      int       `json:"amount"`
	OfferCode   string    `json:"offerCode"`
	OfferName   string    `json:"offerName"`
	NameDevice  string    `json:"nameDevice"`
	Iccid       string    `json:"iccid"`
}

type OperationHistory struct {
	UserID        int64     `json:"userId"`
	OperationType string    `json:"operationType"`
	DepositType   string    `json:"depositType,omitempty"`
	DepositSource string    `json:"depositSource,omitempty"`
	ActualDate    time.Time `json:"actualDate"`
	Amount        struct {
		Amount       float64 `json:"amount"`
		CurrencyCode string  `json:"currencyCode"`
	} `json:"amount"`
	WriteDownInitOperation string `json:"writeDownInitOperation,omitempty"`
	WriteDownType          string `json:"writeDownType,omitempty"`
	OfferCode              string `json:"offerCode,omitempty"`
	OfferName              string `json:"offerName,omitempty"`
	ResourceName           string `json:"resourceName,omitempty"`
	Iccid                  string `json:"iccid,omitempty"`
}
