package microcode

import (
	"bytes"
	"strings"
)

// Command defines the structure of a configuration command.
type Command struct {
	Name   string
	Args   []string
	Action Action
}

// UnmarshalText unmarshals a textual representation of a [Command].
func (c *Command) UnmarshalText(p []byte) error {
	input := bytes.TrimSpace(p)
	if len(input) == 0 {
		return ErrEmptyLine
	}

	// See if there was an action specified.
	var (
		action    Action
		actionSep = bytes.Index(input, []byte("&&"))
	)
	if actionSep > 0 {
		start := actionSep + len("&&")
		fields := bytes.Fields(input[start:])

		var (
			name = string(fields[0])
			args []string
		)
		if len(fields) > 1 {
			args = make([]string, len(fields[1:]))
			for i, field := range fields[1:] {
				args[i] = string(field)
			}
		}
		action = Action{
			Name: name,
			Args: args,
		}
	}
	if actionSep == -1 {
		actionSep = len(input)
	}
	c.Action = action

	fields := bytes.Fields(input[:actionSep])

	c.Name = string(fields[0])

	var args []string
	if len(fields) > 1 {
		args = make([]string, len(fields[1:]))
		for i, field := range fields[1:] {
			args[i] = string(field)
		}
	}
	c.Args = args

	return nil
}

// MarshalText returns a textual representation of c that is parseable with
// [Command.UnmarshalText].
func (c Command) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// String returns a string representation of c.
func (c Command) String() string {
	parts := []string{
		c.Name,
		strings.Join(c.Args, " "),
	}
	if !c.Action.IsZero() {
		parts = append(parts, "&&", c.Action.String())
	}
	return strings.Join(parts, " ")
}
