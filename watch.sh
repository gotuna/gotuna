#!/bin/bash

# hot-reload golang projects (https://github.com/githubnemo/CompileDaemon)
CompileDaemon -exclude-dir=.git -include="*.html" -command "./main" -build="go build -o main cmd/main/main.go"

