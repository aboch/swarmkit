package networkallocator

import (
	"github.com/docker/libnetwork/drivers/bridge/brmanager"
	"github.com/docker/libnetwork/drivers/host"
	"github.com/docker/libnetwork/drivers/ipvlan/ivmanager"
	"github.com/docker/libnetwork/drivers/macvlan/mvmanager"
	"github.com/docker/libnetwork/drivers/overlay/ovmanager"
	"github.com/docker/libnetwork/drivers/remote"
)

func getInitializers() []initializer {
	return []initializer{
		{remote.Init, "remote"},
		{ovmanager.Init, "overlay"},
		{mvmanager.Init, "macvlan"},
		{brmanager.Init, "bridge"},
		{ivmanager.Init, "ipvlan"},
		{host.Init, "host"},
	}
}

// GetPredefinedNetworks returns the list of predefined network structures
func GetPredefinedNetworks() []PredefinedNetwork {
	return []PredefinedNetwork{
		{"bridge", "bridge"},
		{"host", "host"},
	}
}
