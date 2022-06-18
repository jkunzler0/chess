package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

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

	// Note that u.Win and u.Loss will be 0 if they are not set in the JSON

	err = database.Put(u.UserID, database.Score{Win: u.Win, Loss: u.Loss})
	if err != nil {
		c.JSON(400, gin.H{"error": "could not get User"})
		return
	}

	transact.WritePut(u.UserID, database.Score{Win: u.Win, Loss: u.Loss})

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

func setupTransactionLog(cfg *config) error {
	var err error

	if cfg.localLog {
		// Initialize a local transaction log file
		transact, err = database.NewFileTransactionLogger("test.txt")
	} else {
		// Initialize a Postgres transaction log
		transact, err = database.NewPostgresTransactionLogger(cfg.dbParams)
	}
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	err = transact.StartTransactionLog()
	if err != nil {
		return fmt.Errorf("failed to start transaction logger: %w", err)
	}

	return err
}

// #######################################################################
// (Section 3) Main Fucntion  ############################################
// #######################################################################

func main() {

	help := flag.Bool("help", false, "Display Help")
	cfg := parseFlags()

	if *help {
		fmt.Printf("Config Info: %+v\n", cfg)
		os.Exit(0)
	}

	// Initializes the transaction log and loads any existing data
	// Blocks until all data is read
	err := setupTransactionLog(cfg)
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
	router.Run(fmt.Sprintf(":%d", cfg.listenPort))
}
