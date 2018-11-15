package model

import "time"

// RequestInput struct
type RequestInput struct {
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

// Request struct
type Request struct {
	DateFrom time.Time
	DateTo   time.Time
}

// ======================== Qraphql Structs ================================ //

// Report struct
type Report struct {
	FuelSales *FuelSales `json:"fuelSaleRangeSummary"`
}

// FuelSales struct
type FuelSales struct {
	DateStart int
	DateEnd   int
	DateFrom  time.Time
	DateTo    time.Time
	Sales     *[]SalesSummary `json:"summary"`
}

// SalesSummary struct
type SalesSummary struct {
	StationID   string `json:"stationID"`
	StationName string `json:"stationName"`
	HasDSL      bool
	Fuels       struct {
		NL  float64
		DSL float64
	} `json:"fuels"`
}
