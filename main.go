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

var url = os.Getenv("WEBHOOK_URL")
var awsRegion = os.Getenv("DYNAMODB_REGION")
var awsTable = os.Getenv("DYNAMODB_TABLE")

// var msgDeu = "Heute ist" + time.Now().Weekday().String() + ", der *" + time.Now().Format("02.01.2006") + "*, hier ist das Mittagsmenü für heute.\n*Guten Appetit!* :drooling_face:"
// var shortDeu = "Hier ist das Mittagsmenü für heute!"
var msgEng = "Today is " + time.Now().Weekday().String() + ", the *" + time.Now().Format("01/02/2006") + "*, here is today's lunch menu.\n*Enjoy!* :drooling_face:"
var shortEng = "Here is today's lunch menu!"

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

	msg := ""

	// if time.Now().Weekday().String() == "Wednesday" {
	// 	msg = getMessage(f, msgEng, shortEng)
	// } else {
	// 	msg = getMessage(f, msgDeu, shortDeu)
	// }

	// every day is English wednesday
	msg = getMessage(f, msgEng, shortEng)

	jsonStr := []byte(msg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

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

	log.Printf("sending %s to %s, got %d: %s", msg, url, resp.StatusCode, string(data))
}

func main() {
	lambda.Start(HandleRequest)
}
