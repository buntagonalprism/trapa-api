package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/buntagonalprism/trapa/api/common"
	"github.com/buntagonalprism/trapa/api/locations"
	"github.com/buntagonalprism/trapa/api/trips"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
)

var fbApp *firebase.App
var fsClient *firestore.Client
var authClient *fbAuth.Client
var mapsClient *maps.Client

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

	mapsKeyFilePath := os.Getenv("GOOGLE_MAPS_API_KEY_FILE")
	mapsKeyData, _ := os.ReadFile(mapsKeyFilePath)
	mapsKey := string(mapsKeyData)
	mapsClient, err = maps.NewClient(maps.WithAPIKey(mapsKey))
	if err != nil {
		log.Fatalf("error initializing Google Maps client: %v\n", err)
	}
}

func main() {

	defer fsClient.Close()

	router := gin.Default()
	router.Use(common.CORSMiddleware())
	router.Use(common.FirebaseAuth(authClient))

	cache := common.NewCache()
	locationRouter := locations.NewLocationRouter(locations.NewLocationsService(cache, mapsClient))
	tripRouter := trips.NewTripRouter(trips.NewTripService())

	v1 := router.Group("/v1")
	locationRouter.RegisterRoutes(v1)
	tripRouter.RegisterRoutes(v1)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run("localhost:3000")
}
