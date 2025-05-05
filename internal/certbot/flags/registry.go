package flags

// registry holds the registered flag generator instances. Order might matter for CLI readability.
var registry = []FlagGenerator{}

// Register adds a flag generator instance to the registry.
// Called from init() functions in implementation files.
func Register(generator FlagGenerator) {
	//typeName := reflect.TypeOf(generator).Elem().Name() // Get struct name
	//log.Printf("Registering flag generator: %s", typeName)
	registry = append(registry, generator)
}

// GetAll retrieves all registered flag generators.
func GetAll() []FlagGenerator {
	// Return a copy to prevent external modification? For now, return direct slice.
	return registry
}
