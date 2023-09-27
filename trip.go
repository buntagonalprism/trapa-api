package main

type Trip struct {
	Id            string   `json:"id"               firestore:"-"`
	Name          string   `json:"name"             firestore:"name"`
	StartDate     string   `json:"startDate"        firestore:"startDate"`
	EndDate       string   `json:"endDate"          firestore:"endDate"`
	SingleCountry string   `json:"singleCountry"    firestore:"singleCountry,omitempty"`
	Owner         string   `json:"owner"            firestore:"owner"`
	Editors       []string `json:"editors"          firestore:"editors"`
}

type CreateTripRequest struct {
	Name              string `json:"name"               binding:"required"`
	StartDate         string `json:"startDate"          binding:"required,date"`
	EndDate           string `json:"endDate"            binding:"required,date"`
	SingleCountryCode string `json:"singleCountryCode"`
}
