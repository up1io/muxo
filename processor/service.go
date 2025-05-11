package processor

// Service represents a unit with dependency installation and processing steps.
type Service interface {
	// Install ensures dependencies are met before execution.
	Install() error
	// Process runs the service's main logic.
	Process() error
}
