package processor

import (
	"log"
)

// Processor manages and runs registered services.
type Processor struct {
	services []Service
}

// New creates a new Processor.
func New() *Processor {
	return &Processor{services: make([]Service, 0)}
}

// Add installs and registers a Service. Terminates on install failure.
func (p *Processor) Add(s Service) {
	if err := s.Install(); err != nil {
		log.Fatalf("installation failed: %s", err)
	}

	p.services = append(p.services, s)
}

// Run executes all registered services in order. Stops on first error.
func (p *Processor) Run() {
	for _, svc := range p.services {
		if err := svc.Process(); err != nil {
			log.Fatalf("service failed: %s", err)
		}
	}
}
