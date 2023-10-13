package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"connectrpc.com/connect"
	apiv1 "github.com/pirosiki197/emoine/pkg/proto/api/v1"
	"github.com/pirosiki197/emoine/pkg/proto/api/v1/apiv1connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	scanner    *bufio.Scanner
	httpClient *http.Client
)

func main() {
	scanner = bufio.NewScanner(os.Stdin)
	httpClient = http.DefaultClient
	client := apiv1connect.NewAPIServiceClient(httpClient, "http://emoine.trapti.tech")

	fmt.Println("connected to server")

	fmt.Println("1: create event")
	fmt.Println("2: get events")
	fmt.Println("3: send comment")
	fmt.Println("4: get comments")
	fmt.Println("5: connect to stream")
	fmt.Println("6: exit")

	var n int
	for {
		fmt.Println("=====================================")
		fmt.Println("enter number:")
		fmt.Scan(&n)

		switch n {
		case 1:
			createEvent(client)
		case 2:
			getEvents(client)
		case 3:
			sendComment(client)
		case 4:
			getComments(client)
		case 5:
			connectToStream(client)
		case 6:
			fmt.Println("bye")
			return
		default:
			fmt.Println("invalid number")
		}
	}
}

func createEvent(c apiv1connect.APIServiceClient) {
	var title string
	fmt.Println("enter title:")
	scanner.Scan()
	title = scanner.Text()

	res, err := c.CreateEvent(context.Background(), connect.NewRequest(&apiv1.CreateEventRequest{
		Title:   title,
		StartAt: timestamppb.Now(),
	}))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("created event")
	fmt.Printf("id: %s\n", res.Msg.Id)
}

func getEvents(c apiv1connect.APIServiceClient) {
	res, err := c.GetEvents(context.Background(), connect.NewRequest(&apiv1.GetEventsRequest{}))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("total: %d\n", len(res.Msg.Events))
	fmt.Println("events:")
	for _, e := range res.Msg.Events {
		fmt.Printf("%s: %s\n", e.Id, e.Title)
	}
}

func sendComment(c apiv1connect.APIServiceClient) {
	var eventID string
	fmt.Println("enter event id:")
	fmt.Scan(&eventID)

	var text string
	fmt.Println("enter text:")
	scanner.Scan()
	text = scanner.Text()

	res, err := c.SendComment(context.Background(), connect.NewRequest(&apiv1.SendCommentRequest{
		EventId: eventID,
		UserId:  "guest",
		Text:    text,
	}))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("sent comment")
	fmt.Printf("id: %s\n", res.Msg.Id)
}

func getComments(c apiv1connect.APIServiceClient) {
	var eventID string
	fmt.Println("enter event id:")
	fmt.Scan(&eventID)

	res, err := c.GetComments(context.Background(), connect.NewRequest(&apiv1.GetCommentsRequest{
		EventId: eventID,
	}))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("comments:")
	for _, c := range res.Msg.Comments {
		fmt.Printf("%s: %s\n", c.Id, c.Text)
	}
}

func connectToStream(c apiv1connect.APIServiceClient) {
	var eventID string
	fmt.Println("enter event id:")
	fmt.Scan(&eventID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := c.ConnectToStream(ctx, connect.NewRequest(&apiv1.ConnectToStreamRequest{
		EventId: eventID,
	}))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("connected to stream")

	go func() {
		fmt.Println("to exit from stream, enter 'exit'")
		for {
			scanner.Scan()
			if scanner.Text() == "exit" {
				fmt.Println("exiting from stream...")
				cancel()
				return
			}
		}
	}()

	for stream.Receive() {
		msg := stream.Msg()

		if e := msg.GetEvent(); e != nil {
			fmt.Printf("event: %s\n", e.Title)
		}
		if c := msg.GetComment(); c != nil {
			fmt.Printf("comment: %s\n", c.Text)
		}
	}
	if err := stream.Err(); err != nil {
		log.Println(err)
		return
	}
}
