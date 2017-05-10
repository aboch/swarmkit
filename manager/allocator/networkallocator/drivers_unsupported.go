// +build !linux,!darwin,!windows

package networkallocator

func getInitializers() []initializer {
	return nil
}

// Returns the list of predefined networks as
// a list of {network name, driver name} tuples
func GetPredefinedNetworks() []PredefinedNetwork {
	return nil
}
