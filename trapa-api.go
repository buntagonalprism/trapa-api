package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

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

	authClient, err = fbApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firebase auth client: %v\n", err)
	}

	fsClient, err = firestore.NewClient(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("error initializing Firestore: %v\n", err)
	}

	mapsKeySecret := fmt.Sprintf("projects/%s/secrets/places-api-key/versions/latest", os.Getenv("GOOGLE_CLOUD_PROJECT_NUMBER"))
	mapsKeyBytes, err := AccessSecretVersion(mapsKeySecret)
	if err != nil {
		log.Fatalf("error getting Google Places API key: %v\n", err)
	}
	mapsKey := string(mapsKeyBytes)
	mapsClient, err = maps.NewClient(maps.WithAPIKey(mapsKey))
	if err != nil {
		log.Fatalf("error initializing Google Maps client: %v\n", err)
	}
}

func main() {

	defer fsClient.Close()

	router := gin.Default()
	router.Use(common.CORSMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	cache := common.NewCache()
	locationRouter := locations.NewLocationRouter(locations.NewLocationsService(cache, mapsClient))
	tripRouter := trips.NewTripRouter(trips.NewTripService(fsClient))

	v1 := router.Group("/v1")
	v1.Use(common.FirebaseAuth(authClient))
	locationRouter.RegisterRoutes(v1)
	tripRouter.RegisterRoutes(v1)

	port := os.Getenv("PORT")
	if runtime.GOOS == "windows" {
		// For local windows development, to avoid needing to run as admin
		router.Run("localhost:" + port)
	} else {
		router.Run(":" + port)
	}

}
