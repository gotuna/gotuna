package doubles

import (
	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
)

func NewStubTemplatingEngine(template string, options util.Options) templating.TemplatingEngine {
	return templating.GetEngine(options).
		MountFS(
			NewFileSystemStub(
				map[string][]byte{
					"view.html": []byte(template),
				}))
}
