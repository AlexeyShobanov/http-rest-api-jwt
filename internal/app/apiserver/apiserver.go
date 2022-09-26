package apiserver

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/ash/http-rest-api/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	file, err := os.OpenFile(config.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		srv.logger.Fatal(err)
	}

	defer file.Close()

	srv.configureLogger(config.LogLevel, file)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(DatabaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
