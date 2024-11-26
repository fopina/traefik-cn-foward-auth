FROM scratch

COPY traefik-cn-foward-auth /usr/bin/traefik-cn-foward-auth
ENTRYPOINT ["/usr/bin/traefik-cn-foward-auth"]
