package main

import (
	"fmt"
	"log"
	"net"
	library "personal-learning/go-library/pkg/librarylib"
)

const (
	port = ":50051" // PORT on which the GRPC server will listen.
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			fmt.Printf("Failed to close %s %s: %v", "tcp", port, err)
		}
	}()

	library.InitializeLibrary()
	library.StartLibraryServer(lis)
}
