package apiserver

import (
	"github.com/evd1ser/go-homework-finish/internal/app/middleware"
	"github.com/evd1ser/go-homework-finish/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prefix string = ""
)

// type for APIServer object for instancing server
type APIServer struct {
	//Unexported field
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

//APIServer constructor
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start http server and connection to db and logger confs
func (api *APIServer) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}
	api.logger.Info("starting api server at port :", api.config.BindAddr)
	api.configureRouter()
	if err := api.configureStore(); err != nil {
		return err
	}
	return http.ListenAndServe(api.config.BindAddr, api.router)
}

//func for configureate logger, should be unexported
func (api *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(api.config.LogLevel)
	if err != nil {
		return nil
	}
	api.logger.SetLevel(level)

	return nil
}

//func for configure Router
func (api *APIServer) configureRouter() {
	//Было до JWT
	//Теперь требует наличия JWT
	//api.router.Handle(prefix+"/articles"+"/{id}", middleware.JwtMiddleware.Handler(
	//	http.HandlerFunc(s.GetArticleById),
	//)).Methods("GET")
	//

	api.router.HandleFunc(prefix+"/register", api.PostUserRegister).Methods("POST")
	api.router.HandleFunc(prefix+"/auth", api.PostToAuth).Methods("POST")

	//приватные методы
	apiRouter := api.router.PathPrefix(prefix + "/").Subrouter()
	apiRouter.Use(middleware.JwtMiddleware)

	//stock,
	apiRouter.HandleFunc(prefix+"/stock", api.GetStock).Methods("GET")
	apiRouter.HandleFunc(prefix+"/auto/{mark}", api.PostAutoCreate).Methods("POST")
	apiRouter.HandleFunc(prefix+"/auto/{mark}", api.GetAuto).Methods("GET")
	apiRouter.HandleFunc(prefix+"/auto/{mark}", api.PutAutoUpdate).Methods("PUT")
	apiRouter.HandleFunc(prefix+"/auto/{mark}", api.DeleteAuto).Methods("DELETE")
	//auto/<string:mark>
}

//configureStore method
func (api *APIServer) configureStore() error {
	st := store.New(api.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	api.store = st
	return nil
}
