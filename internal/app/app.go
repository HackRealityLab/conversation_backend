package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"Hackathon/internal/app/grpclient"
	"Hackathon/internal/config"
	"Hackathon/internal/transport/rest"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	serverAddr = ":8000"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	filesCh := make(chan grpclient.FileRequest)

	minioClient := setupMinio(cfg.MinioConfig)

	go func() {
		runRestServer(cfg, minioClient, filesCh)
	}()
	go func() {
		runGrpcClient(cfg.AIConfig, filesCh)
	}()

	<-ctx.Done()
	close(filesCh)
}

func runRestServer(cfg *config.Config, minioClient *minio.Client, filesCh chan<- grpclient.FileRequest) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	conversationHandler := rest.NewConversationHandler(validate, minioClient, cfg.MinioConfig, filesCh)

	mux := newServeMux(conversationHandler)

	server := http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	log.Printf("Run server on %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stop server")
}

func newServeMux(
	conversationHandler *rest.ConversationHandler,
) *http.ServeMux {
	mux := &http.ServeMux{}
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /conversation/text", conversationHandler.LoadConversationText)
	mux.HandleFunc("POST /conversation/file", conversationHandler.LoadConversationFile)
	mux.HandleFunc("GET /conversation/file/{name}", conversationHandler.GetConversationFile)
	mux.HandleFunc("POST /conversation/file/send_ai/{name}", conversationHandler.SendFileToAI)

	return mux
}

func runGrpcClient(aiConfig *config.AIConfig, filesCh <-chan grpclient.FileRequest) {
	cl := grpclient.NewGRPCClient(aiConfig)
	cl.SendFileToAI(filesCh)
}

func setupMinio(minioCfg *config.MinioConfig) *minio.Client {
	// Initialize minio client object.
	minioClient, err := minio.New(minioCfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioCfg.ServerAccessKey, minioCfg.ServerSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	bucketName := minioCfg.ConversationBucket
	location := "us-east-1"

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return minioClient
}
