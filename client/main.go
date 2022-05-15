package main

import (
	"context"
	"io"
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

	//DoUnary(client)
	//DoClientStreaming(client)
	//DoServerStreaming(client)
	DoBidirectionalStreaming(client)

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

func DoServerStreaming(client testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "t1",
	}

	stream, err := client.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetStudentsPerTest: %v", err)
	}

	for {
		student, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("student: %v", student)
	}
}

func DoBidirectionalStreaming(client testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "42",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := client.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("error while calling TakeTest: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}
			log.Printf("response from TakeTest: %v", res)
		}
		close(waitChannel)
	}()
	<-waitChannel
}
