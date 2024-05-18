package grpclient

import (
	"context"
	"fmt"
	"io"
	"log"

	"Hackathon/internal/config"
	conversation "Hackathon/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppClient struct {
	client conversation.ConversationClient
}

type FileRequest struct {
	UUID      string
	FileName  string
	FileBytes []byte
}

func NewGRPCClient(aiConfig *config.AIConfig) *AppClient {
	target := fmt.Sprintf("%s:%s", aiConfig.Host, aiConfig.Port)

	transportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(target, transportOpt)
	if err != nil {
		log.Fatal(err)
	}

	client := conversation.NewConversationClient(conn)
	return &AppClient{
		client: client,
	}
}

func (c *AppClient) SendFileToAI(filesCh <-chan FileRequest) {
	streamGreater, err := c.client.AnalyzeAudio(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go asyncClientBidirectionalRPC(streamGreater, doneCh)

	for request := range filesCh {
		err = streamGreater.Send(&conversation.ConversationRequest{
			ConversationID: request.UUID,
			FileName:       request.FileName,
			File:           request.FileBytes,
		})

		if err != nil {
			log.Println(err)
		}
	}

	err = streamGreater.CloseSend()
	if err != nil {
		log.Fatal()
	}

	<-doneCh
}

func (c *AppClient) AnalyzeTest() {
	streamGreater, err := c.client.AnalyzeAudio(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go asyncClientBidirectionalRPC(streamGreater, doneCh)

	err = streamGreater.Send(&conversation.ConversationRequest{
		ConversationID: uuid.New().String(),
		File:           []byte{1, 2, 3, 4, 5, 6},
	})
	if err != nil {
		log.Println(err)
	}

	err = streamGreater.CloseSend()
	if err != nil {
		log.Fatal()
	}

	<-doneCh
	log.Println("Stop bidi streaming")
}

func asyncClientBidirectionalRPC(
	streamGreater conversation.Conversation_AnalyzeAudioClient,
	doneCh chan struct{},
) {
	for {
		reply, err := streamGreater.Recv()
		if err == io.EOF {
			break
		}

		log.Printf("Received reply: %s\n", reply.Text)
	}

	doneCh <- struct{}{}
}
