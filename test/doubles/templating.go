package doubles

import "github.com/alcalbg/gotdd"

func NewStubTemplatingEngine(template string) gotdd.TemplatingEngine {
	app := gotdd.NewApp(gotdd.App{})
	return app.GetEngine().
		MountFS(
			NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(template),
				}))
}
