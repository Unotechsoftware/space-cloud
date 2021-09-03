package modules

import (
	"github.com/spaceuptech/space-cloud/gateway/managers"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/global"
)

// Module is an object that sets up the modules
type Module struct {
	auth      *auth.Module

	// Global Modules
	GlobalMods *global.Global

	// Managers
	Managers *managers.Managers
}

func newModule(projectID, clusterID, nodeID string, managers *managers.Managers, globalMods *global.Global) (*Module, error) {
	// Get managers
	adminMan := managers.Admin()
	syncMan := managers.Sync()
	integrationMan := managers.Integration()

	a := auth.Init(clusterID, nodeID, adminMan, integrationMan)
	a.SetMakeHTTPRequest(syncMan.MakeHTTPRequest)

	return &Module{auth: a, Managers: managers, GlobalMods: globalMods}, nil
}
