# Web Framework for Go
Experimental Go web framework. Used to experiment with Go, try the latest language features, explore, and learn more.

NOTE: This project is under heavy development and is not ready for use.

# Features
- MVC pattern
- Router (Gorilla)
- Middleware support
- Session managment (Gorilla)
- CSRF protection (Gorilla)
- Native view rendering (html/template)
- User authentication scaffolding
- Multi-language support
- Sample layout with login forms (Bulma CSS)

# External dependencies
External modules are mostly used when the feature is too complex to build or maintain - Router, Secure cookies

# TODO Features
- CSRF
- Cache (Redis)
- DB abstraction
- Validation

# Installation & Hot reloading
- Make sure you have Go >= 1.16 installed
- Configure `.env` based on `.env.example`
- Install CompileDaemon for running the watcher (https://github.com/githubnemo/CompileDaemon)
- Run the watcher script `./watch.sh`
- Open http://127.0.0.1:8888

# Licence
This project is licensed under the MIT License.
