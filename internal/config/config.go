package config

import "os"

type Config struct {
	Port           string
	WorkerCount    string // int
	JobQueueDepth  string // int
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MongoURI       string
	MongoDB        string
}

func Load() Config {
	return Config{
		Port:           os.Getenv("PORT"),
		WorkerCount:    os.Getenv("WORKER_COUNT"),
		JobQueueDepth:  os.Getenv("JOB_QUEUE_DEPTH"),
		MinioEndpoint:  os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:    os.Getenv("MINIO_BUCKET"),
		MongoURI:       os.Getenv("MONGO_URI"),
		MongoDB:        os.Getenv("MONGO_DB"),
	}
}
