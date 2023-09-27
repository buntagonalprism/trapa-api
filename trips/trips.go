package trips

import (
	"log"

	"cloud.google.com/go/firestore"
	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/buntagonalprism/trapa/api/common"
	"github.com/gin-gonic/gin"
)

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

type TripService interface {
	CreateTrip(ctx *gin.Context, req CreateTripRequest) (*Trip, error)
}

func NewTripService() TripService {
	return &tripService{}
}

type tripService struct {
	fs *firestore.Client
}

func (s *tripService) CreateTrip(ctx *gin.Context, req CreateTripRequest) (*Trip, error) {
	newDoc := s.fs.Collection("trips").NewDoc()
	userId := ctx.MustGet(common.FirebaseUserKey).(*fbAuth.Token).UID
	trip := &Trip{
		Name:          req.Name,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		SingleCountry: req.SingleCountryCode,
		Owner:         userId,
		Editors:       []string{userId},
	}
	_, err := newDoc.Set(ctx, trip)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	trip.Id = newDoc.ID
	return trip, err
}
