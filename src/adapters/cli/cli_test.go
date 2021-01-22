package cli

import (
	"github.com/gitaction/docopt.go"
	"testing"
)

var input = "input"
var usage = "Usage: echo <input>"

func testFunc(opts docopt.Opts) {
	userInput, _ := opts.String("<input>")
	input = userInput
}

func TestRegister(t *testing.T) {
	_ = Register("echo", testFunc, usage)
	
	if found :=findCmd("echo"); found == false {
		t.Errorf("findCmd('echo') = false, want true")
	}

	if found :=findCmd("echo2"); found == true {
		t.Errorf("findCmd('echo2') = false, want true")
	}
}

func TestRun(t *testing.T) {
	if _ = Run("echo", []string{"hello"}); input != "hello" {
		t.Errorf("Run('echo', []string{'hello'}) wrong, want echo hello")
	}
}
