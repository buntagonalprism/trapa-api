package locations

import "github.com/gin-gonic/gin"

type LocationRouter struct {
	locationsService LocationsService
}

func NewLocationRouter(locationsService LocationsService) *LocationRouter {
	return &LocationRouter{locationsService: locationsService}
}

func (r *LocationRouter) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/locations")
	g.GET("/search", r.searchPlacesHandler)
	g.GET("/:placeId", r.getPlaceDetailsHandler)
}

func (r *LocationRouter) searchPlacesHandler(c *gin.Context) {
	query := c.Query("query")
	session := LocationSearchSession{sessionKey: c.Query("sessionKey")}
	places, _, _ := r.locationsService.SearchPlaces(query, session)
	print(places)
}

func (r *LocationRouter) getPlaceDetailsHandler(c *gin.Context) {
	placeId := c.Param("placeId")
	session := LocationSearchSession{sessionKey: c.Query("sessionKey")}
	details, _ := r.locationsService.GetPlaceDetails(placeId, session)
	print(details.Name)
}
