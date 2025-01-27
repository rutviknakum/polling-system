package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var client *firestore.Client

func main() {
	r := gin.Default()

	if err := initFirestore(); err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer client.Close()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Polling System API",
			"endpoints": []string{
				"POST /create_poll - Create a new poll",
				"POST /vote - Vote in a poll",
				"GET /view_results - View poll results",
			},
		})
	})

	r.POST("/create_poll", createPoll)
	r.POST("/vote", vote)
	r.GET("/view_results", viewResults)

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initFirestore() error {
	ctx := context.Background()
	serviceAccountPath := "firebase/serviceAccountKey.json"
	var err error
	client, err = firestore.NewClient(ctx, "polling-system-4e5c6", option.WithCredentialsFile(serviceAccountPath))
	if err != nil {
		return err
	}
	log.Println("Firestore initialized successfully")
	return nil
}

func createPoll(c *gin.Context) {
	var poll struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
	}

	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if poll.Question == "" || len(poll.Options) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question and at least two options are required"})
		return
	}

	ctx := context.Background()
	votes := make(map[string]int)
	for _, option := range poll.Options {
		votes[option] = 0
	}

	doc, _, err := client.Collection("polls").Add(ctx, map[string]interface{}{
		"question": poll.Question,
		"options":  poll.Options,
		"votes":    votes,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create poll"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"poll_id": doc.ID, "message": "Poll created successfully!"})
}

func vote(c *gin.Context) {
	var voteData struct {
		PollID         string `json:"poll_id"`
		SelectedOption string `json:"selected_option"`
	}

	if err := c.ShouldBindJSON(&voteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if voteData.PollID == "" || voteData.SelectedOption == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID and Selected Option are required"})
		return
	}

	ctx := context.Background()
	docRef := client.Collection("polls").Doc(voteData.PollID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	data := doc.Data()
	votes, ok := data["votes"].(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process votes"})
		return
	}

	if _, exists := votes[voteData.SelectedOption]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid option selected"})
		return
	}

	currentVotes, ok := votes[voteData.SelectedOption].(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid vote count type"})
		return
	}

	votes[voteData.SelectedOption] = currentVotes + 1

	_, err = docRef.Set(ctx, map[string]interface{}{
		"votes": votes,
	}, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your vote has been recorded!"})
}

func viewResults(c *gin.Context) {
	pollID := c.Query("poll_id")
	if pollID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Poll ID is required"})
		return
	}

	ctx := context.Background()
	doc, err := client.Collection("polls").Doc(pollID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poll not found"})
		return
	}

	data := doc.Data()
	c.JSON(http.StatusOK, gin.H{
		"poll_id":  pollID,
		"question": data["question"],
		"results":  data["votes"],
	})
}
