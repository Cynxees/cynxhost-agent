package controller

import (
	"cynxhostagent/internal/app"
	"cynxhostagent/internal/controller/usercontroller"
	"cynxhostagent/internal/middleware"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.elastic.co/apm/module/apmhttp/v2"
)

type HttpServer struct {
	*http.Server
}

func NewHttpServer(app *app.App) (*HttpServer, error) {

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                            // replace with your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, // allowed methods
		AllowedHeaders: []string{"Content-Type"},                 // allowed headers
	})

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	routerPath := app.Dependencies.Config.Router.Default
	debug := app.Dependencies.Config.App.Debug

	handleRouterFunc := func(path string, handler middleware.HandlerFuncWithHelper, requireAuth bool) *mux.Route {
		wrappedHandler := middleware.WrapHandler(handler, debug)

		if requireAuth && !debug {
			wrappedHandler = middleware.AuthMiddleware(app.Dependencies.JWTManager, wrappedHandler, debug)
		}
		return r.HandleFunc(routerPath+path, wrappedHandler).Methods("POST", "GET")
	}

	userController := usercontroller.New(app.Usecases.UserUseCase, app.Dependencies.Validator)

	// User
	handleRouterFunc("user/bypass-login", userController.BypassLoginUser, false)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			app.Dependencies.Logger.Errorf("Failed to write response: %v", err)
		}
	})

	corsHandler := c.Handler(r)

	address := app.Dependencies.Config.App.Address + ":" + strconv.Itoa(app.Dependencies.Config.App.Port)
	app.Dependencies.Logger.Infof("Starting http server on %s\n", address)

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      apmhttp.Wrap(corsHandler),
	}

	return &HttpServer{srv}, nil
}

func (s *HttpServer) Start() error {
	return s.ListenAndServe()
}

func (s *HttpServer) Stop() error {
	return errors.New("http stop not implemented")
}
