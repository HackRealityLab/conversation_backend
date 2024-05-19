package grpclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"Hackathon/internal/config"
	"Hackathon/internal/service"
	conversation "Hackathon/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AppClient struct {
	client  conversation.ConversationClient
	service service.ConversationService
}

func NewGRPCClient(
	aiConfig *config.AIConfig,
	service service.ConversationService,
) *AppClient {
	target := fmt.Sprintf("%s:%s", aiConfig.Host, aiConfig.Port)

	transportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(target, transportOpt)
	if err != nil {
		log.Fatal(err)
	}

	client := conversation.NewConversationClient(conn)
	return &AppClient{
		client:  client,
		service: service,
	}
}

func (c *AppClient) SendFileToAI(filesCh <-chan conversation.ConversationRequest) {
	streamRequest, err := c.client.AnalyzeAudio(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go c.asyncClientBidirectionalRPC(streamRequest, doneCh)

	for request := range filesCh {
		log.Printf("Send file with name: %s", request.FileName)
		err = streamRequest.Send(&request)

		if err != nil {
			log.Printf("Got error: %s", err.Error())
		} else {
			log.Printf("No error: %v", err)
		}
	}

	err = streamRequest.CloseSend()
	if err != nil {
		log.Fatal()
	}

	<-doneCh
}

func (c *AppClient) asyncClientBidirectionalRPC(
	streamConversation conversation.Conversation_AnalyzeAudioClient,
	doneCh chan struct{},
) {
	for {
		log.Println("start recv")
		reply, err := streamConversation.Recv()
		if err == io.EOF {
			break
		}
		if reply == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		err = c.service.InsertAdditionRecordInfo(
			int(reply.ConversationID),
			reply.Text,
			int(reply.GoodPercent),
			int(reply.BadPercent),
		)

		log.Printf("Err while insert additional info: %s", err.Error())

		log.Printf("Received reply: %+v\n", reply)
	}

	doneCh <- struct{}{}
}
