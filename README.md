# traefik-cn-foward-auth

[![goreference](https://pkg.go.dev/badge/github.com/fopina/traefik-cn-foward-auth.svg)](https://pkg.go.dev/github.com/fopina/traefik-cn-foward-auth)
[![release](https://img.shields.io/github/v/release/fopina/traefik-cn-foward-auth)](https://github.com/fopina/traefik-cn-foward-auth/releases)
[![downloads](https://img.shields.io/github/downloads/fopina/traefik-cn-foward-auth/total.svg)](https://github.com/fopina/traefik-cn-foward-auth/releases)
[![ci](https://github.com/fopina/traefik-cn-foward-auth/actions/workflows/publish-main.yml/badge.svg)](https://github.com/fopina/traefik-cn-foward-auth/actions/workflows/publish-main.yml)
[![test](https://github.com/fopina/traefik-cn-foward-auth/actions/workflows/test.yml/badge.svg)](https://github.com/fopina/traefik-cn-foward-auth/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/fopina/traefik-cn-foward-auth/graph/badge.svg)](https://codecov.io/github/fopina/traefik-cn-foward-auth)

Simple HTTP service to validate TLS Info header using Forward Auth middleware.

While originally developed for whitelisting client certificates `commonName`s, it can also be used to whitelist any custom header values (via config)

> Note:  
> Keep in mind to always *set* these headers in a previous middleware, otherwise they can be set externally...

## Usage

...


## Build

Check out [CONTRIBUTING.md](CONTRIBUTING.md)
