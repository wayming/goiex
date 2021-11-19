package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const (
	dbHost = "db"
	dbPort = 5432
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	token      = os.Getenv("IEX_TOKEN")
)

func getSymbols(c *gin.Context) {

	if len(token) == 0 {
		log.Fatal("Failed to read IEX_SANDBOX_TOKEN environment variable")
	}

	url := "https://sandbox.iexapis.com/stable/ref-data/symbols?token=" + token

	iexClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "goiex")

	res, getErr := iexClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var symbols []map[string]interface{}
	jsonErr := json.Unmarshal(body, &symbols)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	for _, symbol := range symbols {
		fmt.Println(symbol["symbol"])
	}
	c.IndentedJSON(http.StatusOK, symbols)
}
func ping(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to ", psqlInfo)
	c.IndentedJSON(http.StatusOK, psqlInfo)
}
func main() {
	router := gin.Default()
	router.GET("/symbols", getSymbols)
	router.GET("/ping", ping)
	router.Run(":8080")
}
