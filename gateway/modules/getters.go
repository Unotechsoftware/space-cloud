package modules

import (
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/caching"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

// Auth returns the auth module
func (m *Modules) Auth(projectID string) (*auth.Module, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.auth, nil
}

// GetAuthModuleForSyncMan returns auth module for sync manager
func (m *Modules) GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error) {
	module, err := m.loadModule(projectID)
	if err != nil {
		return nil, err
	}
	return module.auth, nil
}

// LetsEncrypt returns the letsencrypt module
func (m *Modules) LetsEncrypt() *letsencrypt.LetsEncrypt {
	return m.GlobalMods.LetsEncrypt()
}

// Routing returns the routing module
func (m *Modules) Routing() *routing.Routing {
	return m.GlobalMods.Routing()
}

// Caching returns the caching module
func (m *Modules) Caching() *caching.Cache {
	return m.GlobalMods.Caching()
}