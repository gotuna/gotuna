<p align="center">
<img src="https://avatars.githubusercontent.com/u/82163094?s=200&v=4">
</p>


<p align="center">
<a href="https://pkg.go.dev/github.com/gotuna/gotuna"><img src="https://pkg.go.dev/badge/github.com/gotuna/gotuna" alt="PkgGoDev"></a>
<a href="https://github.com/gotuna/gotuna/actions"><img src="https://github.com/gotuna/gotuna/workflows/Tests/badge.svg" alt="Tests Status" /></a>
<a href="https://goreportcard.com/report/github.com/gotuna/gotuna"><img src="https://goreportcard.com/badge/github.com/gotuna/gotuna" alt="Go Report Card" /></a>
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

# Testing
```
go test -race -v ./...
```

# Running examples
```
go run examples/fullapp/cmd/main.go
```

# External dependencies
External modules are mostly used when the feature is too complex to build or maintain - Router, Secure cookies

# TODO
- Validation (input/forms)
- Cache

# Licence
This project is licensed under the MIT License.
