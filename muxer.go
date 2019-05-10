package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	type CurRates struct {
		USD float32 `json:"USD"`
		JPY float32 `json:"JPY"`
		RUB float32 `json:"RUB"`
		MAD float32 `json:"MAD"`
		CAD float32 `json:"CAD"`
		AUD float32 `json:"AUD"`
		IRR float32 `json:"IRR"`
		STD float32 `json:"STD"`
		MXN float32 `json:"MXN"`
		INR float32 `json:"INR"`
		ZWD float32 `json:"ZWD"`
		NPR float32 `json:"NPR"`
		EUR float32 `json:"EUR"`
	}

	type Currency struct {
		TimeStamp int      `json:"timestamp"`
		CurRates  CurRates `json:"rates"`
		Base      string   `json:"base"`
		Date      string   `json:"date"`
	}
	//From my worker, here is the json structure

	currentsession := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	//Starts a new AWS connection
	svc := dynamodb.New(currentsession)
	//Connect to dynamodb
	muxer := mux.NewRouter()
	muxer.HandleFunc("/mfernan4/status", func(writer http.ResponseWriter, muxer *http.Request) {
		//Assign endpoint
		req := &dynamodb.DescribeTableInput{TableName: aws.String("CurrencyDB")}
		//Finds dynamodb table
		data, err := svc.DescribeTable(req)
		//Finds my table
		if err != nil {
			fmt.Print("err 1:   %s", err.Error())
		}
		fmt.Fprintf(writer, "", data.Table)
		//Make data for endpoint
	})

	//This is all the same as above, but we are working on the other endpoint
	muxer.HandleFunc("/mfernan4/all", func(writer http.ResponseWriter, muxer *http.Request) {
		req := &dynamodb.ScanInput{TableName: aws.String("CurrencyDB")}
		data, err := svc.Scan(req)
		if err != nil {
			fmt.Print("err 2:   %s", err.Error())
		}
		entitydata := []Currency{}
		dynamodbattribute.UnmarshalListOfMaps(data.Items, &entitydata)
		json.NewEncoder(writer).Encode(entitydata)
		//Define structure for the data from structs
		//Uses the "data" and the defined structure to unmarshal the objects from the database

	})

	for {
		http.ListenAndServe(":8080", muxer)
	}
}

