package server

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"

	"github.com/spaceuptech/space-cloud/gateway/server/handlers"
)

func (s *Server) routes(profiler bool, staticPath string, restrictedHosts []string) http.Handler {
	// Add a '*' to the restricted hosts if length is zero
	if len(restrictedHosts) == 0 {
		restrictedHosts = append(restrictedHosts, "*")
	}

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/v1/config/credentials").HandlerFunc(handlers.HandleGetCredentials(s.managers.Admin()))
	router.Methods(http.MethodGet).Path("/v1/config/permissions").HandlerFunc(handlers.HandleGetPermissions(s.managers.Admin()))

	router.Methods(http.MethodPost).Path("/v1/config/generate-token").HandlerFunc(handlers.HandleGenerateAdminToken(s.managers.Admin()))

	router.Methods(http.MethodPost).Path("/v1/config/integrations").HandlerFunc(handlers.HandlePostIntegration(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/integrations").HandlerFunc(handlers.HandleGetIntegrations(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/integrations/{name}").HandlerFunc(handlers.HandleDeleteIntegration(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/integrations/tokens").HandlerFunc(handlers.HandleGetIntegrationTokens(s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/integrations/{name}/hooks").HandlerFunc(handlers.HandleAddIntegrationHook(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/integrations/{name}/hooks").HandlerFunc(handlers.HandleGetIntegrationHooks(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/integrations/{name}/hooks/{id}").HandlerFunc(handlers.HandleDeleteIntegrationHook(s.managers.Admin(), s.managers.Sync()))

	// Initialize the routes for config management
	router.Methods(http.MethodGet).Path("/v1/config/env").HandlerFunc(handlers.HandleLoadEnv(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/login").HandlerFunc(handlers.HandleAdminLogin(s.managers.Admin()))
	router.Methods(http.MethodGet).Path("/v1/config/refresh-token").HandlerFunc(handlers.HandleRefreshToken(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleGetProjectConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleApplyProject(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}").HandlerFunc(handlers.HandleDeleteProjectConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/generate-internal-token").HandlerFunc(handlers.HandleGenerateTokenForMissionControl(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/cluster").HandlerFunc(handlers.HandleGetClusterConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/cluster").HandlerFunc(handlers.HandleSetClusterConfig(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/caching/config").HandlerFunc(handlers.HandleGetCacheConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/caching/config/{id}").HandlerFunc(handlers.HandleSetCacheConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/external/caching/connection-state").HandlerFunc(handlers.HandleGetCacheConnectionState(s.managers.Admin(), s.modules.Caching()))
	router.Methods(http.MethodDelete).Path("/v1/external/projects/{project}/caching/purge-cache").HandlerFunc(handlers.HandlePurgeCache(s.managers.Admin(), s.modules.Caching()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/letsencrypt/config").HandlerFunc(handlers.HandleGetEncryptWhitelistedDomain(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/letsencrypt/config/{id}").HandlerFunc(handlers.HandleLetsEncryptWhitelistedDomain(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing/ingress").HandlerFunc(handlers.HandleGetProjectRoute(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/ingress/global").HandlerFunc(handlers.HandleSetGlobalRouteConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodGet).Path("/v1/config/projects/{project}/routing/ingress/global").HandlerFunc(handlers.HandleGetGlobalRouteConfig(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodPost).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleSetProjectRoute(s.managers.Admin(), s.managers.Sync()))
	router.Methods(http.MethodDelete).Path("/v1/config/projects/{project}/routing/ingress/{id}").HandlerFunc(handlers.HandleDeleteProjectRoute(s.managers.Admin(), s.managers.Sync()))

	router.Methods(http.MethodPost).Path("/v1/config/batch-apply").HandlerFunc(handlers.HandleBatchApplyConfig(s.managers.Admin()))

	// Health check
	router.Methods(http.MethodGet).Path("/v1/api/health-check").HandlerFunc(handlers.HandleHealthCheck(s.managers.Sync()))

	// Register pprof handlers if profiler set to true
	if profiler {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
	}

	// forward request for project mutation, websocket, getting cluster type
	runnerRouter := router.PathPrefix("/v1/runner").HandlerFunc(s.managers.Sync().HandleRunnerRequests(s.managers.Admin())).Subrouter()
	// secret routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}").HandlerFunc(s.managers.Sync().HandleRunnerApplySecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/secrets").HandlerFunc(s.managers.Sync().HandleRunnerListSecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/root-path").HandlerFunc(s.managers.Sync().HandleRunnerSetFileSecretRootPath(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteSecret(s.managers.Admin()))
	runnerRouter.Methods(http.MethodPost).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.managers.Sync().HandleRunnerSetSecretKey(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/secrets/{id}/{key}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteSecretKey(s.managers.Admin()))
	// service
	runnerRouter.Methods(http.MethodPost).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.managers.Sync().HandleRunnerApplyService(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/services").HandlerFunc(s.managers.Sync().HandleRunnerGetServices(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/services/{serviceId}/{version}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteService(s.managers.Admin()))
	// service routes
	runnerRouter.Methods(http.MethodPost).Path("/{project}/service-routes/{serviceId}").HandlerFunc(s.managers.Sync().HandleRunnerServiceRoutingRequest(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/service-routes").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceRoutingRequest(s.managers.Admin()))

	// service role
	runnerRouter.Methods(http.MethodPost).Path("/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.managers.Sync().HandleRunnerSetServiceRole(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/service-roles").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceRoleRequest(s.managers.Admin()))
	runnerRouter.Methods(http.MethodDelete).Path("/{project}/service-roles/{serviceId}/{roleId}").HandlerFunc(s.managers.Sync().HandleRunnerDeleteServiceRole(s.managers.Admin()))

	runnerRouter.Methods(http.MethodGet).Path("/{project}/services/logs").HandlerFunc(s.managers.Sync().HandleRunnerGetServiceLogs(s.managers.Admin()))
	runnerRouter.Methods(http.MethodGet).Path("/{project}/services/status").HandlerFunc(s.managers.Sync().HandleRunnerGetDeploymentStatus(s.managers.Admin()))

	if staticPath != "" {
		// Add handler for mission control
		router.PathPrefix("/mission-control").HandlerFunc(handlers.HandleMissionControl(staticPath))
	}

	// Add handler for routing module
	router.PathPrefix("/").HandlerFunc(s.modules.Routing().HandleRoutes(s.modules))
	return s.restrictDomainMiddleware(restrictedHosts, router)
}
