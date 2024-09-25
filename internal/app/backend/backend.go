package backend

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/bconfig"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/controllers"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/mclient"
	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/trsclient"
	"github.com/julienschmidt/httprouter"
)

const (
	apiV1Prefix             string        = "/api/v1"
	ServerReadHeaderTimeout time.Duration = 10 * time.Second
)

type Backend struct {
	srv    *http.Server
	config *bconfig.BackendConfig
}

type Option func(options *Backend) error

// WithConfig инициализирует конфиг
func WithConfig() Option {
	return func(options *Backend) error {
		config, err := bconfig.LoadBackendConfig()
		if err != nil {
			return err
		}

		options.config = config

		return nil
	}
}

// WithLogger инициализирует логгер
func WithLogger() Option {
	return func(options *Backend) error {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		slog.SetDefault(logger)

		return nil
	}
}

// New создает новый объект Backend
func New(opts ...Option) *Backend {
	backend := &Backend{}

	for _, opt := range opts {
		err := opt(backend)
		if err != nil {
			slog.Error("ошибка инициализации backend", "error", err)
			os.Exit(1)
		}
	}

	controller := &controllers.Controller{
		TRSParserClient: &trsclient.MockTRSClient{},
		ModelClient:     &mclient.MockModelCient{},
	}

	router := initRoutes(controller)

	backend.srv = &http.Server{
		Handler:           router,
		Addr:              ":" + backend.config.Port,
		ReadHeaderTimeout: ServerReadHeaderTimeout,
	}

	return backend
}

// initRoutes производит инициализацию ручек сервиса
func initRoutes(controller *controllers.Controller) *httprouter.Router {
	router := httprouter.New()

	router.GET(apiV1Prefix+"/trs", controller.TRSCheck)
	router.GET(apiV1Prefix+"/knowledge_base", controller.KnowledgeBase)

	return router
}

func (b *Backend) Run() {
	slog.Info("backend has successfully started at port " + b.config.Port)

	if err := b.srv.ListenAndServe(); err != nil {
		slog.Error("server was crushed", "error", err)
	}
}
