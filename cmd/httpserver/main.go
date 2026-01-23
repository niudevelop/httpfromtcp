package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/niudevelop/httpfromtcp/internal/response"
	"github.com/niudevelop/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := func(req *server.Request, w io.Writer) *server.HandlerError {
		switch req.Target {
		case "/yourproblem":
			return &server.HandlerError{
				Status:  response.StatusCode400,
				Message: "Your problem is not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				Status:  response.StatusCode500,
				Message: "Woopsie, my bad\n",
			}
		default:
			_, _ = io.WriteString(w, "All good, frfr\n")
			return nil
		}
	}
	server, err := server.Serve(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
