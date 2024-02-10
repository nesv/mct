package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

func main() {
	fmt.Println("Miek's Configuration Tool")

	cmds, err := ParseCommands(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', 0)
	defer tw.Flush()
	for _, cmd := range cmds {
		fmt.Fprintln(tw, strings.Join([]string{
			fmt.Sprintf("%8s", cmd.Name),
			fmt.Sprintf("%v", cmd.Args),
			"&&",
			cmd.Action.String(),
		}, "\t"))
	}
}

func ParseCommands(r io.Reader) ([]Command, error) {
	var (
		cmds    []Command
		scanner = bufio.NewScanner(r)
		line    int
	)
	for scanner.Scan() {
		line++
		c, err := ParseCommandString(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("parse command (line %d): %w", line, err)
		}
		cmds = append(cmds, c)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return cmds, nil
}

func ParseCommandString(input string) (Command, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return Command{}, ErrEmptyLine
	}

	var c Command
	if err := c.UnmarshalText([]byte(input)); err != nil {
		return Command{}, fmt.Errorf("unmarshal text: %w", err)
	}
	return c, nil
}

// ErrEmptyLine is used to indicate an empty string was passed to
// [ParseCommandString] or [Command.UnmarshalText], and does not necessarily
// indicate an error.
var ErrEmptyLine = errors.New("empty line")

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
// c.UnmarshalText.
func (c Command) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

// String returns a string representation of c.
func (c Command) String() string {
	return fmt.Sprintf("%8s\t+%v\t&& %s", c.Name, c.Args, c.Action)
}

type Action struct {
	Name string
	Args []string
}

func (a Action) String() string {
	return fmt.Sprintf("%8s\t%+v", a.Name, a.Args)
}

func (a Action) Equals(other Action) bool {
	return a.Name == other.Name && slices.Compare(a.Args, other.Args) == 0
}
