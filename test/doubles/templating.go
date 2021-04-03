package doubles

import "github.com/alcalbg/gotdd"

func NewStubTemplatingEngine(template string, options gotdd.Options) gotdd.TemplatingEngine {
	return gotdd.GetEngine(options).
		MountFS(
			NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(template),
				}))
}
