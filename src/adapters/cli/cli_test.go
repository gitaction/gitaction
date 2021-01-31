package cli

import (
	"github.com/gitaction/docopt.go"
	"testing"
)

var input = "input"
func testFunc(opts docopt.Opts) {
	userInput, _ := opts.String("<input>")
	input = userInput
}

func TestRegister(t *testing.T) {
	_ = Register("testR", testFunc)
	
	if found :=findCmd("testR"); found == false {
		t.Errorf("findCmd('testR') = false, want true")
	}

	if found :=findCmd("testR2"); found == true {
		t.Errorf("findCmd('testR2') = false, want true")
	}
}

func TestRun(t *testing.T) {
	_ = Register("testRun", testFunc)
	opts, _ := docopt.ParseArgs("Usage: example testRun <input>", []string{"testRun", "hello"}, "0.0.1")
	
	if _ = Run(opts); input != "hello" {
		t.Errorf("Run('testRun', []string{'hello'}) wrong, want echo hello")
	}
}
