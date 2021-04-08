package doubles

import "github.com/gotuna/gotuna"

func NewStubTemplatingEngine(template string) gotuna.TemplatingEngine {
	return gotuna.App{
		ViewFiles: NewFileSystemStub(
			map[string][]byte{
				"view.html": []byte(template),
			}),
	}.NewTemplatingEngine()
}
