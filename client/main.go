package main

import (
	"context"
	"log"
	"time"

	"github.com/acaldo/grpc/testpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5070", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := testpb.NewTestServiceClient(conn)

	DoUnary(client)
	DoClientStreaming(client)

}

func DoUnary(client testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}
	res, err := client.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest: %v", err)
	}
	log.Printf("response from GetTest: %v", res)
}

func DoClientStreaming(client testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q1",
			Answer:   "Cyan",
			Question: "Color asociado a Golang",
			TestId:   "t1",
		},
		{
			Id:       "q2",
			Answer:   "Google",
			Question: "Empresa que desarrollo Go",
			TestId:   "t1",
		},
		{
			Id:       "q3",
			Answer:   "Golang",
			Question: "Lenguaje de programacion",
			TestId:   "t1",
		},
	}
	stream, err := client.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("error while calling SetQuestions: %v", err)
	}
	for _, question := range questions {
		log.Println("sending question: ", question)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}
	log.Printf("response from SetQuestions: %v", res)
}
