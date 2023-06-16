package main

import (
	"context"
	"fmt"
	"github.com/bincooo/openai-wapi"
	"io"
	"time"
)

func main() {
	token := "xxx"
	ru := "https://api.pawan.krd/backend-api"
	chat, err := chatgpt.New(token, ru)
	if err != nil {
		panic(err)
	}

	prompt := "hi"
	fmt.Println("You: ", prompt)
	partialResponse, err := chat.Reply(context.Background(), prompt)
	if err != nil {
		panic(err)
	}
	Println(partialResponse)

	prompt = "who are you?"
	fmt.Println("You: ", prompt)
	partialResponse, err = chat.Reply(context.Background(), prompt)
	if err != nil {
		panic(err)
	}
	Println(partialResponse)

	prompt = "what can you do?"
	fmt.Println("You: ", prompt)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()
	partialResponse, err = chat.Reply(ctx, prompt)
	if err != nil {
		panic(err)
	}
	Println(partialResponse)
}

func Println(partialResponse chan chatgpt.PartialResponse) {
	for {
		message, ok := <-partialResponse
		if !ok {
			return
		}

		if message.Error != nil {
			if message.Error == io.EOF {
				return
			}
			panic(message.Error)
		}

		if message.ResponseError != nil {
			fmt.Println(message.ResponseError)
			return
		}

		s := message.Message.Content.Parts[0]
		fmt.Println(s)
		fmt.Println("===============")
	}
}
