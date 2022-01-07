package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"banisaeid.com/letsgo/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "Address to listen to")
	dsn := flag.String("dsn", "web:web@/snippetbox?parseTime=true", "MySQL connection string (DSN)")

	// Secret for http sessions. Must be 32 bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret ket")

	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize the templates cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  app.errorLog,
		TLSConfig: tlsConfig,
	}

	app.infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
