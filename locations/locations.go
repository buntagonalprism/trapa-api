package locations

type LocationsService interface {
	SearchPlaces(query string, session LocationSearchSession) ([]Place, LocationSearchSession, error)
	GetPlaceDetails(placeId string, session LocationSearchSession) (PlaceDetails, error)
}

func NewLocationsService() LocationsService {
	return &locationsService{}
}

type Place struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
	Coordinates `json:"coordinates"`
}

type LocationSearchSession struct {
	sessionKey string
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
}

func (s *locationsService) SearchPlaces(query string, session LocationSearchSession) ([]Place, LocationSearchSession, error) {
	return nil, LocationSearchSession{}, nil
}

func (s *locationsService) GetPlaceDetails(placeId string, session LocationSearchSession) (PlaceDetails, error) {
	return PlaceDetails{}, nil
}
