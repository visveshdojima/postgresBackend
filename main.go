package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {

	type Person struct {
		ID   string
		Name string
	}

	type Txn struct {
	}

	// To fetch data from api
	response, err := http.Get("http://localhost:1317/cosmos/tx/v1beta1/txs/block/97")

	if err != nil {
		fmt.Print(err.Error())
		// os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewSharedCredentials("/Users/visveshnaraharisetty/.aws/credentials", "default"),
	})

	if err != nil {
		fmt.Printf("%s", err)
	}

	svc := dynamodb.New(sess)

	req := &dynamodb.DescribeTableInput{
		TableName: aws.String("Person"),
	}

	result, err := svc.DescribeTable(req)

	if err != nil {
		fmt.Printf("%s", err)
	}

	table := result.Table
	fmt.Printf("%s", table)

	params := &dynamodb.ScanInput{
		TableName: aws.String("Person"),
	}

	r, e := svc.Scan(params)

	if e != nil {
		fmt.Errorf("failed to Make query %v", e)
	}

	o := []Person{}

	e = dynamodbattribute.UnmarshalListOfMaps(r.Items, &o)

	if e != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", e))
	}

	fmt.Println(o)

	rec := Person{
		ID:   "11",
		Name: "QQA",
	}

	av, err := dynamodbattribute.MarshalMap(rec)

	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	tableName := "Person"

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

}
