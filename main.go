package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var webhookURL = os.Getenv("WEBHOOK_URL")
var awsRegion = os.Getenv("DYNAMODB_REGION")
var awsTable = os.Getenv("DYNAMODB_TABLE")
var deepLTargetLang = os.Getenv("DEEPL_TARGET_LANG")
var deepLURL = os.Getenv("DEEPL_URL")
var deepLKey = os.Getenv("DEEPL_KEY")

var longMsg = "Today is " + time.Now().Weekday().String() + ", the *" + time.Now().Format("01/02/2006") + "*, here is today's lunch menu.\n*Enjoy!* :drooling_face:"
var shortMsg = "Here is today's lunch menu!"

// DBEntry is the entry in our DynamoDB table for a particular day.
type DBEntry struct {
	Canteen  string     `json:"canteen"`
	SpecDiet bool       `json:"spec_diet"`
	Date     string     `json:"date"`
	Items    []FoodItem `json:"items"`
}

// FoodItem is one menu item.
type FoodItem struct {
	Name       string `json:"name"`
	StudPrice  int    `json:"studprice"`
	ProfPrice  int    `json:"profprice"`
	Vegan      bool   `json:"vgn"`
	Vegetarian bool   `json:"vgt"`
	Fish       bool   `json:"fish"`
	Climate    bool   `json:"climate"`
}

// HandleRequest handles one request to the Lambda function.
func HandleRequest(ctx context.Context, event events.CloudWatchEvent) {

	timezone := os.Getenv("MENSA_TIMEZONE")

	tz, err := time.LoadLocation(timezone)

	if err != nil {
		log.Fatal(err)
	}

	// see if this event was triggered by the DST eventbridge rule
	if strings.Contains(event.Resources[0], "dst") != time.Now().In(tz).IsDST() {
		return
	}

	f, err := getFood(awsRegion, awsTable)

	if err != nil {
		panic(err)
	}

	f, err = translateFood(f, deepLTargetLang, deepLURL, deepLKey)

	if err != nil {
		panic(err)
	}

	msg := ""

	// every day is English wednesday
	msg = getMessage(f, longMsg, shortMsg)

	jsonStr := []byte(msg)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	log.Printf("sending %s to %s, got %d: %s", msg, webhookURL, resp.StatusCode, string(data))
}

func main() {
	lambda.Start(HandleRequest)
}
