package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func run(ctx context.Context) error {
	// listener
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("failed to load `PORT` env: %v", err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", port, err)
	}
	// mux
	mux, err := NewMux(ctx)
	if err != nil {
		return err
	}
	// server
	s := NewServer(l, mux)
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("starting server at: %v", url)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("faled to terminate server: %v", err)
		os.Exit(1)
	}
}
