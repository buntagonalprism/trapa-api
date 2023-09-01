package main

type Trip struct {
	Name          string `json:"name" binding:"required" firestore:"name"`
	StartDate     string `json:"startDate" binding:"required,date" firestore:"startDate"`
	EndDate       string `json:"endDate" binding:"required,date" firestore:"endDate"`
	SingleCountry string `json:"singleCountry" firestore:"singleCountry,omitempty"`
	Owner         string `firestore:"owner"`
}
