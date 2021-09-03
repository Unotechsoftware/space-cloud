package syncman

import (
	"context"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

// AdminSyncmanInterface is an interface consisting of functions of admin module used by eventing module
type AdminSyncmanInterface interface {
	GetInternalAccessToken() (string, error)
	IsTokenValid(ctx context.Context, token, resource, op string, attr map[string]string) (model.RequestParams, error)
	SetServices(eventType string, services model.ScServices)
	ValidateProjectSyncOperation(c *config.Config, project *config.ProjectConfig) bool
	SetIntegrationConfig(integrations config.Integrations)

	// For integrations
	GetIntegrationToken(id string) (string, error)
}

type integrationInterface interface {
	SetConfig(integrations config.Integrations, integrationHooks config.IntegrationHooks) error
	SetIntegrations(integrations config.Integrations) error
	SetIntegrationHooks(integrationHooks config.IntegrationHooks)
	InvokeHook(context.Context, model.RequestParams) config.IntegrationAuthResponse
}

// ModulesInterface is an interface consisting of functions of the modules module used by syncman
type ModulesInterface interface {
	// SetInitialProjectConfig sets the config all modules
	SetInitialProjectConfig(ctx context.Context, config config.Projects) error

	// SetProjectConfig sets specific project config
	SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error

	SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error

	SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error
	SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error

	// Getters
	GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error)
	LetsEncrypt() *letsencrypt.LetsEncrypt
	Routing() *routing.Routing
	Caching() *caching.Cache

	// Delete
	Delete(projectID string)
}

// GlobalModulesInterface is an interface consisting of functions of the global modules
type GlobalModulesInterface interface {
	// SetMetricsConfig set the config of the metrics module
	SetMetricsConfig(isMetricsEnabled bool)
}
