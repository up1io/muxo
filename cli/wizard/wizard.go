package wizard

import "context"

// Note(john):
// We will later have two types of wizards.
// The first is the project wizard which set up the foundation project.
// Then there is also the module wizard which encapsulated a specified logic in a module which can be added to project.

// Wizard define process to set up or load functionalities to a muxo project.
type Wizard interface {
	Execute(ctx context.Context) error
}
