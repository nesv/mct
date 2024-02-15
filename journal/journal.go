// Package journal provides encoding and decoding microcode logs (a.k.a. "journals").
package journal

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/nesv/mct/instruction"
)

// ReadFrom reads a journal from r, sending each parsed entry to the provided channel.
// A non-nil error will be returned if an entry cannot be parsed.
// The entries channel will be closed when ReadFrom returns.
// Callers can force an early return by cancelling ctx.
func ReadFrom(ctx context.Context, r io.Reader, entries chan<- Entry) error {
	defer close(entries)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var entry Entry
		if err := entry.UnmarshalText(scanner.Bytes()); err != nil {
			return fmt.Errorf("unmarshal text: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case entries <- entry:
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan: %w", err)
	}
	return nil
}

type Entry struct {
	Command Command
	Action  *Action
	Revert  *Command
}

func (e *Entry) UnmarshalText(text []byte) error {
	input := bytes.TrimSpace(text)
	if len(input) == 0 {
		return ErrEmptyLine
	}
	// "REM" (remark) commands do not have an action
	if bytes.HasPrefix(input, []byte(instruction.Rem.String())) {
		return e.Command.UnmarshalText(input)
	}
	actionSep := bytes.Index(input, []byte("&&"))
	if actionSep == -1 {
		return ErrMissingAction
	}
	cmdText := bytes.TrimSpace(input[:actionSep])
	var cmd Command
	if err := cmd.UnmarshalText(cmdText); err != nil {
		return fmt.Errorf("parse command: %w", err)
	}
	undoSep := bytes.LastIndex(input, []byte("&&"))
	if undoSep == actionSep {
		// There is no undo part for the entry.
		undoSep = len(input)
	}
	actionStart := actionSep + len("&& ")
	actionText := bytes.TrimSpace(input[actionStart:undoSep])
	var action Action
	if err := action.UnmarshalText(actionText); err != nil {
		return fmt.Errorf("parse action: %w", err)
	}
	if e == nil {
		e = &Entry{
			Command: cmd,
			Action:  &action,
		}
	} else {
		e.Command, e.Action = cmd, &action
	}
	if undoSep == len(input) {
		// No undo part to parse.
		return nil
	}
	undoText := bytes.TrimSpace(input[undoSep+len("&&"):])
	var undo Command
	if err := undo.UnmarshalText(undoText); err != nil {
		return fmt.Errorf("parse undo/revert command: %w", err)
	}
	e.Revert = &undo
	return nil
}

var (
	ErrEmptyLine     = errors.New("empty line")
	ErrMissingAction = errors.New("missing action")
)

type Action struct {
	Instruction ActionInstruction
	Args        [][]byte
}

func (a *Action) UnmarshalText(input []byte) error {
	if len(bytes.TrimSpace(input)) == 0 {
		return errors.New("empty action")
	}
	fields := bytes.Fields(input)
	var in ActionInstruction
	if err := in.UnmarshalText(fields[0]); err != nil {
		return fmt.Errorf("unmarshal instruction: %w", err)
	}
	if len(fields) < 2 {
		// No arguments.
		return nil
	}
	args := make([][]byte, len(fields[1:]))
	for i, f := range fields {
		arg := make([]byte, len(f))
		copy(arg, f)
		args[i] = arg
	}
	if a == nil {
		a = &Action{
			Instruction: in,
			Args:        args,
		}
		return nil
	}
	a.Instruction, a.Args = in, args
	return nil
}

type ActionInstruction uint32

const (
	Nop ActionInstruction = 1 << iota
	Sysctl
)

func (a *ActionInstruction) UnmarshalText(text []byte) error {
	m := map[string]ActionInstruction{
		"NOP":    Nop,
		"SYSCTL": Sysctl,
	}
	v, ok := m[string(text)]
	if !ok {
		return fmt.Errorf("unknown instruction: %q", text)
	}
	if a == nil {
		a = &v
		return nil
	}
	*a = v
	return nil
}

type Command struct {
	Instruction instruction.Instruction
	Args        [][]byte
}

func (c *Command) UnmarshalText(input []byte) error {
	if len(bytes.TrimSpace(input)) == 0 {
		return errors.New("empty command")
	}
	fields := bytes.Fields(input)
	var in instruction.Instruction
	if err := in.UnmarshalText(fields[0]); err != nil {
		return fmt.Errorf("unmarshal instruction: %w", err)
	}
	if len(fields) < 2 {
		// Nothing else to parse.
		return nil
	}
	args := make([][]byte, len(fields[1:]))
	for i, f := range fields[1:] {
		// Copy the byte slices in case they are reused elsewhere.
		arg := make([]byte, len(f))
		copy(arg, f)
		args[i] = arg
	}
	if c == nil {
		c = &Command{
			Instruction: in,
			Args:        args,
		}
		return nil
	}
	c.Instruction, c.Args = in, args
	return nil
}
