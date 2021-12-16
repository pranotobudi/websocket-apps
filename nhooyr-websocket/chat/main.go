package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.SetFlags(0)

	err := run()
	if err != nil {
		log.Fatal(err)
	}
	//go run main.go chat.go localhost:8080
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func run() error {
	// go run main.go chat.go localhost:8000
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	cs := newChatServer()
	s := &http.Server{
		Handler:      cs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	// preparing signal from os like Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	// apps will terminate due to two condition: os signal or errc from s.Server(l)
	// apps will wait here until something happen, because no "default" in this select.
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//even in it's official document function example, they don't call cancel(),
	// so, I doubt it is a must.

	return s.Shutdown(ctx)
}
