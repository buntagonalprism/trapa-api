package locations

import (
	"github.com/buntagonalprism/trapa/api/common"
)

type LocationsService interface {
	SearchPlaces(userId string, query string) ([]Place, error)
	GetPlaceDetails(userId string, placeId string) (PlaceDetails, error)
}

func NewLocationsService(cache common.Cache) LocationsService {
	return &locationsService{cache: cache}
}

type Place struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	Coordinates `json:"coordinates"`
}
type PlaceDetails struct {
	Place   `json:"place"`
	Website string `json:"website"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type locationsService struct {
	cache common.Cache
}

func (s *locationsService) SearchPlaces(userId string, query string) ([]Place, error) {
	sessionKey, _ := s.cache.Get(userSessionCacheKey(userId))
	print(sessionKey)
	return nil, nil
}

func (s *locationsService) GetPlaceDetails(userId string, placeId string) (PlaceDetails, error) {
	s.cache.Delete(userSessionCacheKey(userId))
	return PlaceDetails{}, nil
}

func userSessionCacheKey(userId string) string {
	return "userSessionKey_" + userId
}
