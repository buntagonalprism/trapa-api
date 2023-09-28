package trips

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/buntagonalprism/trapa/api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type tripServiceFake struct {
	trip Trip
	err  error
}

func (s *tripServiceFake) CreateTrip(ctx *gin.Context, req CreateTripRequest) (*Trip, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &s.trip, nil
}

func TestCreateTripHandlerErrors(t *testing.T) {
	s := &tripServiceFake{
		trip: Trip{
			Id:        "trip123",
			Name:      "test",
			StartDate: "2021-11-19",
			EndDate:   "2021-11-21",
			Owner:     "testId123",
			Editors:   []string{"testId123"},
		},
	}
	r := NewTripRouter(s)

	g := gin.New()
	g.Use(SetUser("testId123"))

	r.RegisterRoutes(g.Group(""))

	tests := []struct {
		name     string
		reqBody  string
		respCode int
		respBody string
	}{
		{
			name:     "Empty request body",
			reqBody:  "",
			respCode: 400,
			respBody: `{"error":"Request body must not be empty"}`,
		},
		{
			name:     "Invalid JSON",
			reqBody:  `{"name": "test"`,
			respCode: 400,
			respBody: `{"error":"Invalid JSON"}`,
		},
		{
			name:     "Missing fields",
			reqBody:  `{}`,
			respCode: 400,
			respBody: `{"fieldErrors":[{"field":"name","error":"name is required"},{"field":"startDate","error":"startDate is required"},{"field":"endDate","error":"endDate is required"}]}`,
		},
		{
			name:     "Invalid dates",
			reqBody:  `{"name": "test", "startDate": "foobar", "endDate": "11-19-21"}`,
			respCode: 400,
			respBody: `{"fieldErrors":[{"field":"startDate","error":"startDate must be a valid date in the format YYYY-MM-DD"},{"field":"endDate","error":"endDate must be a valid date in the format YYYY-MM-DD"}]}`,
		},
		{
			name:     "Valid request",
			reqBody:  `{"name": "test", "startDate": "2021-11-19", "endDate": "2021-11-21"}`,
			respCode: 200,
			respBody: `{"id":"trip123","name":"test","startDate":"2021-11-19","endDate":"2021-11-21","owner":"testId123","editors":["testId123"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			g.ServeHTTP(w, httptest.NewRequest("PUT", "/trips", bytes.NewBufferString(tt.reqBody)))
			assert.Equal(t, tt.respCode, w.Code)
			assert.Equal(t, tt.respBody, w.Body.String())
		})
	}
}

func SetUser(userId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.FirebaseUserIdKey, userId)
	}
}
