package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"go-snp/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://root:secret@localhost/goSnpDB?sslmode=disable", "Postgres data source dsn")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, nativeDB, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()
	defer nativeDB.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(nativeDB)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, *sql.DB, error) {
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("connect to postgres")

	nativeDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, err
	}

	err = nativeDB.PingContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("native connect to postgres")

	return conn, nativeDB, nil
}
