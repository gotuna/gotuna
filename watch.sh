#!/bin/bash

# hot-reload golang projects (https://github.com/githubnemo/CompileDaemon)
CompileDaemon -exclude-dir=.git -include="*.html" -command "./webapp" -build="go build -o webapp cmd/main/main.go"

