package authenticators

import (
	"fmt"
	"log"
	"strings"
)

// registry holds the registered authenticator plugins.
var registry = make(map[string]Authenticator)

// Register adds an authenticator plugin to the registry.
// It should be called from the init() function of each plugin implementation.
func Register(name string, plugin Authenticator) {
	normalizedName := strings.ToLower(name)
	if _, exists := registry[normalizedName]; exists {
		log.Printf("Warning: Authenticator plugin '%s' is already registered. Overwriting.", normalizedName)
	}
	//log.Printf("Registering authenticator plugin: %s", normalizedName)
	registry[normalizedName] = plugin
}

// Get retrieves an authenticator plugin from the registry by name.
// Returns the plugin and true if found, otherwise nil and false.
func Get(name string) (Authenticator, error) {
	normalizedName := strings.ToLower(name)
	plugin, exists := registry[normalizedName]
	if !exists {
		// Provide a list of registered authenticators in the error
		known := make([]string, 0, len(registry))
		for k := range registry {
			known = append(known, k)
		}
		return nil, fmt.Errorf("unknown authenticator '%s' requested (known: %v)", normalizedName, known)
	}
	return plugin, nil
}
