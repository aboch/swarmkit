package networkallocator

import (
	"github.com/docker/libnetwork/drivers/overlay/ovmanager"
	"github.com/docker/libnetwork/drivers/remote"
)

func getInitializers() []initializer {
	return []initializer{
		{remote.Init, "remote"},
		{ovmanager.Init, "overlay"},
	}
}

// GetPredefinedNetworks returns the list of predefined network structures
func GetPredefinedNetworks() []PredefinedNetwork {
	return nil
}
