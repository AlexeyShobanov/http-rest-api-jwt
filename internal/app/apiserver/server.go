// легковесная реализация сервера без запуска http-сервера
package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ash/http-rest-api/internal/app/model"
	"github.com/ash/http-rest-api/internal/app/store"
	pkg "github.com/ash/http-rest-api/pkg/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

//  создаем тип данных для контекста

type ctxKey int8

type server struct {
	logger   *logrus.Logger
	router   *mux.Router
	store    store.Store
	jwtToken *pkg.ConfToken
}

func newServer(store store.Store, ct *pkg.ConfToken) *server {
	s := &server{
		router:   mux.NewRouter(),
		logger:   logrus.New(),
		store:    store,
		jwtToken: ct,
	}

	s.configureRouter()

	return s
}

func (s *server) configureLogger(l string, file *os.File) {
	level, err := logrus.ParseLevel(l)
	if err != nil {
		s.logger.Fatal(err)
	}

	s.logger.SetLevel(level)

	s.logger.Formatter = &logrus.JSONFormatter{}
	s.logger.SetOutput(os.Stdout)
	s.logger.SetOutput(file)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	// вызываем мидлвару которая устанавливает id запроса
	s.router.Use(s.setRequestID)
	// вызываем мидлвару которая логирует запрос
	s.router.Use(s.logRequest)
	// вызываем мидлвару CORS для всех роутов
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// стандартные роуты
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	// создаем подроутер закрытый аутентификацией (получим роут вида /private/...)
	private := s.router.PathPrefix("/private").Subrouter()
	// указываем мидлвару которую хотим применить к этому подроутеру
	private.Use(s.authenticateUser)
	// далее описываем роуты этого подроутера
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
}

// мидлвара для генерации id запроса
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// мидлвара для логирования запросов
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWiter{w, http.StatusOK} //создаем экземпляр кастомного респонс-райтера
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		headerParts := strings.Split(authHeader[0], " ")

		if len(headerParts) != 2 {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		if headerParts[0] != "Bearer" {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		id, err := s.jwtToken.ParseToken(headerParts[1])
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u))) // прокидываем пользователя в контексте, чтобы не повторять его поиск при дальнейших запросах
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		// декодируем в эту структуру входящий payload от пользователя
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		// создаем токен jwt
		token, err := s.jwtToken.GenerateToken(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
