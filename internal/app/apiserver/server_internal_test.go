package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ash/http-rest-api/internal/app/model"
	"github.com/ash/http-rest-api/internal/app/store/sqlstore"
	pkg "github.com/ash/http-rest-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestServer_AuthenticateUser(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	store := sqlstore.New(db)

	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		userId       int
		expectedCode int
	}{
		{
			name:         "authenticated",
			userId:       u.ID,
			expectedCode: http.StatusOK,
		},
		// {
		// 	name:         "not authenticated",
		// 	userId:       u.ID + 1,
		// 	expectedCode: http.StatusUnauthorized,
		// },
	}

	token := pkg.Init("/Users/alexey_sh/go/src/learn-project/http-rest-api/pkg/jwt/config.toml")
	s := newServer(store, token) // создаем сервер

	// sc := securecookie.New(secretKey, nil)
	// создаем фейковый хендлер next для authenticateUser с функцией которая просто записывает в хедер статус 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := s.authenticateUser(handler)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			// создаем куку в виде строки
			jwtToken, _ := s.jwtToken.GenerateToken(tc.userId)
			// добавляем созданную куку в запрос в виде заголовка Cookie
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
			mw.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUsersCreate(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	token := pkg.Init("/Users/alexey_sh/go/src/learn-project/http-rest-api/pkg/jwt/config.toml")

	srv := newServer(sqlstore.New(db), token)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}

func TestServer_HandleSessionsCreate(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")
	store := sqlstore.New(db)
	token := pkg.Init("/Users/alexey_sh/go/src/learn-project/http-rest-api/pkg/jwt/config.toml")
	srv := newServer(store, token)

	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

		})
	}
}
