
# Changelog
All notable changes to this project will be documented in this file.

## Upcoming...

## v0.7.0 - 2024-12-27

- Go version bump to >=1.23

## v0.6.0 - 2021-05-04

- Tests refactor

## v0.5.0 - 2021-04-22

Breaking changes:
- StoreToContext middleware split to: StoreUserToContext and StoreParamsToContext
- Session.SetUserLocale renamed to Session.SetLocale
- Session.GetUserLocale renamed to Session.GetLocale

Better UI on example app

## v0.4.0 - 2021-04-20

Breaking changes:
- Configurable session name added, gotuna.NewSession signature changed

## v0.3.0 - 2021-04-19

Breaking changes:
- StoreUserToContext middleware renamed to StoreToContext

- NewMuxRouter constructor added for the underlying mux.Router
- ContextWithParams / GetParam added for easier input data retrieval

