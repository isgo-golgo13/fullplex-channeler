package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"fullplex-channeler/svckit"
)

type Server struct {
	port string
}

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func NewServer(opts ...Option) *Server {
	s := &Server{port: "8080"} // Default port is 8080
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func main() {
	// Initialize the server with a configurable port using functional options.
	server := NewServer(WithPort("8080"))

	ln, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer ln.Close()

	// Channel to handle OS signals for graceful shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Context for managing connection handling and graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown.
	go func() {
		<-sigs
		fmt.Println("Received shutdown signal, shutting down...")
		cancel()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("Server shutdown complete")
				return
			default:
				log.Println("Failed to accept connection:", err)
			}
			continue
		}

		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	channel := svckit.NewFullPlexChanneler(conn)

	// Receiving data into multiple writers using RecvAll
	var requestBuffer1 bytes.Buffer
	var requestBuffer2 bytes.Buffer
	writers := []io.Writer{&requestBuffer1, &requestBuffer2}
	if _, err := channel.RecvAll(ctx, writers, 1024); err != nil {
		log.Println("RecvAll error:", err)
		return
	}
	fmt.Println("Received from client (Buffer 1):", requestBuffer1.String())
	fmt.Println("Received from client (Buffer 2):", requestBuffer2.String())

	// Sending a response back to the client
	response := strings.NewReader("Message from Server!")
	if _, err := channel.Send(ctx, response, int64(response.Len())); err != nil {
		log.Println("Send error:", err)
	}

	channel.Close(ctx)
}
