package handlers

import (
	"encoding/json"
	"time"

	"go-server/models"
	"go-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
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

func generateSlug() string {
	return utils.RandomString(8)
}

func getPollById(rdb *redis.Client, pollID string) (models.Poll, error) {
	pollJSON, err := rdb.Get("poll:" + pollID).Result()
	if err != nil {
		return models.Poll{}, err
	}

	var poll models.Poll
	err = json.Unmarshal([]byte(pollJSON), &poll)
	if err != nil {
		return models.Poll{}, err
	}

	return poll, nil
}

func increaseOptionVotes(rdb *redis.Client, pollID string, optionID string) error {
	poll, err := getPollById(rdb, pollID)
	if err != nil {
		return err
	}

	for i, option := range poll.Options {
		if option.ID == optionID {
			poll.Options[i].Votes++
			break
		}
	}

	pollJSON, err := json.Marshal(poll)

	if err != nil {
		return err
	}

	err = rdb.Set("poll:"+poll.ID, pollJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (handler PollsHandler) CreatePoll(ctx *gin.Context) {
	rdb := handler.RedisClient
	var poll Poll

	err := ctx.BindJSON(&poll)
	poll.ID = uuid.NewString()
	poll.Slug = generateSlug()
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pollJSON, err := json.Marshal(poll)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = rdb.Set("poll:"+poll.ID, pollJSON, 0).Err()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"id": poll.ID, "slug": poll.Slug})
}

func (handler PollsHandler) GetPollsByUserUBH(ctx *gin.Context) {
	rdb := handler.RedisClient
	unique_browser_hash := ctx.Param("unique_browser_hash")

	pollIDs, err := rdb.Keys("poll:*").Result()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var polls []models.Poll
	for _, pollID := range pollIDs {
		pollJSON, err := rdb.Get(pollID).Result()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var poll models.Poll
		err = json.Unmarshal([]byte(pollJSON), &poll)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if poll.AuthorUniqueBrowserHash == unique_browser_hash {
			polls = append(polls, poll)
		}
	}

	ctx.JSON(200, polls)
}

func (handler PollsHandler) GetPoll(ctx *gin.Context) {
	rdb := handler.RedisClient
	pollID := ctx.Param("id")

	pollJSON, err := rdb.Get("poll:" + pollID).Result()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var poll models.Poll
	err = json.Unmarshal([]byte(pollJSON), &poll)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, poll)
}

func (handler PollsHandler) GetPollBySlug(ctx *gin.Context) {
	rdb := handler.RedisClient
	pollSlug := ctx.Param("slug")

	pollIDs, err := rdb.Keys("poll:*").Result()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var poll models.Poll
	for _, pollID := range pollIDs {

		pollJSON, err := rdb.Get(pollID).Result()
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = json.Unmarshal([]byte(pollJSON), &poll)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if poll.Slug == pollSlug {
			ctx.JSON(200, poll)
			return
		}
	}

	ctx.JSON(404, gin.H{"error": "poll not found"})
}

func (handler PollsHandler) UpdatePoll(ctx *gin.Context) {
	rdb := handler.RedisClient
	pollID := ctx.Param("id")

	var poll Poll

	err := ctx.BindJSON(&poll)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	poll.ID = pollID
	pollJSON, err := json.Marshal(poll)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = rdb.Set("poll:"+poll.ID, pollJSON, 0).Err()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, poll)
}

func (handler PollsHandler) DeletePoll(ctx *gin.Context) {
	rdb := handler.RedisClient
	pollID := ctx.Param("id")

	err := rdb.Del("poll:" + pollID).Err()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(204, nil)
}

//------------------------------------

// VOTE USE CASES ------------------------------------
func getVotesByPollId(rdb *redis.Client, pollID string) ([]models.Vote, error) {
	voteIDs, err := rdb.Keys("vote:*").Result()
	if err != nil {
		return nil, err
	}

	var votes []models.Vote
	for _, voteID := range voteIDs {
		voteJSON, err := rdb.Get(voteID).Result()
		if err != nil {
			return nil, err
		}

		var vote models.Vote
		err = json.Unmarshal([]byte(voteJSON), &vote)
		if err != nil {
			return nil, err
		}

		if vote.PollID == pollID {
			votes = append(votes, vote)
		}
	}

	return votes, nil
}

func (handler PollsHandler) Vote(ctx *gin.Context) {
	var vote Vote
	err := ctx.BindJSON(&vote)
	vote.ID = uuid.NewString()
	vote.UserIP = ctx.ClientIP()

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if vote.Option == "" {
		ctx.JSON(400, gin.H{"error": "An option is required."})
		return
	} else {
		// code to check if the option exists in the poll
		optionExists, err := doesOptionExist(handler.RedisClient, vote.PollID, vote.Option)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !optionExists {
			ctx.JSON(400, gin.H{"error": "Option does not exist in the poll."})
			return
		}

		// code to check if the user has already voted
		userHasVoted, err := hasUserVoted(handler.RedisClient, vote.PollID, vote)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if userHasVoted {
			ctx.JSON(400, gin.H{"error": "You have already voted!"})
			return
		}

		// code to check if the poll is still open
		pollEnded, err := hasPollEnded(handler.RedisClient, vote.PollID)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if pollEnded {
			ctx.JSON(400, gin.H{"error": "Poll is closed."})
			return
		}

		// increase vote
		err = increaseOptionVotes(handler.RedisClient, vote.PollID, vote.Option)

		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

	}

	voteJSON, err := json.Marshal(vote)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = handler.RedisClient.Set("vote:"+vote.ID, voteJSON, 0).Err()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"id": vote.ID})
}

func (handler PollsHandler) GetVotesByPollId(ctx *gin.Context) {
	pollID := ctx.Param("id")

	votes, err := getVotesByPollId(handler.RedisClient, pollID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, votes)
}

func (handler PollsHandler) HasUserVoted(ctx *gin.Context) {
	pollID := ctx.Param("poll_id")
	clientUBH := ctx.Param("client_ubh")
	userIP := ctx.ClientIP()

	votes, err := getVotesByPollId(handler.RedisClient, pollID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, vote := range votes {
		if vote.UserIP == userIP && vote.UserUniqueBrowserHash == clientUBH {
			ctx.JSON(200, gin.H{"voted": true})
			return
		}
	}

	ctx.JSON(200, gin.H{"voted": false})
}

func (handler PollsHandler) HasPollEnded(ctx *gin.Context) {
	pollID := ctx.Param("id")

	poll, err := getPollById(handler.RedisClient, pollID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if poll.EndDate == "" {
		ctx.JSON(200, gin.H{"ended": false})
		return
	}

	pollTime, err := time.Parse(time.RFC3339, poll.EndDate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if time.Now().After(pollTime) {
		ctx.JSON(200, gin.H{"ended": true})
		return
	}

	ctx.JSON(200, gin.H{"ended": false})
}

func hasUserVoted(rdb *redis.Client, pollID string, vote Vote) (bool, error) {
	votes, err := getVotesByPollId(rdb, pollID)
	if err != nil {
		return false, err
	}

	for _, v := range votes {
		if v.UserIP == vote.UserIP && v.UserUniqueBrowserHash == vote.UserUniqueBrowserHash {
			return true, nil
		}
	}

	return false, nil
}

func hasPollEnded(rdb *redis.Client, pollID string) (bool, error) {
	poll, err := getPollById(rdb, pollID)
	if err != nil {
		return false, err
	}

	if poll.EndDate != "" {
		pollTime, err := time.Parse(time.RFC3339, poll.EndDate)
		if err != nil {
			return false, err
		}
		return time.Now().After(pollTime), nil
	}
	return false, nil
}

func doesOptionExist(rdb *redis.Client, pollID string, optionID string) (bool, error) {
	poll, err := getPollById(rdb, pollID)
	if err != nil {
		return false, err
	}

	for _, option := range poll.Options {
		if option.ID == optionID {
			return true, nil
		}
	}

	return false, nil
}
