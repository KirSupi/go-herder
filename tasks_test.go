package herder

import "testing"

func TestPythonTask(t *testing.T) {
	tc := DefaultPythonScriptTask("./script.py")
	t.Log(tc)
}
