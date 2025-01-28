package app

import (
	"cynxhostagent/internal/dependencies"
	"cynxhostagent/internal/repository/micro/cynxhostcentral"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	Logger *logrus.Logger
	Config *dependencies.Config

	Validator *validator.Validate

	RedisClient    *redis.Client
	DatabaseClient *dependencies.DatabaseClient
	AWSClient      *dependencies.AWSClient

	CynxhostCentral *cynxhostcentral.CynxhostCentral

	JWTManager    *dependencies.JWTManager
	OSManager     *dependencies.OSManager
	DockerManager *dependencies.DockerManager
}

func NewDependencies(configPath string) *Dependencies {

	log.Println("Loading Config")
	config, err := dependencies.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Initializing Logger")
	logger := dependencies.NewLogger(config)

	logger.Infoln("Initializing Validator")
	validator := validator.New()

	logger.Infoln("Connecting to Redis")
	redis := dependencies.NewRedisClient(config)

	logger.Infoln("Connecting to AWS")
	awsManager := dependencies.NewAWSClient(config.Aws.AccessKeyId, config.Aws.AccessKeySecret)

	logger.Infoln("Connecting to JWT")
	jwtManager := dependencies.NewJWTManager(config.Security.JWT.Secret, time.Hour*time.Duration(config.Security.JWT.ExpiresInHour))

	logger.Infoln("Initializing OS Manager")
	osManager := dependencies.NewOsManager()

	logger.Infoln("Initializing Docker Manager")
	dockerManager := dependencies.NewDockerManager()

	logger.Infoln("Connecting to Database")
	databaseClient, err := dependencies.NewDatabaseClient(config)
	if err != nil {
		logger.Fatalln("Failed to connect to database: ", err)
	}

	logger.Infoln("Connecting to Cynxhost Central")
	cynxhostCentral := cynxhostcentral.New(&config.Central)

	logger.Infoln("Dependencies initialized")
	return &Dependencies{
		Config:          config,
		DatabaseClient:  databaseClient,
		Validator:       validator,
		RedisClient:     redis,
		Logger:          logger,
		AWSClient:       awsManager,
		JWTManager:      jwtManager,
		OSManager:       osManager,
		CynxhostCentral: cynxhostCentral,
		DockerManager:   dockerManager,
	}
}
