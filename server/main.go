package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jkunzler0/chess/server/database"
	_ "github.com/lib/pq"
)

// #######################################################################
// (Section 1) Handlers ##################################################
// #######################################################################

type GameResult struct {
	WinnerID string `json:"WinnerID" binding:"required"`
	LoserID  string `json:"LoserID" binding:"required"`
}

type User struct {
	UserID string `json:"UserID" binding:"required"`
	Win    int    `json:"Win"`
	Loss   int    `json:"Loss"`
}

func postHandler(c *gin.Context) {
	var u User
	var err error

	if err = c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	err = database.Put(u.UserID, database.Score{u.Win, u.Loss})
	if err != nil {
		c.JSON(400, gin.H{"error": "could not get User"})
		return
	}

	transact.WritePut(u.UserID, database.Score{u.Win, u.Loss})

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func deleteHandler(c *gin.Context) {
	var u User
	var err error

	if err = c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	err = database.Delete(u.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "could not get User"})
		return
	}

	transact.WriteDelete(u.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func incrHandler(c *gin.Context) {
	var gr GameResult
	var err error

	if err = c.ShouldBindJSON(&gr); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	err = database.IncrWinLoss(gr.WinnerID, gr.LoserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "could not increment win/loss"})
		return
	}

	transact.WriteIncr(gr.WinnerID, gr.LoserID)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func getHandler(c *gin.Context) {
	var u User
	var err error
	var s database.Score

	if err = c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "could not bind JSON"})
		return
	}

	s, err = database.Get(u.UserID)
	if err != nil {
		c.JSON(400, gin.H{"error": "could not get User"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": s})
}

func getAllHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

// #######################################################################
// (Section 2) Transaction Log Setup  ####################################
// #######################################################################

var transact database.TransactionLogger

func setupTransactionLog() error {
	var err error

	transact, err = database.NewPostgresTransactionLogger(database.PostgresDbParams{
		Host:     "localhost",
		DBName:   "chessScoreboard",
		User:     "admin",
		Password: "pass123",
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	err = database.StartTransactionLog(transact)
	if err != nil {
		return fmt.Errorf("failed to start transaction logger: %w", err)
	}

	return err
}

// #######################################################################
// (Section 3) Main Fucntion  ############################################
// #######################################################################

func main() {

	// Initializes the transaction log and loads any existing data
	// Blocks until all data is read
	err := setupTransactionLog()
	if err != nil {
		panic(err)
	}

	// Setup gin router
	router := gin.Default()
	router.POST("/post", postHandler)
	router.POST("/delete", deleteHandler)
	router.POST("/incr", incrHandler)
	router.GET("/get", getHandler)
	router.GET("/getAll", getAllHandler)
	router.Run(":5000")
}
