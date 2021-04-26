<p align="center">
<img src="https://avatars.githubusercontent.com/u/82163094?s=200&v=4">
</p>


<p align="center">
<a href="https://pkg.go.dev/github.com/gotuna/gotuna"><img src="https://pkg.go.dev/badge/github.com/gotuna/gotuna" alt="PkgGoDev"></a>
<a href="https://github.com/gotuna/gotuna/actions"><img src="https://github.com/gotuna/gotuna/workflows/Tests/badge.svg" alt="Tests Status" /></a>
<a href="https://goreportcard.com/report/github.com/gotuna/gotuna"><img src="https://goreportcard.com/badge/github.com/gotuna/gotuna" alt="Go Report Card" /></a>
<a href="https://codecov.io/gh/gotuna/gotuna"><img src="https://codecov.io/gh/gotuna/gotuna/branch/main/graph/badge.svg?token=QG7CG4MSPC" alt="Go Report Card" /></a>
</p>

# GoTuna - Web framework for Go
Please visit [https://gotuna.org](https://gotuna.org)  for the latest documentation, examples, and more.


# Features
- Router (gorilla)
- Standard `http.Handler` interface
- Middleware support
- User session management (gorilla)
- Session flash messages
- Native view rendering (html/template) with helpers
- Static file server included with the configurable prefix
- Standard logger interface
- Request logging and panic recovery
- Full support for embedded templates and static files
- User authentication (via user provider interface)
- Sample InMemory user provider included
- Multi-language support
- Database agnostic

# Requirements
- Make sure you have Go >= 1.16 installed

# Quick Start
Initialize new app and install GoTuna:

```shell
mkdir testapp
cd testapp
go get -u github.com/gotuna/gotuna
```

Now create two files `main.go` and `app.html` as an example:

```go
// main.go

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gotuna/gotuna"
)

func main() {
	app := gotuna.App{
		ViewFiles: os.DirFS("."),
	}
	app.Router = gotuna.NewMuxRouter()
	app.Router.Handle("/", handlerHome(app))
	app.Router.Handle("/login", handlerLogin(app)).Methods(http.MethodGet, http.MethodPost)

	fmt.Println("Running on http://localhost:8888")
	http.ListenAndServe(":8888", app.Router)
}

func handlerHome(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.NewTemplatingEngine().
			Render(w, r, "app.html")
	})
}

func handlerLogin(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Login form...")
	})
}
```

This will be your app's html layout:

```html
// app.html

{{- define "app" -}}
<!DOCTYPE html>
<html>
  <head></head>
  <body>
    <a href="/login">Please login!</a>
  </body>
</html>
{{- end -}}
```

Run this simple app and visit http://localhost:8888 in your browser:
```shell
go run main.go
```


# Running example apps
GoTuna comes with few working examples. Make sure you have git and Go >= 1.16 installed.
```shell
git clone https://github.com/gotuna/gotuna.git
cd gotuna
go run examples/fullapp/cmd/main.go
```

# Testing

```shell
go test -race -v ./...
```

# Licence
This project is licensed under the MIT License.
