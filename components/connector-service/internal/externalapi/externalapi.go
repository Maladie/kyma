package externalapi

import (
	"net/http"

	"github.com/kyma-project/kyma/components/connector-service/internal/httphelpers"

	"github.com/kyma-project/kyma/components/connector-service/internal/certificates"

	"github.com/kyma-project/kyma/components/connector-service/internal/clientcontext"
	"github.com/kyma-project/kyma/components/connector-service/internal/tokens"

	"github.com/gorilla/mux"
	"github.com/kyma-project/kyma/components/connector-service/internal/errorhandler"
)

type Config struct {
	Middlewares          []mux.MiddlewareFunc
	TokenService         tokens.Service
	ContextExtractor     clientcontext.ConnectorClientExtractor
	CertificateURL       string
	GetInfoURL           string
	ConnectorServiceHost string
	Subject              certificates.CSRSubject
	CertService          certificates.Service
}

type SignatureHandler interface {
	SignCSR(w http.ResponseWriter, r *http.Request)
}

type CSRGetInfoHandler interface {
	GetCSRInfo(w http.ResponseWriter, r *http.Request)
}

type MngmtGetInfoHandler interface {
	GetManagementInfo(w http.ResponseWriter, r *http.Request)
}

const apiSpecPath = "connectorapi.yaml"

func NewHandler(appHandlerCfg, runtimeHandlerCfg, appMngmtInfoHandlerCfg, runtimeMngmtInfoHandlerCfg Config, globalMiddlewares []mux.MiddlewareFunc) http.Handler {
	router := mux.NewRouter()

	httphelpers.WithMiddlewares(globalMiddlewares, router)

	router.Path("/v1").Handler(http.RedirectHandler("/v1/api.yaml", http.StatusMovedPermanently)).Methods(http.MethodGet)
	router.Path("/v1/api.yaml").Handler(NewStaticFileHandler(apiSpecPath)).Methods(http.MethodGet)

	applicationInfoHandler := NewCSRInfoHandler(appHandlerCfg.TokenService, appHandlerCfg.ContextExtractor, appHandlerCfg.CertificateURL, appHandlerCfg.GetInfoURL, appHandlerCfg.ConnectorServiceHost, appHandlerCfg.Subject, AppURLFormat)
	applicationManagementInfoHandler := NewManagementInfoHandler(appMngmtInfoHandlerCfg.ContextExtractor, appMngmtInfoHandlerCfg.ConnectorServiceHost)
	applicationSignatureHandler := NewSignatureHandler(appHandlerCfg.TokenService, appHandlerCfg.CertService, appHandlerCfg.ContextExtractor)

	csrApplicationRouter := router.PathPrefix("/v1/applications/signingRequests").Subrouter()
	csrApplicationRouter.HandleFunc("/info", applicationInfoHandler.GetCSRInfo).Methods(http.MethodGet)

	certApplicationRouter := router.PathPrefix("/v1/applications/certificates").Subrouter()
	certApplicationRouter.HandleFunc("/", applicationSignatureHandler.SignCSR).Methods(http.MethodPost)

	httphelpers.WithMiddlewares(appHandlerCfg.Middlewares, csrApplicationRouter, certApplicationRouter)

	mngmtApplicationRouter := router.PathPrefix("/v1/applications/management").Subrouter()
	mngmtApplicationRouter.HandleFunc("/info", applicationManagementInfoHandler.GetManagementInfo).Methods(http.MethodGet)

	httphelpers.WithMiddlewares(appMngmtInfoHandlerCfg.Middlewares, mngmtApplicationRouter)

	runtimeInfoHandler := NewCSRInfoHandler(runtimeHandlerCfg.TokenService, runtimeHandlerCfg.ContextExtractor, runtimeHandlerCfg.CertificateURL, runtimeHandlerCfg.GetInfoURL, runtimeHandlerCfg.ConnectorServiceHost, runtimeHandlerCfg.Subject, RuntimeURLFormat)
	runtimeManagementInfoHandler := NewManagementInfoHandler(runtimeMngmtInfoHandlerCfg.ContextExtractor, runtimeMngmtInfoHandlerCfg.ConnectorServiceHost)
	runtimeSignatureHandler := NewSignatureHandler(runtimeHandlerCfg.TokenService, runtimeHandlerCfg.CertService, runtimeHandlerCfg.ContextExtractor)

	csrRuntimesRouter := router.PathPrefix("/v1/runtimes/signingRequests").Subrouter()
	csrRuntimesRouter.HandleFunc("/info", runtimeInfoHandler.GetCSRInfo).Methods(http.MethodGet)

	certRuntimesRouter := router.PathPrefix("/v1/runtimes/certificates").Subrouter()
	certRuntimesRouter.HandleFunc("/", runtimeSignatureHandler.SignCSR).Methods(http.MethodPost)

	httphelpers.WithMiddlewares(runtimeHandlerCfg.Middlewares, csrRuntimesRouter, certApplicationRouter)

	mngmtRuntimeRouter := router.PathPrefix("/v1/runtimes/management").Subrouter()
	mngmtRuntimeRouter.HandleFunc("/info", runtimeManagementInfoHandler.GetManagementInfo).Methods(http.MethodGet)

	router.NotFoundHandler = errorhandler.NewErrorHandler(404, "Requested resource could not be found.")
	router.MethodNotAllowedHandler = errorhandler.NewErrorHandler(405, "Method not allowed.")

	return router
}
