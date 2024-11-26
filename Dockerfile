FROM scratch

COPY golang-template /usr/bin/traefik-cn-foward-auth
ENTRYPOINT ["/usr/bin/traefik-cn-foward-auth"]
