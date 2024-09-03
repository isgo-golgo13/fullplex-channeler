package main

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"fullplex-channeler/svckit"
)

func main() {
	conn, err := net.Dial("tcp", "server:8080")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	ctx := context.Background()
	channel := svckit.NewFullPlexChanneler(conn)

	// 1. Sending a simple message to the server using Send()
	data := strings.NewReader("Message from Client using Send!")
	if _, err := channel.Send(ctx, data, int64(data.Len())); err != nil {
		log.Fatal("Send error:", err)
	}

	// 2. Sending a message with a timeout using SendTimeout()
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := channel.SendTimeout(timeoutCtx, strings.NewReader("Timed Message from Client"), 1024, 5*time.Second); err != nil {
		log.Fatal("SendTimeout error:", err)
	}

	// 3. Sending multiple messages using SendAll()
	messages := []io.Reader{
		strings.NewReader("First part of a message from SendAll."),
		strings.NewReader("Second part of a message from SendAll."),
		strings.NewReader("Third part of a message from SendAll."),
	}
	if _, err := channel.SendAll(ctx, messages, 1024); err != nil {
		log.Fatal("SendAll error:", err)
	}

	// Close the channel
	channel.Close(ctx)
}
