package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
	psqlInfo   = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db         *sql.DB
)

func loadSymbols(c *gin.Context) {
	log.Println("request [loadSymbols]")
	if len(token) == 0 {
		log.Fatal("Failed to read IEX_SANDBOX_TOKEN environment variable")
	}

	url := "https://cloud.iexapis.com//stable/ref-data/symbols?token=" + token

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
		log.Println("Inserting symbol: ", symbol)
		keys := []string{}
		vals := []string{}
		for k, v := range symbol {
			keys = append(keys, k)
			vals = append(vals, fmt.Sprintf("%v", v))
		}
		log.Println(keys)
		log.Println(vals)

		sql := "INSERT INTO iex.symbols (" + strings.Join(keys, ",") + ") VALUES(" + strings.Join(vals, ",") + ")"
		log.Println("Execute sql:", sql)
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}
	//retMessage := fmt.Sprintf("%d symbols loaded into database.", len(symbols))
	c.IndentedJSON(http.StatusOK, symbols)
}

func getSymbols(c *gin.Context) {
	log.Println("request [getSymbols]")
	if len(token) == 0 {
		log.Fatal("Failed to read IEX_SANDBOX_TOKEN environment variable")
	}
	var symbols []map[string]interface{}
	c.IndentedJSON(http.StatusOK, symbols)
}

func ping(c *gin.Context) {
	log.Println("request [ping]")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully ping", db.Driver())
	c.IndentedJSON(http.StatusOK, "Successfully ping "+dbName)
}

func main() {
	fileWriter, err := os.OpenFile("goiex.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	outputWriter := io.MultiWriter(fileWriter, os.Stdout)
	log.SetOutput(outputWriter)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.GET("/load", loadSymbols)
	router.GET("/symbols", getSymbols)
	router.GET("/ping", ping)
	router.Run(":8080")
}
