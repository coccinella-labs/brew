package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CloudAnalytics struct {
	db     *dynamodb.Client
	coffee *CoffeeMaker
}

func NewCloudAnalytics() *CloudAnalytics {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	return &CloudAnalytics{
		db:     dynamodb.NewFromConfig(cfg),
		coffee: NewCoffeeMaker(),
	}
}

func (ca *CloudAnalytics) logBrewCycle() {
	timestamp := time.Now().Unix()

	item := map[string]types.AttributeValue{
		"device_id":   &types.AttributeValueMemberS{Value: "coffee-maker-001"},
		"timestamp":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", timestamp)},
		"state":       &types.AttributeValueMemberS{Value: ca.coffee.state},
		"temperature": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ca.coffee.temperature)},
		"water_level": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ca.coffee.waterLevel)},
		"brew_time":   &types.AttributeValueMemberN{Value: "180"}, // 3 minutes
	}

	_, err := ca.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("coffee-analytics"),
		Item:      item,
	})

	if err != nil {
		log.Printf("DynamoDB error: %v", err)
	} else {
		fmt.Println("📊 Brew cycle logged to cloud")
	}
}

func (ca *CloudAnalytics) run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastState := ca.coffee.state

	for {
		select {
		case <-ticker.C:
			ca.coffee.tick()

			// Log when brew completes
			if lastState == "BREWING" && ca.coffee.state == "DONE" {
				ca.logBrewCycle()
			}

			lastState = ca.coffee.state
			fmt.Printf("State: %s | Temp: %d°C | Water: %d\n",
				ca.coffee.state, ca.coffee.temperature, ca.coffee.waterLevel)
		}
	}
}

func main() {
	analytics := NewCloudAnalytics()
	fmt.Println("☁️ Cloud analytics started")
	analytics.run()
}
