package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
	dev  bool
}

type application struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg config

	// Default values for production
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.BoolVar(&cfg.dev, "dev", false, "Development mode")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.port),
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
