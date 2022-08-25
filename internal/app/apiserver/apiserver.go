// api server
package apiserver

import (
	"io"
	"net/http"

	"github.com/ash/http-rest-api/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// функция New возвращает сконфигурированный интсатнс структуры APIServer
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// метод запускает http сервер, подключение к БД и тд
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting api server")

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

// конфигурируем логер
func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

// опсания роутинга
func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

// конфигурируем Store
func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

// обработчик для /hello
func (s *APIServer) handleHello() http.HandlerFunc {
	// здесь можно определить переменные или типы конкретно для этого обработчика, поэтому это функция возвращает интерфейс, а не сразу ту ф-цию, котороая в return

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello!!!")
	}
}
