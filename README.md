# GoTuna - Web framework written in Go
Experimental Go web framework. Used to experiment with Go, try the latest language features, explore, and learn more.

NOTE: This project is under heavy development and is not ready for use.

# Features
- MVC pattern
- Router (Gorilla)
- Middleware support
- Session managment (Gorilla)
- User authentication (via user provider interface)
- Native view rendering (html/template)
- Multi-language support

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
- Cache
- Validation

# Licence
This project is licensed under the MIT License.
