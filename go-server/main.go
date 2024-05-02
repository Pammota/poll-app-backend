package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"go-server/handlers"
)

func main() {
	password := os.Getenv("REDIS_PASSWORD")
	// password := "redditpassword"

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	// client.Set("testKey", "hello", 0)
	got := client.Get("anotherKey")

	fmt.Println(got.Val())

	handler := handlers.CreatePollsHandler(client)

	router := gin.Default()

	// CORS Setup

	config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://localhost:3000", "http://polls.localhost:3000", "https://polls.gligor.dev"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "active",
		})
	})

	// Polls CRUD
	router.GET("/polls/user/:unique_browser_hash", handler.GetPollsByUserUBH)
	router.POST("/polls", handler.CreatePoll)
	router.GET("/polls/:id", handler.GetPoll)
	router.GET("/polls/slug/:slug", handler.GetPollBySlug)
	router.PUT("/polls/:id", handler.UpdatePoll)
	router.DELETE("/polls/:id", handler.DeletePoll)
	router.GET("/polls/:id/hasEnded", handler.HasPollEnded)

	// Vote use cases
	router.POST("/polls/:id/vote", handler.Vote)
	router.GET("/polls/:id/votes", handler.GetVotesByPollId)
	router.GET("/polls/:id/votes/:user_id", handler.HasUserVoted)

	router.Run("0.0.0.0:8090")
}
