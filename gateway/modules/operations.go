package modules

import (
	"context"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// SetInitialProjectConfig sets the config all modules
func (m *Module) SetInitialProjectConfig(ctx context.Context, projects config.Projects) error {
	for projectID, project := range projects {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting config of auth module", nil)
		if err := m.auth.SetConfig(ctx, project.FileStoreConfig.StoreType, project.ProjectConfig, project.DatabaseRules, project.DatabasePreparedQueries, project.FileStoreRules, project.RemoteService, project.EventingRules); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set auth module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting config of lets encrypt module", nil)
		if err := m.GlobalMods.LetsEncrypt().SetProjectDomains(projectID, project.LetsEncrypt); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set letsencypt module config", err, nil)
		}

		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting config of ingress routing module", nil)
		if err := m.GlobalMods.Routing().SetProjectRoutes(projectID, project.IngressRoutes); err != nil {
			_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to set routing module config", err, nil)
		}
		m.GlobalMods.Routing().SetGlobalConfig(project.IngressGlobal)
		m.GlobalMods.Caching().AddDBRules(projectID, project.DatabaseRules)
	}
	return nil
}

// SetProjectConfig set project config
func (m *Module) SetProjectConfig(ctx context.Context, p *config.ProjectConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting project config", nil)
	if err := m.auth.SetProjectConfig(p); err != nil {
		return err
	}
	return nil
}

// SetLetsencryptConfig set the config of letsencrypt module
func (m *Module) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting letsencrypt config of project", nil)
	return m.GlobalMods.LetsEncrypt().SetProjectDomains(projectID, c)
}

// SetIngressRouteConfig set the config of routing module
func (m *Module) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting config of routing module", nil)
	return m.GlobalMods.Routing().SetProjectRoutes(projectID, routes)
}

// SetIngressGlobalRouteConfig set config of routing module
func (m *Module) SetIngressGlobalRouteConfig(ctx context.Context, _ string, c *config.GlobalRoutesConfig) error {
	helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "Setting config of global routing", nil)
	m.GlobalMods.Routing().SetGlobalConfig(c)
	return nil
}
