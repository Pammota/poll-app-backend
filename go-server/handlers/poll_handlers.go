package handlers

import (
	"go-server/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Vote = models.Vote
type Options = models.Options
type Poll = models.Poll

type PollsHandler struct {
	RedisClient *redis.Client
}

func CreatePollsHandler(client *redis.Client) PollsHandler {
	return PollsHandler{
		RedisClient: client,
	}
}

// POLLS CRUD ------------------------------------

func (PollsHandler) CreatePoll(GinContext *gin.Context) {
}

func (PollsHandler) GetPolls(GinContext *gin.Context) {
	// code to retrieve polls from Redis
}

func (PollsHandler) GetPoll(GinContext *gin.Context) {
	// code to retrieve a single poll from Redis
}

func (PollsHandler) UpdatePoll(GinContext *gin.Context) {
	// code to update a poll in Redis
}

func (PollsHandler) DeletePoll(GinContext *gin.Context) {
	// code to delete a poll from Redis
}

//------------------------------------

// VOTE USE CASES ------------------------------------
func (PollsHandler) Vote(GinContext *gin.Context) {
	// code to vote on a poll and save the vote to Redis
}

func (PollsHandler) GetVotesByPollId(GinContext *gin.Context) {
	// code to retrieve votes for a poll from Redis
}

func (PollsHandler) HasUserVoted(GinContext *gin.Context) {
	// code to retrieve votes for a user from Redis
}
