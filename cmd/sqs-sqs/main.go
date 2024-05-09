package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/ThreeDotsLabs/watermill-amazonsqs/connection"
	"github.com/ThreeDotsLabs/watermill-amazonsqs/sqs"
)

const myTopic = "my-sqs-test"

func main() {
	ctx := context.Background()
	logger := watermill.NewStdLogger(true, true)
	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(os.Getenv("AWS_SQS_REGION")),
		connection.SetEndPoint(os.Getenv("AWS_SQS_ENDPOINT")),
		awsconfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
				SecretAccessKey: os.Getenv("AWS_ACCESS_SECRET"),
			},
		}),
	)
	if err != nil {
		log.Fatalf("LoadDefaultConfig() error: %v", err)
	}
	pub, err := sqs.NewPublisher(sqs.PublisherConfig{
		AWSConfig:              cfg,
		CreateQueueIfNotExists: true,
		Marshaler:              sqs.DefaultMarshalerUnmarshaler{},
	}, logger)
	if err != nil {
		log.Fatalf("NewPublisher() error: %v", err)
	}
	_ = pub

	sub, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig:                    cfg,
		CreateQueueInitializerConfig: sqs.QueueConfigAtrributes{},
		Unmarshaler:                  sqs.DefaultMarshalerUnmarshaler{},
	}, logger)
	if err != nil {
		log.Fatalf("NewSubscriber() error: %v", err)
	}

	err = sub.SubscribeInitialize(myTopic)
	if err != nil {
		log.Printf("SubscribeInitialize() error: %v", err)
	}

	messages, err := sub.Subscribe(ctx, myTopic)
	if err != nil {
		log.Fatalf("Subscribe() error: %v", err)
	}

	go func() {
		for m := range messages {
			logger.Info("Received message", watermill.LogFields{"UUID": m.UUID, "Metadata": m.Metadata, "Payload": string(m.Payload)})
			m.Ack()
		}
	}()

	for {
		msg := message.NewMessage(watermill.NewULID(), []byte(`{"some_json": "body"}`))
		err := pub.Publish(myTopic, msg)
		if err != nil {
			log.Fatalf("Publish() error: %v", err)
		}
		time.Sleep(time.Second)
	}
}
