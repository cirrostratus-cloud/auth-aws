package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/cirrostratus-cloud/auth-aws/user/repository"
	user_service "github.com/cirrostratus-cloud/auth-aws/user/service"
	"github.com/cirrostratus-cloud/auth/user"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

var fiberLambda *fiberadapter.FiberLambda
var app *fiber.App

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
	stage := os.Getenv("AWS_STAGE")
	inValue, err := strconv.Atoi(os.Getenv("USER_MIN_PASSWORD_LENGTH"))
	if err != nil {
		log.Fatal("Error parsing USER_MIN_PASSWORD_LENGTH")
		panic(err)
	}
	minPasswordLength := inValue
	boolValue, err := strconv.ParseBool(os.Getenv("USER_UPPER_CASE_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_UPPER_CASE_REQUIRED")
		panic(err)
	}
	upperCaseRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_LOWER_CASE_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_LOWER_CASE_REQUIRED")
		panic(err)
	}
	lowerCaseRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_NUMBER_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_NUMBER_REQUIRED")
		panic(err)
	}
	numberRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_SPECIAL_CHARACTER_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_SPECIAL_CHARACTER_REQUIRED")
		panic(err)
	}
	specialCharacterRequired := boolValue
	topicArnPrefix := os.Getenv("TOPIC_ARN_PREFIX")
	if topicArnPrefix == "" {
		log.Fatal("TOPIC_ARN_PREFIX is required")
		panic("TOPIC_ARN_PREFIX is required")
	}
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("PRIVATE_KEY is required")
		panic("PRIVATE_KEY is required")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Error loading AWS config")
		panic(err)
	}
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
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)
	snsEventBus := user_service.NewSNSEventBus(snsClient, topicArnPrefix)
	userRepository := repository.NewDynamoUserRepository(dynamodbClient, tableName)
	validatePasswordService := user.NewValidatePasswordService(userRepository, upperCaseRequired, lowerCaseRequired, numberRequired, specialCharacterRequired, minPasswordLength)
	createUserService := user.NewCreateUserService(userRepository, snsEventBus, validatePasswordService)
	getUserUseCase := user.NewGetUserService(userRepository)
	updateProfileUseCase := user.NewUpdateUserProfileService(userRepository)
	confirmateEmailService := user.NewConfirmateEmailService(userRepository, snsEventBus, []byte(privateKey))
	userAPI := newUserAPI(createUserService, getUserUseCase, updateProfileUseCase, confirmateEmailService)
	app = fiber.New()
	userAPI.setUp(app, stage)
	fiberLambda = fiberadapter.New(app)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.
		WithField("Path", req.RequestContext.Path).
		Info("Processing request.")
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	stage := os.Getenv("AWS_STAGE")
	if stage == "local" {
		address := os.Getenv("USER_HTTP_ADDR")
		log.Fatal(app.Listen(address))
	} else {
		lambda.Start(handler)
	}
}
