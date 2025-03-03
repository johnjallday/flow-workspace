package commands

// Command defines the interface that all commands must implement.
type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

// BaseCommand provides a simple implementation of Command using a handler function.
type BaseCommand struct {
	name        string
	description string
	handler     func(args []string) error
}

func (c BaseCommand) Name() string {
	return c.name
}

func (c BaseCommand) Description() string {
	return c.description
}

func (c BaseCommand) Execute(args []string) error {
	return c.handler(args)
}

// NewCommand is a helper to create a new Command.
func NewCommand(name, description string, handler func(args []string) error) Command {
	return BaseCommand{name: name, description: description, handler: handler}
}
