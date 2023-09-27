package locations

import (
	"github.com/buntagonalprism/trapa/api/common"
	"github.com/gin-gonic/gin"
)

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
	userId := c.MustGet(common.FirebaseUserIdKey).(string)
	places, _ := r.locationsService.SearchPlaces(userId, query)
	print(places)
}

func (r *LocationRouter) getPlaceDetailsHandler(c *gin.Context) {
	placeId := c.Param("placeId")
	userId := c.MustGet(common.FirebaseUserIdKey).(string)
	details, _ := r.locationsService.GetPlaceDetails(userId, placeId)
	print(details.Name)
}
