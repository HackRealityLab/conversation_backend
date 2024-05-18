package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	aiHostKey                  = "AI_HOST"
	aiPortKey                  = "AI_PORT"
	minioServerAccessKey       = "MINIO_SERVER_ACCESS_KEY"
	minioServerSecretKey       = "MINIO_SERVER_SECRET_KEY"
	minioEndpointKey           = "MINIO_ENDPOINT"
	minioConversationBucketKey = "MINIO_CONVERSATION_BUCKET"
)

type Config struct {
	AIConfig    *AIConfig
	MinioConfig *MinioConfig
}

type AIConfig struct {
	Host string
	Port string
}

type MinioConfig struct {
	EndPoint           string
	ServerAccessKey    string
	ServerSecretKey    string
	ConversationBucket string
}

func LoadConfig() (*Config, error) {
	aiConfig, err := loadAIConfig()
	if err != nil {
		return nil, err
	}

	minioConfig, err := loadMinioConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		AIConfig:    aiConfig,
		MinioConfig: minioConfig,
	}, nil
}

func loadAIConfig() (*AIConfig, error) {
	host, err := getEnvString(aiHostKey)
	if err != nil {
		return nil, err
	}

	port, err := getEnvString(aiPortKey)
	if err != nil {
		return nil, err
	}

	return &AIConfig{
		Host: host,
		Port: port,
	}, nil
}

func loadMinioConfig() (*MinioConfig, error) {
	serverAccessKey, err := getEnvString(minioServerAccessKey)
	if err != nil {
		return nil, err
	}

	serverSecretKey, err := getEnvString(minioServerSecretKey)
	if err != nil {
		return nil, err
	}

	endpoint, err := getEnvString(minioEndpointKey)
	if err != nil {
		return nil, err
	}

	conversationBucket, err := getEnvString(minioConversationBucketKey)
	if err != nil {
		return nil, err
	}

	return &MinioConfig{
		EndPoint:           endpoint,
		ServerAccessKey:    serverAccessKey,
		ServerSecretKey:    serverSecretKey,
		ConversationBucket: conversationBucket,
	}, nil
}

func getEnvString(name string) (string, error) {
	env := os.Getenv(name)
	if env == "" {
		msg := fmt.Sprintf("no env with name: %s", env)
		return "", errors.New(msg)
	}

	return env, nil
}
