package main

import (
	"context"
	"fmt"
	"github.com/bincooo/openai-wapi"
	"io"
	"time"
)

func main() {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiJiaW5jbzAwMDAwMDAyQG91dGxvb2suY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWV9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItRHRsaU9TWjIyM1duQ0JOSmRHT2FIUzJsIn0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJhdXRoMHw2M2ExMmRhZTUzMDRmY2NlMmE0MGU0NDkiLCJhdWQiOlsiaHR0cHM6Ly9hcGkub3BlbmFpLmNvbS92MSIsImh0dHBzOi8vb3BlbmFpLm9wZW5haS5hdXRoMGFwcC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjg2NjQyNDgzLCJleHAiOjE2ODc4NTIwODMsImF6cCI6InBkbExJWDJZNzJNSWwycmhMaFRFOVZWOWJOOTA1a0JoIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBtb2RlbC5yZWFkIG1vZGVsLnJlcXVlc3Qgb3JnYW5pemF0aW9uLnJlYWQgb2ZmbGluZV9hY2Nlc3MifQ.H1d87aGyNfuDBfT9-RVoUEbzsDYmYtuguIk9zxaAfOvdGmgGFmEf2O0P6AsO5iIaoaMEIoTUgB_BF1V5AKZAhcbtl4dgl3arZjSGv_QTWiFxbJuMkClnpr7VAIGgENZkqry6VUnuRvEMdPBSrWYvaoeSGmIh7_erulkbDwjDTRDAn96RLdx4l6p-2WlP5PJRTRhQWiE1vptGdf7yKzeIpQRuhzgX7LzdXfTzcSPHmTj53v_7UCRxcRZ2xA-UVtyXLUG0_wa45HrxhqNv9cElnQdhf0oBjeUdEWk3aCvMPHxBKgJQADxDNesWdlTsxyDMleqw_AQ8r5mwVOHIQ9RMAA"
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
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
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
