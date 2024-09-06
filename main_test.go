package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find an available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Start the server in a goroutine with the available port
	serverExited := make(chan struct{})
	go func() {
		run(port)
		close(serverExited)
	}()

	// Give the server a moment to start
	time.Sleep(time.Second)

	// Create a client and make a request to ensure the server is running
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect to the server: %v", err)
	}
	resp.Body.Close()

	// Simulate a SIGINT signal to trigger graceful shutdown
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find the current process: %v", err)
	}
	err = p.Signal(syscall.SIGINT)
	if err != nil {
		t.Fatalf("Failed to send SIGINT signal: %v", err)
	}

	// Wait for the server to exit
	select {
	case <-serverExited:
		// Server exited successfully
	case <-time.After(5 * time.Second):
		t.Fatal("Server did not shut down within the expected time")
	}

	// Try to make another request, which should fail if the server is shut down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d", port), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatalf("Expected an error due to server shutdown, but got none")
	}

	// Check if the error is related to connection refused or timeout
	if !isConnectionRefusedOrTimeout(err) {
		t.Fatalf("Expected connection refused or timeout error, but got: %s", err.Error())
	}
}

func isConnectionRefusedOrTimeout(err error) bool {
	return strings.Contains(err.Error(), "connect: connection refused") || strings.Contains(err.Error(), "context deadline exceeded")
}
