package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	fbAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

var app *firebase.App
var fs *firestore.Client
var auth *fbAuth.Client

func init() {
	var err error
	app, err = firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v\n", err)
	}

	fs, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firestore: %v\n", err)
	}

	auth, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initializing Firebase auth client: %v\n", err)
	}
}

func main() {

	defer fs.Close()

	router := gin.Default()
	router.Use(FirebaseAuth())

	v1 := router.Group("/v1")
	{
		v1.PUT("/trips", createTripHandler)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run("localhost:3000")
}

func createTripHandler(c *gin.Context) {
	var trip Trip
	if err := c.ShouldBindJSON(&trip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": getFieldDisplayErrors(err)})
		return
	}

	tripId, err := addTrip(c, trip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"trip_id": tripId,
	})
}

func addTrip(ctx *gin.Context, trip Trip) (string, error) {
	newDoc := fs.Collection("trips").NewDoc()
	trip.Owner = ctx.MustGet(FirebaseUserKey).(*fbAuth.Token).UID
	_, err := newDoc.Set(ctx, trip)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return newDoc.ID, err
}

func FirebaseAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": "Authorization header missing"})
			return
		}
		idToken, err := extractBearerToken(authorizationHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": err})
			return
		}
		token, err := auth.VerifyIDToken(c, idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"unauthorized": err})
			return
		}
		c.Set(FirebaseUserKey, token)
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

const FirebaseUserKey = "firebase_user"
