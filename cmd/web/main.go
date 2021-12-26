package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "Address to listen to")
	flag.Parse()

	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	app.infoLog.Printf("Starting server on %s\n", *addr)
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLog,
	}

	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
