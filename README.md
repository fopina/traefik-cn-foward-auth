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

docker image [ghcr.io/fopina/traefik-cn-foward-auth](https://github.com/fopina/traefik-cn-foward-auth/pkgs/container/traefik-cn-foward-auth) or binaries in [releases](https://github.com/fopina/traefik-cn-foward-auth/releases) available.

```
➜  traefik-cn-foward-auth -h
Run HTTP service that validates headers

Usage:
  traefik-cn-foward-auth [flags]

Flags:
  -a, --allow-header string             Name of the header that will container a list of allowed common names - or raw values, if --raw (default "X-Allow-CN")
  -s, --allow-header-separator string   Separator character of the values in allow-header. Use special value "json" if you prefer to use a JSON string array to specify the list (default ",")
  -b, --bind-addr string                Address to bind the web server to (default ":8080")
      --debug                           Log all failed validations to help debug
  -n, --header string                   Name of the header with values to be validated (default "X-Forwarded-Tls-Client-Cert-Info")
  -h, --help                            help for traefik-cn-foward-auth
      --raw                             By default, values in --allow-header are expected to be the common names, eg: for "CN=mobile01,OU=..." it should have "mobile01". Using --raw no such parsing is done and allowed values are expected to match exactly the value sent in --header
  -v, --version                         version for traefik-cn-foward-auth
```

### CommonName required

This is to use as a [forwardauth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) Traefik middleware, chained with [passtlsclientcert](https://doc.traefik.io/traefik/middlewares/http/passtlsclientcert/) and [headers](https://doc.traefik.io/traefik/middlewares/http/headers/).

* passtlsclientcert: add Subject CN to an header
* headers: add allowed CNs to another header
* forwardauth: call this service and only proceed if it returns 200 => Subject CN is part of allowed CNs

Full working example available in [example/docker-compose.yml](example/docker-compose.yml) (`whoamix-cn` router)

```yaml
services:
  traefik:
    image: traefik:v2.9
    ...
    command:
      ...
      - --entrypoints.webtls.address=:443
      - --entrypoints.webtls.http.tls.options=mtlsTest@file
    ...

  whoami:
    image: traefik/whoami
    labels:
      traefik.http.services.whoami.loadbalancer.server.port: 80
      traefik.http.middlewares.pass-certificate-info.passTLSClientCert.info.subject.commonName: 'true'
      traefik.http.middlewares.allow-cns.headers.customrequestheaders.X-Allow-CN: auth-client
      traefik.http.middlewares.validate-cn.forwardauth.address: http://validate-cn:8080/
      traefik.http.routers.whoamix-cn.rule: Host(`whoami-mtls.7f000001.nip.io`)
      traefik.http.routers.whoamix-cn.entrypoints: webtls
      traefik.http.routers.whoamix-cn.tls: "true"
      traefik.http.routers.whoamix-cn.tls.options: mtlsTest@file
      traefik.http.routers.whoamix-cn.middlewares: pass-certificate-info,allow-cns,validate-cn
      ...
  
  validate-cn:
    image: ghcr.io/fopina/traefik-cn-foward-auth
```

This only allows clients with CN `auth-client` or `auth2-client`.

### Headers required

`--raw` mode matches allowed header values directly with header value, allowing any other use cases for authorizing based on specific headers

```yaml
services:
  traefik:
    image: traefik:v2.9
    ...
    command:
      ...
      - --entrypoints.webtls.address=:443
      - --entrypoints.webtls.http.tls.options=mtlsTest@file
    ...

  whoami:
    image: traefik/whoami
    labels:
      traefik.http.services.whoami.loadbalancer.server.port: 80
      traefik.http.middlewares.allow-header.headers.customrequestheaders.X-Allowed-Passphrase: let me in,me me me
      traefik.http.middlewares.validate-raw.forwardauth.address: http://validate-raw:8080/
      traefik.http.routers.whoamix-raw.rule: Host(`whoami-mtls-raw.7f000001.nip.io`)
      traefik.http.routers.whoamix-raw.entrypoints: webtls
      traefik.http.routers.whoamix-raw.tls: "true"
      traefik.http.routers.whoamix-raw.tls.options: mtlsTest@file
      traefik.http.routers.whoamix-raw.middlewares: allow-header,validate-raw
      ...


  validate-raw:
    image: ghcr.io/fopina/traefik-cn-foward-auth
    command: --raw --allow-header "X-Allowed-Passphrase" --header "X-Passphrase"
```

## Build

Check out [CONTRIBUTING.md](CONTRIBUTING.md)
