package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ash/http-rest-api/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)

	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	srv := newServer(sqlstore.New(db))
	srv.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
}
