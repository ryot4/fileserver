package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	listenAddress string
	documentRoot  string
}

const version = "0.0.1"

func main() {
	var config Config

	flag.StringVar(&config.listenAddress, "l", "0.0.0.0:8080", "Address to listen on")
	flag.StringVar(&config.documentRoot, "d", ".", "Directory to serve")
	printVersion := flag.Bool("v", false, "Print the version and exit")
	flag.Parse()

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if err := checkDirectory(config.documentRoot); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	runServer(&config)
}

func checkDirectory(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	} else if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

func newServer(config *Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", &LoggingHandler{http.FileServer(http.Dir(config.documentRoot))})
	return &http.Server{Addr: config.listenAddress, Handler: mux}
}

func runServer(config *Config) {
	srv := newServer(config)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		signal.Notify(sigCh, syscall.SIGTERM)

		s := <-sigCh
		switch s {
		case os.Interrupt:
			log.Println("received SIGINT. shutting down")
		case syscall.SIGTERM:
			log.Println("received SIGTERM. shutting down")
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("failed to shutdown the server gracefully:", err.Error())
		}
		close(idleConnsClosed)
	}()

	log.Printf("serving %s on %s\n", config.documentRoot, config.listenAddress)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	<-idleConnsClosed
}
