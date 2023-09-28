package locations

import (
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
	if query == "" {
		c.JSON(400, gin.H{"error": "query parameter is required"})
		return
	}
	country := c.Query("country")
	if country == "" {
		c.JSON(400, gin.H{"error": "country parameter is required"})
		return
	}
	places, err := r.locationsService.SearchPlaces(c, query, country)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, places)
}

func (r *LocationRouter) getPlaceDetailsHandler(c *gin.Context) {
	placeId := c.Param("placeId")
	if placeId == "" {
		c.JSON(400, gin.H{"error": "placeId parameter is required"})
		return
	}
	details, _ := r.locationsService.GetPlaceDetails(c, placeId)
	c.JSON(200, details)
}
