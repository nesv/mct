package instruction

import "fmt"

// Instruction is a command that the system can perform.
type Instruction uint32

const (
	Rem Instruction = 1 << iota
	Mkdir
	Copy
	Chmod
	Chown
	Chgrp
	Rm
	Exec
)

// String returns the textual representation of an instruction.
func (i Instruction) String() string {
	switch i {
	case Rem:
		return "REM"
	case Mkdir:
		return "MKDIR"
	case Copy:
		return "COPY"
	case Chmod:
		return "CHMOD"
	case Chown:
		return "CHOWN"
	case Chgrp:
		return "CHGRP"
	case Rm:
		return "RM"
	case Exec:
		return "Exec"
	default:
		return "???"
	}
}

func (i Instruction) MarshalText() ([]byte, error) {
	s := i.String()
	if s == "???" {
		return nil, fmt.Errorf("unknown instruction: %d", int32(i))
	}
	return []byte(s), nil
}

func (i *Instruction) UnmarshalText(p []byte) error {
	m := map[string]Instruction{
		"REM":   Rem,
		"MKDIR": Mkdir,
		"COPY":  Copy,
		"CHMOD": Chmod,
		"CHOWN": Chown,
		"CHGRP": Chgrp,
		"RM":    Rm,
		"EXEC":  Exec,
	}
	v, ok := m[string(p)]
	if !ok {
		return fmt.Errorf("unknown command: %q", p)
	}
	*i = v
	return nil
}
