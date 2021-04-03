package doubles

import "github.com/alcalbg/gotdd"

func NewStubTemplatingEngine(template string) gotdd.TemplatingEngine {
	return gotdd.App{}.GetEngine().
		MountViews(
			NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(template),
				}))
}
