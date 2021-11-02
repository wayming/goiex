package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	token := os.Getenv("IEX_TOKEN")
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
}
