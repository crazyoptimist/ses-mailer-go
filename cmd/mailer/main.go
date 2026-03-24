package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("unable to load .env file: %v", err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to get workdir: %v", err)
	}

	f, err := excelize.OpenFile(filepath.Join(path, "emails.xlsx"))
	if err != nil {
		log.Fatalf("unable to open excel file: %v", err)
		return
	}
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Fatalf("unable to get rows: %v", err)
		return
	}

	accessKeyId := viper.GetString("AWS_ACCESS_KEY_ID")
	secretAccessKey := viper.GetString("AWS_SECRET_ACCESS_KEY")
	region := viper.GetString("AWS_REGION")
	session := ""

	// Setup Static Credentials
	staticProvider := credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, session)

	// Load Config with the credentials provider
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(staticProvider),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create the SES v2 client
	client := sesv2.NewFromConfig(cfg)

	FROM_EMAIL := viper.GetString("FROM_EMAIL")
	REPLY_TO_EMAIL := viper.GetString("REPLY_TO_EMAIL")
	EMAIL_SUBJECT := viper.GetString("EMAIL_SUBJECT")
	EMAIL_TEXT_BODY := viper.GetString("EMAIL_TEXT_BODY")

	for _, row := range rows {
		input := ComposeEmail(FROM_EMAIL, []string{REPLY_TO_EMAIL}, []string{row[0]}, EMAIL_SUBJECT, EMAIL_TEXT_BODY)

		// Send the email
		out, err := client.SendEmail(context.TODO(), input)
		if err != nil {
			log.Printf("failed to send email, %v", err)
		}

		fmt.Printf("Email sent successfully! Message ID: %s\n", *out.MessageId)
	}
}

func ComposeEmail(fromEmail string, replyToAdresses []string, toAddresses []string, emailSubject string, emailText string) *sesv2.SendEmailInput {
	return &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(fromEmail),
		ReplyToAddresses: replyToAdresses,
		Destination: &types.Destination{
			ToAddresses: toAddresses,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(emailSubject),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(emailText),
					},
				},
			},
		},
	}
}
