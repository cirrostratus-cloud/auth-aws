package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/cirrostratus-cloud/auth-aws/user/repository"
	user_service "github.com/cirrostratus-cloud/auth-aws/user/service"
	"github.com/cirrostratus-cloud/auth/user"
	"github.com/cirrostratus-cloud/common/event"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

var re *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("arn:aws:sqs:[a-z]{2}-[a-z]*-[0-9]:[0-9]{12}:%s_user_(.*)", os.Getenv("CIRROSTRATUS_AUTH_MODULE_NAME")))

var snsEventBus *user_service.SNSEventBus

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	log.SetOutput(os.Stdout)
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Error loading AWS config")
		panic(err)
	}
	topicArnPrefix := os.Getenv("TOPIC_ARN_PREFIX")
	if topicArnPrefix == "" {
		log.Fatal("TOPIC_ARN_PREFIX is required")
		panic("TOPIC_ARN_PREFIX is required")
	}
	emailFrom := os.Getenv("SES_EMAIL_FROM")
	if emailFrom == "" {
		log.Fatal("SES_EMAIL_FROM is required")
		panic("SES_EMAIL_FROM is required")
	}
	emailConfirmationURL := os.Getenv("EMAIL_CONFIRMATION_URL")
	if emailConfirmationURL == "" {
		log.Fatal("EMAIL_CONFIRMATION_URL is required")
		panic("EMAIL_CONFIRMATION_URL is required")
	}
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY is required")
		panic("PRIVATE_KEY is required")
	}
	maxAgeInSeconds, err := strconv.Atoi(os.Getenv("MAX_AGE_IN_SECONDS"))
	if err != nil {
		log.Fatal("Error parsing MAX_AGE_IN_SECONDS")
		panic(err)
	}
	snsClient := sns.NewFromConfig(cfg)
	sesClient := ses.NewFromConfig(cfg)
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	snsEventBus = user_service.NewSNSEventBus(snsClient, topicArnPrefix)
	moduleName := os.Getenv("CIRROSTRATUS_AUTH_MODULE_NAME")
	if moduleName == "" {
		log.Fatal("CIRROSTRATUS_AUTH_MODULE_NAME is required")
		panic("CIRROSTRATUS_AUTH_MODULE_NAME is required")
	}
	userTable := os.Getenv("CIRROSTRATUS_AUTH_USER_TABLE")
	if userTable == "" {
		log.Fatal("CIRROSTRATUS_AUTH_USER_TABLE is required")
		panic("CIRROSTRATUS_AUTH_USER_TABLE is required")
	}
	tableName := fmt.Sprintf("%s-%s", moduleName, userTable)
	userRepository := repository.NewDynamoUserRepository(dynamodbClient, tableName)
	mailService := user_service.NewSESMailService(sesClient)
	user.NewNotifyPasswordChangedService(userRepository, mailService, snsEventBus, emailFrom)
	user.NewNotifyPasswordRecoveredService(userRepository, mailService, snsEventBus, emailFrom)
	user.NewNotifyUserCreatedService(userRepository, mailService, snsEventBus, emailFrom)
	user.NewNotifyEmailConfirmationService(userRepository, mailService, snsEventBus, emailFrom, emailConfirmationURL, []byte(privateKey), maxAgeInSeconds)
}

func handler(ctx context.Context, req events.SQSEvent) {
	for _, record := range req.Records {
		log.WithFields(log.Fields{
			"MessageID":      record.MessageId,
			"Body":           record.Body,
			"EventSourceARN": record.EventSourceARN,
		}).Info("Received message")
		var recordBody map[string]interface{}
		err := json.Unmarshal([]byte(record.Body), &recordBody)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error unmarshalling record body")
			return
		}
		queueArn := record.EventSourceARN
		matches := re.FindStringSubmatch(queueArn)
		if matches == nil {
			log.WithFields(log.Fields{
				"QueueArn": queueArn,
			}).Error("Regex did not match queueArn")
			return
		}
		if len(matches) < 2 {
			log.WithFields(log.Fields{
				"QueueArn": queueArn,
			}).Error("Error parsing queueArn")
			return
		}
		eventName := "user/" + matches[1]
		log.WithFields(log.Fields{
			"EventName": eventName,
		}).Info("Triggering event")
		err = snsEventBus.Trigger(event.EventName(eventName), recordBody["Message"].(string))
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error triggering event")
		}
	}
}

func main() {
	stage := os.Getenv("AWS_STAGE")
	if stage == "local" {
		address := os.Getenv("USER_EVENT_ADDR")
		app := fiber.New()
		app.Post("/", func(c *fiber.Ctx) error {
			var payload map[string]interface{}
			err := json.Unmarshal([]byte(c.Body()), &payload)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error unmarshalling record body")
				return err
			}
			eventName := payload["eventName"]
			message, err := json.Marshal(payload["message"])
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Error marshalling message")
				return err
			}
			return snsEventBus.Trigger(event.EventName(eventName.(string)), string(message))
		})
		log.Fatal(app.Listen(address))
	} else {
		lambda.Start(handler)
	}
}
