package auth

import (
	"context"
	"log"
)

type AuthServer struct {

}

func (s *AuthServer) Test(ctx context.Context, message *Message) (*Message, error) {
	log.Printf("Received %s", message.Body)
	return &Message{Body: "return message"}, nil
}

