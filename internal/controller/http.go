package controller

import (
	"cynxhostagent/internal/app"
	"cynxhostagent/internal/controller/persistentnodecontroller"
	"cynxhostagent/internal/controller/usercontroller"
	"cynxhostagent/internal/middleware"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"go.elastic.co/apm/module/apmhttp/v2"
)

type HttpServer struct {
	http *http.Server
	ws   *http.Server
}

func NewHttpServer(app *app.App) (*HttpServer, error) {

	c := cors.New(cors.Options{
		AllowedOrigins: app.Dependencies.Config.Security.CORS.Origins,       // replace with your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // allowed methods
		AllowedHeaders: []string{"Content-Type"},                            // allowed headers
	})

	r := mux.NewRouter()
	wsRouter := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	routerPath := app.Dependencies.Config.Router.Default
	debug := app.Dependencies.Config.App.Debug

	handleDefaultRouterFunc := func(path string, handler middleware.HandlerFuncWithHelper, requireAuth bool) *mux.Route {
		wrappedHandler := middleware.WrapHandler(handler, debug)

		if requireAuth && !debug {
			wrappedHandler = middleware.AuthMiddleware(app.Dependencies.JWTManager, wrappedHandler, debug)
		}
		return r.HandleFunc(routerPath+path, wrappedHandler).Methods("POST", "GET")
	}

	handleBasicRouterFunc := func(path string, handler func(http.ResponseWriter, *http.Request), requireAuth bool) *mux.Route {

		if requireAuth && !debug {
			handler = middleware.AuthMiddleware(app.Dependencies.JWTManager, handler, debug)
		}

		return r.HandleFunc(routerPath+path, handler).Methods("POST", "GET")
	}

	handleWebsocketFunc := func(path string, handler func(w http.ResponseWriter, r *http.Request, conn *websocket.Conn)) *mux.Route {
		fmt.Println("Registering websocket handler for", routerPath+path)
		return wsRouter.HandleFunc("/cynxws"+routerPath+path, func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				app.Dependencies.Logger.Errorf("Failed to upgrade connection: %v", err)
				return
			}
			defer conn.Close()
			handler(w, r, conn)
		})
	}

	userController := usercontroller.New(app.Usecases.UserUseCase, app.Dependencies.Validator)
	persistentNodeController := persistentnodecontroller.New(app.Usecases.PersistentNodeUseCase, app.Dependencies.Validator, app.Dependencies.Config)

	// User
	handleDefaultRouterFunc("user/bypass-login", userController.BypassLoginUser, false)

	// Persistent Node
	handleDefaultRouterFunc("persistent-node/run-template-script", persistentNodeController.RunPersistentNodeTemplateScript, false)

	// Dashboard
	handleDefaultRouterFunc("persistent-node/dashboard/container-stats", persistentNodeController.GetNodeContainerStats, false)
	handleDefaultRouterFunc("persistent-node/dashboard/send-single-docker-command", persistentNodeController.SendSingleDockerCommand, false)

	// Console
	handleDefaultRouterFunc("persistent-node/dashboard/console/create-session", persistentNodeController.CreateSession, false)
	handleDefaultRouterFunc("persistent-node/dashboard/console/send-command", persistentNodeController.SendCommand, false)

	// Files
	handleBasicRouterFunc("persistent-node/dashboard/files/download-file", persistentNodeController.DownloadFile, false)
	handleDefaultRouterFunc("persistent-node/dashboard/files/upload-file", persistentNodeController.UploadFile, false)
	handleDefaultRouterFunc("persistent-node/dashboard/files/remove-file", persistentNodeController.RemoveFile, false)
	handleDefaultRouterFunc("persistent-node/dashboard/files/list-directory", persistentNodeController.ListDirectory, false)
	handleDefaultRouterFunc("persistent-node/dashboard/files/create-directory", persistentNodeController.CreateDirectory, false)
	handleDefaultRouterFunc("persistent-node/dashboard/files/remove-directory", persistentNodeController.RemoveDirectory, false)

	// Images
	handleDefaultRouterFunc("persistent-node/dashboard/images/push", persistentNodeController.PushImage, false)
	handleDefaultRouterFunc("persistent-node/dashboard/images/list", persistentNodeController.ListImages, false)

	// Server Properties
	// handleDefaultRouterFunc("persistent-node/dashboard/properties/get", persistentNodeController.GetServerProperties, false)
	// handleDefaultRouterFunc("persistent-node/dashboard/properties/set", persistentNodeController.SetServerProperties, false)

	// Websocket
	handleWebsocketFunc("persistent-node/logs", persistentNodeController.GetPersistentNodeRealTimeLogs)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			app.Dependencies.Logger.Errorf("Failed to write response: %v", err)
		}
	})

	corsHandler := c.Handler(r)

	address := app.Dependencies.Config.App.Address + ":" + strconv.Itoa(app.Dependencies.Config.App.Port)
	app.Dependencies.Logger.Infof("Starting http server on %s", address)

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      apmhttp.Wrap(corsHandler),
	}

	wsSrv := &http.Server{
		Addr:    app.Dependencies.Config.App.Address + ":" + strconv.Itoa(app.Dependencies.Config.App.WebsocketPort),
		Handler: wsRouter,
	}

	return &HttpServer{
		http: srv,
		ws:   wsSrv,
	}, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Add origin check logic if needed
		return true
	},
}

func (s *HttpServer) Start() error {

	go s.ws.ListenAndServe()
	return s.http.ListenAndServe()
}

func (s *HttpServer) Stop() error {
	return errors.New("http stop not implemented")
}
