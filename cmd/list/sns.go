package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/spf13/cobra"
)

func runSns(cmd *cobra.Command, args []string) {
	sdkConfig, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(os.Getenv("AWS_SNS_REGION")),
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
	fmt.Println("Let's list SNS for your account.")
	var subscriptions []types.Subscription
	snsClient := sns.NewFromConfig(sdkConfig)
	paginator := sns.NewListSubscriptionsPaginator(snsClient, &sns.ListSubscriptionsInput{})
	for paginator.HasMorePages() {
		output, e := paginator.NextPage(context.TODO())
		if e != nil {
			err = e
			break
		}
		subscriptions = append(subscriptions, output.Subscriptions...)
	}
	if err != nil {
		log.Fatalf("List queue error: %v", err)
	}
	if len(subscriptions) <= 0 {
		fmt.Println("You don't have any subscriptions!")
		os.Exit(0)
	}
	for _, queueUrl := range subscriptions {
		fmt.Printf("\t%v\n", queueUrl)
	}
}
