package main

import (
	"fmt"
	"os"

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

	fmt.Println("Redis credential: ", password)
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	handler := handlers.CreatePollsHandler(client)

	router := gin.Default()

	// Polls CRUD
	router.GET("/polls", handler.GetPolls)
	router.POST("/polls", handler.CreatePoll)
	router.GET("/polls/:id", handler.GetPoll)
	router.PUT("/polls/:id", handler.UpdatePoll)
	router.DELETE("/polls/:id", handler.DeletePoll)

	// Vote use cases
	router.POST("/polls/:id/vote", handler.Vote)
	router.GET("/polls/:id/votes", handler.GetVotesByPollId)
	router.GET("/polls/:id/votes/:user_id", handler.HasUserVoted)

	router.Run("0.0.0.0:8080")
}
