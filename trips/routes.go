package trips

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TripRouter struct {
	tripService TripService
}

func NewTripRouter(tripService TripService) *TripRouter {
	return &TripRouter{tripService: tripService}
}

func (r *TripRouter) RegisterRoutes(router *gin.RouterGroup) {
	g := router.Group("/trips")
	g.PUT("", r.createTripHandler)
}

func (r *TripRouter) createTripHandler(c *gin.Context) {
	var req CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": getFieldDisplayErrors(err)})
		return
	}

	trip, err := r.tripService.CreateTrip(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, trip)
}
