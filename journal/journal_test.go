package journal

import (
	"bytes"
	"testing"

	"github.com/nesv/mct/instruction"
)

func TestCommandUnmarshalText(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		instr instruction.Instruction
		args  [][]byte
		fail  bool
	}{
		{
			name:  "rem",
			input: []byte(`REM hello, there!`),
			instr: instruction.Rem,
			args: [][]byte{
				[]byte(`hello,`),
				[]byte(`there!`),
			},
		},
		{
			name:  "exec",
			input: []byte(`EXEC apt install -y coredns`),
			instr: instruction.Exec,
			args: [][]byte{
				[]byte(`apt`),
				[]byte(`install`),
				[]byte(`-y`),
				[]byte(`coredns`),
			},
		},
	}
	for _, tt := range tests {
		var cmd Command
		err := cmd.UnmarshalText(tt.input)
		if err == nil && tt.fail {
			t.Fatal("should have failed")
		} else if err != nil && !tt.fail {
			t.Fatal(err)
		} else if err != nil {
			return
		}
		if cmd.Instruction != tt.instr {
			t.Errorf("wrong instruction: want=%q, got=%q", tt.instr, cmd.Instruction)
		}
		if want, got := len(tt.args), len(cmd.Args); want != got {
			t.Errorf("wrong number of arguments: want=%d, got=%d", want, got)
		}
		for i, got := range cmd.Args {
			want := tt.args[i]
			if !bytes.Equal(want, got) {
				t.Errorf("arg[%d] mismatch: want=%q, got=%q", i, want, got)
			}
		}
	}
}
