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
	dbHost          = "db"
	dbPort          = 5432
	pgTypeDate      = "date"
	pgTypeString    = "character varying"
	pgTypeBool      = "boolean"
	colAttrDataType = "datatype"
	tableSymbols    = "symbols"
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName     = os.Getenv("DB_NAME")
	token      = os.Getenv("IEX_TOKEN")
	psqlInfo   = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db         *sql.DB
)

type colDefinition struct {
	table_name  string
	column_name string
	data_type   string
}

func loadColumnDefinition(tableName string) (columnsDefn map[string]map[string]string) {
	sql := `SELECT  table_name, column_name, data_type
              FROM  information_schema.columns
             WHERE  table_name = $1 `
	rows, err := db.Query(sql, tableName)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}

	columnsDefn = make(map[string]map[string]string)
	for rows.Next() {
		var defn colDefinition
		err = rows.Scan(&defn.table_name, &defn.column_name, &defn.data_type)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(defn)
		columnsDefn[defn.column_name] = make(map[string]string)
		columnsDefn[defn.column_name][colAttrDataType] = defn.data_type
	}
	return
}

func convertToPgBool(value string) (convertedValue string) {
	// No need to convert yet
	log.Println(value)
	convertedValue = value
	log.Println(convertedValue)
	return
}

func convertToPgDate(value string) (convertedValue string) {
	// No need to convert yet
	convertedValue = value
	return
}

func convertToDbCompatible(tableName string, values map[string]interface{}) (convertedValues map[string]interface{}) {
	convertedValues = make(map[string]interface{})
	colDefn := loadColumnDefinition(tableName)
	log.Println(colDefn)
	for k, v := range values {
		log.Println(k, v)
		log.Println(colDefn[k][colAttrDataType])
		switch colDefn[k][colAttrDataType] {
		case pgTypeBool:
			log.Println(convertedValues)
			log.Println(v.(string))
			convertedValues[k] = convertToPgBool(v.(string))
			log.Println(convertedValues)
		case pgTypeDate:
			convertedValues[k] = convertToPgDate(v.(string))
		case pgTypeString:
			convertedValues[k] = v
		default:
			log.Println("Unknown Postgres data type", colDefn[k], "for column", tableName, ":", k)
		}
	}
	log.Println(convertedValues)
	return
}

func generateInsertSQLStatement(tableName string, values map[string]interface{}) (sql string) {
	keys := []string{}
	vals := []string{}
	colDefn := loadColumnDefinition(tableName)
	for k, v := range values {
		keys = append(keys, "\""+k+"\"")
		switch colDefn[k][colAttrDataType] {
		case pgTypeString:
			fallthrough
		case pgTypeDate:
			vals = append(vals, fmt.Sprintf("'%v'", v))
		default:
			vals = append(vals, fmt.Sprintf("%v", v))
		}
	}
	log.Println(keys)
	log.Println(vals)

	sql = "INSERT INTO iex.symbols (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(vals, ",") + ")"
	return
}

func loadSymbols(c *gin.Context) {
	log.Println("request [loadSymbols]")
	if len(token) == 0 {
		log.Fatalln("Failed to read IEX_SANDBOX_TOKEN environment variable")
	}

	url := "https://sandbox.iexapis.com/stable/ref-data/symbols?token=" + token

	iexClient := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "goiex")

	res, getErr := iexClient.Do(req)
	if getErr != nil {
		log.Fatalln(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	var symbols []map[string]interface{}
	jsonErr := json.Unmarshal(body, &symbols)
	if jsonErr != nil {
		log.Fatalln(jsonErr)
	}

	_, err = db.Exec("DELETE FROM iex.symbols")
	if err != nil {
		log.Fatalln(err)
	}

	for _, symbol := range symbols {
		log.Println("Inserting symbol: ", symbol)

		// convertedSymbol := convertToDbCompatible(tableSymbols, symbol)
		// log.Println("Converted to database compatible values: ", convertedSymbol)

		sql := generateInsertSQLStatement(tableSymbols, symbol)
		log.Println("Executing sql: ", sql)

		_, err := db.Exec(sql)
		if err != nil {
			log.Fatalln(err)
		}
	}
	//retMessage := fmt.Sprintf("%d symbols loaded into database.", len(symbols))
	c.IndentedJSON(http.StatusOK, symbols)
}

func getSymbols(c *gin.Context) {
	log.Println("request [getSymbols]")
	if len(token) == 0 {
		log.Fatalln("Failed to read IEX_SANDBOX_TOKEN environment variable")
	}

	sql := "SELECT * FROM iex.symbols"
	log.Println("Executing sql: ", sql)
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("rows:", rows)

	cols, err := rows.Columns()
	log.Println("columns:", cols)
	if err != nil {
		log.Fatalln(err)
	}

	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for idx := range cols {
		scanArgs[idx] = &values[idx]
	}

	var symbols []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("values:", values)
		result := make(map[string]interface{}, len(cols))
		for idx, col := range cols {
			result[col] = values[idx]
		}
		symbols = append(symbols, result)
	}
	c.IndentedJSON(http.StatusOK, symbols)
}

func ping(c *gin.Context) {
	log.Println("request [ping]")
	err := db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully ping", db.Driver())
	c.IndentedJSON(http.StatusOK, "Successfully ping "+dbName)
}

func main() {
	fileWriter, err := os.OpenFile("goiex.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	outputWriter := io.MultiWriter(fileWriter, os.Stdout)
	log.SetOutput(outputWriter)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	router := gin.Default()
	router.GET("/load", loadSymbols)
	router.GET("/symbols", getSymbols)
	router.GET("/ping", ping)
	router.Run(":8080")
}
