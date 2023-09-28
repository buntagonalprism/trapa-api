package locations

import (
	"github.com/buntagonalprism/trapa/api/common"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"googlemaps.github.io/maps"
)

type LocationsService interface {
	SearchPlaces(ctx *gin.Context, query string, country string) ([]Place, error)
	GetPlaceDetails(ctx *gin.Context, placeId string) (PlaceDetails, error)
}

func NewLocationsService(cache common.Cache, mapsClient *maps.Client) LocationsService {
	return &locationsService{cache: cache, mapsClient: mapsClient}
}

type Place struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type PlaceDetails struct {
	Place       `json:"place"`
	Coordinates `json:"coordinates"`
	Website     string `json:"website"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type locationsService struct {
	cache      common.Cache
	mapsClient *maps.Client
}

func (s *locationsService) SearchPlaces(ctx *gin.Context, query string, country string) ([]Place, error) {
	userId := ctx.MustGet(common.FirebaseUserIdKey).(string)
	sessionKey, _ := s.cache.Get(userSessionCacheKey(userId))
	if sessionKey == nil {
		sessionKey = uuid.New()
		s.cache.Set(userSessionCacheKey(userId), sessionKey)
	}
	request := &maps.PlaceAutocompleteRequest{
		Input:        query,
		SessionToken: maps.PlaceAutocompleteSessionToken(sessionKey.(uuid.UUID)),
		Components:   map[maps.Component][]string{maps.Component("country"): {country}},
	}
	response, err := s.mapsClient.PlaceAutocomplete(ctx, request)
	if err != nil {
		return nil, err
	}
	places := make([]Place, len(response.Predictions))
	for i, prediction := range response.Predictions {
		places[i] = Place{
			Id:   prediction.PlaceID,
			Name: prediction.Description,
		}
	}
	return places, nil
}

func (s *locationsService) GetPlaceDetails(ctx *gin.Context, placeId string) (PlaceDetails, error) {
	userId := ctx.MustGet(common.FirebaseUserIdKey).(string)
	sessionKey, _ := s.cache.Get(userSessionCacheKey(userId))
	if sessionKey == nil {
		sessionKey = uuid.New()
	}
	s.cache.Delete(userSessionCacheKey(userId))

	request := &maps.PlaceDetailsRequest{
		PlaceID:      placeId,
		SessionToken: maps.PlaceAutocompleteSessionToken(sessionKey.(uuid.UUID)),
	}
	response, err := s.mapsClient.PlaceDetails(ctx, request)
	if err != nil {
		return PlaceDetails{}, err
	}
	placeDetails := PlaceDetails{
		Place: Place{
			Id:   response.PlaceID,
			Name: response.Name,
		},
		Coordinates: Coordinates{
			Latitude:  response.Geometry.Location.Lat,
			Longitude: response.Geometry.Location.Lng,
		},
	}
	return placeDetails, nil
}

func userSessionCacheKey(userId string) string {
	return "userSessionKey_" + userId
}
