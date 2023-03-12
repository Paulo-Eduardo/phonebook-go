package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Paulo-Eduardo/phonebook-go/handlers"
)

func main() {
	port := flag.String("port", "9000", "Port number to expose the API")
	timeout := flag.Int("timeout", 2, "Timeout in seconds")

	flag.Parse();

	println(*port)

	l := log.New(os.Stdout, "phonebook-api", log.LstdFlags)

	ph := handlers.NewPhonebook(l)

	sm := http.NewServeMux()

	sm.Handle("/list", ph.ListContact())
	sm.Handle("/add", ph.AddContact())
	sm.Handle("/delete/", ph.DeleteContact())
	sm.Handle("/update/", ph.UpdateContact())
	sm.Handle("/find/", ph.FindContact())
	sm.Handle("/find-by-name/", ph.FindContactByName())

	s := http.Server{
		Addr: fmt.Sprintf(":%s", *port),
		Handler: sm,
		ReadTimeout: time.Duration(*timeout) * time.Second,
		WriteTimeout: time.Duration(*timeout) * time.Second,
		ErrorLog: l,
	}

	go func() {
		l.Printf("Starting server on port %s\n", *port)

		if err := s.ListenAndServe(); err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30* time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	 }
}