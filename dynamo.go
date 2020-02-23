package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"os"
	"time"
)

func getFood(region string, table string) ([]DBEntry, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession(&aws.Config{
		Region:aws.String(region),
	}))

	svc := dynamodb.New(sess)

	items := []DBEntry{}

	date := time.Now().Format("2006-01-02")

	filt := expression.Name("date").Equal(expression.Value(date))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		return items, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(os.Getenv("TABLE_NAME")),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	fmt.Println("Result", result)

	if err != nil {
		return items, err
	}

	for _, i := range result.Items {
		item := DBEntry{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return items, err
		}

		items = append(items, item)
	}

	return items, nil
}
