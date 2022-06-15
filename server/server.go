package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jkunzler0/chess/server/database"
)

type gameResult struct {
	winnerID string `json:"winnerID" binding:"required"`
	loserID  string `json:"loserID" binding:"required"`
}

type user struct {
	userID string `json:"userID" binding:"required"`
}

func post(c *gin.Context) {
	var gr gameResult
	var err error

	if err = c.ShouldBindJSON(&gr); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	err = database.IncrementWinLoss(gr.winnerID, gr.loserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "could not increment win/loss"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func get(c *gin.Context) {
	var u user
	var err error
	var s database.Score

	if err = c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	s, err = database.Get(u.userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "could not get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func main() {
	router := gin.Default()
	router.POST("/post", post)
	router.GET("/get", get)
	router.GET("/getAll", getAll)
	router.Run(":5000")
}
