package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"Hackathon/internal/app/grpclient"
	"Hackathon/internal/config"
	"Hackathon/internal/repository"
	"Hackathon/internal/service"
	"Hackathon/internal/transport/rest"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	serverAddr = ":8000"
)

func Run() {
	dbConnStr := "postgresql://user:user@postgres:5432/conversations_db?sslmode=disable"
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	runMigrations(dbConnStr)
	filesCh := make(chan grpclient.FileRequest)
	minioClient := setupMinio(cfg.MinioConfig)

	go func() {
		runRestServer(cfg, minioClient, dbConnStr, filesCh)
	}()

	go func() {
		runGrpcClient(cfg.AIConfig, filesCh)
	}()

	<-ctx.Done()
	close(filesCh)
}

func runRestServer(cfg *config.Config, minioClient *minio.Client, dbConnStr string, filesCh chan<- grpclient.FileRequest) {
	dbConn, err := pgx.Connect(context.Background(), dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())

	validate := validator.New(validator.WithRequiredStructEnabled())
	repo := repository.NewConversationRepo(dbConn)
	converterService := service.NewConversationService(repo)

	conversationHandler := rest.NewConversationHandler(validate, minioClient, cfg.MinioConfig, filesCh, converterService)

	mux := newServeMux(conversationHandler)

	server := http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	log.Printf("Run server on %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stop server")
}

func runMigrations(dbConnStr string) {
	log.Printf("Run migrations on %s\n", dbConnStr)
	m, err := migrate.New("file://migrations", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migrate no change")
	} else if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrate ran successfully")
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

	mux.HandleFunc("GET /conversation/records", conversationHandler.GetRecords)
	mux.HandleFunc("GET /conversation/records/{id}", conversationHandler.GetRecord)

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
