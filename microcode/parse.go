package microcode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

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
