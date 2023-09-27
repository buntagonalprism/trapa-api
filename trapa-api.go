package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/buntagonalprism/trapa/api/common"
	"github.com/buntagonalprism/trapa/api/locations"
	"github.com/buntagonalprism/trapa/api/trips"
	"github.com/gin-gonic/gin"
)

var fbApp *firebase.App
var fsClient *firestore.Client
var authClient *fbAuth.Client

func init() {
	var err error
	fbApp, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	fsClient, err = fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore: %v\n", err)
	}

	authClient, err = fbApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firebase auth client: %v\n", err)
	}
}

func main() {

	defer fsClient.Close()

	router := gin.Default()
	router.Use(common.CORSMiddleware())
	router.Use(common.FirebaseAuth(authClient))

	v1 := router.Group("/v1")
	locationRouter := locations.NewLocationRouter(locations.NewLocationsService())
	locationRouter.RegisterRoutes(v1)
	tripRouter := trips.NewTripRouter(trips.NewTripService())
	tripRouter.RegisterRoutes(v1)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run("localhost:3000")
}
