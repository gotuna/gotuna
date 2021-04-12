package doubles

import "github.com/gotuna/gotuna"

// NewStubTemplatingEngine returns new native HTML engine with a stub template
func NewStubTemplatingEngine(template string) gotuna.TemplatingEngine {
	return gotuna.App{
		ViewFiles: NewFileSystemStub(
			map[string][]byte{
				"view.html": []byte(template),
			}),
	}.NewTemplatingEngine()
}
