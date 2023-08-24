package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

func main() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.PUT("/trips", func(ctx *gin.Context) {
			addDocAsMap(ctx, client)
		})
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}

func createTripEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"trip": "created",
	})
}

func addDocAsMap(ctx context.Context, client *firestore.Client) error {
	_, err := client.Collection("cities").Doc("LA").Set(ctx, map[string]interface{}{
		"name":    "Los Angeles",
		"state":   "CA",
		"country": "USA",
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}
