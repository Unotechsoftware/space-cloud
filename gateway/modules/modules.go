package modules

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
)

// Modules is an object that sets up the modules
type Modules struct {
	lock   sync.RWMutex
	blocks map[string]*Module

	clusterID string
	nodeID    string

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

// New creates a new modules instance
func New(_, clusterID, nodeID string, managers *managers.Managers, globalMods *global.Global) (*Modules, error) {
	return &Modules{
		blocks:     map[string]*Module{},
		clusterID:  clusterID,
		nodeID:     nodeID,
		GlobalMods: globalMods,
		Managers:   managers,
	}, nil
}

// SetInitialProjectConfig sets the config all modules
func (m *Modules) SetInitialProjectConfig(ctx context.Context, projects config.Projects) error {
	for projectID, project := range projects {
		module, err := m.loadModule(projectID)
		if err != nil {
			module, err = m.newModule(project.ProjectConfig)
			if err != nil {
				return err
			}
		}

		if err := module.SetInitialProjectConfig(ctx, config.Projects{projectID: project}); err != nil {
			return err
		}
	}
	return nil
}

// SetProjectConfig sets the config all modules
func (m *Modules) SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error {
	module, err := m.loadModule(config.ID)
	if err != nil {
		module, err = m.newModule(config)
		if err != nil {
			return err
		}
	}
	return module.SetProjectConfig(ctx, config)
}

// SetLetsencryptConfig set the config of letsencrypt module
func (m *Modules) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetLetsencryptConfig(ctx, projectID, c)
}

// SetIngressRouteConfig set the config of routing module
func (m *Modules) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetIngressRouteConfig(ctx, projectID, routes)
}

// SetIngressGlobalRouteConfig set config of routing module
func (m *Modules) SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error {
	module, err := m.loadModule(projectID)
	if err != nil {
		return err
	}
	return module.SetIngressGlobalRouteConfig(ctx, projectID, c)
}

func (m *Modules) projects() *config.Config {
	m.lock.RLock()
	defer m.lock.RUnlock()

	projects := make(config.Projects)
	for id := range m.blocks {
		projects[id] = &config.Project{ProjectConfig: &config.ProjectConfig{ID: id}}
	}
	return &config.Config{Projects: projects}
}

// Delete the project
func (m *Modules) Delete(projectID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if block, p := m.blocks[projectID]; p {
		// Close all the modules here
		helpers.Logger.LogDebug(helpers.GetRequestID(context.TODO()), "Closing config of auth module", nil)
		block.auth.CloseConfig()
	}

	delete(m.blocks, projectID)

	// Remove config from global modules
	_ = m.LetsEncrypt().DeleteProjectDomains(projectID)
	m.Routing().DeleteProjectRoutes(projectID)
}

func (m *Modules) loadModule(projectID string) (*Module, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if module, p := m.blocks[projectID]; p {
		return module, nil
	}

	return nil, fmt.Errorf("project (%s) not found in server state", projectID)
}

func (m *Modules) newModule(config *config.ProjectConfig) (*Module, error) {
	projects := m.projects()
	m.lock.Lock()
	defer m.lock.Unlock()

	if ok := m.Managers.Admin().ValidateProjectSyncOperation(projects, config); !ok {
		helpers.Logger.LogWarn("", "Cannot create new project. Upgrade your plan", nil)
		return nil, errors.New("upgrade your plan to create new project")
	}

	module, err := newModule(config.ID, m.clusterID, m.nodeID, m.Managers, m.GlobalMods)
	if err != nil {
		return nil, err
	}

	m.blocks[config.ID] = module
	return module, nil
}
