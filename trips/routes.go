package trips

import (
	"net/http"

	"github.com/buntagonalprism/trapa/api/common"
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
	req, err := common.BindJson[CreateTripRequest](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	trip, err := r.tripService.CreateTrip(c, *req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(200, trip)
}
