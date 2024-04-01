package network

import (
	"fmt"
	"errors"
	"net"
	"os"
	"os/signal"
)

type server struct {
	listener net.Listener
}

func NewServer() (server, error) {
	server := server{}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return server, err
	}
	server.listener = listener

	return server, nil
}

func (s server) Start() {
	s.handleSignals()
	fmt.Println("server.start(): starting server")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Printf("server.start(): listener closed.")
				os.Exit(0)
			}

			fmt.Printf("server.start() error: %s", err)
		}

		go process(conn)
	}
}

func (s server) Stop() {
	fmt.Println("server.close(): closing server")
	if err := s.listener.Close(); err != nil {
		fmt.Printf("error closing server: %s", err)
		os.Exit(1)
	}
}

func (s server) handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		s.Stop()
		os.Exit(0)
	}()
}