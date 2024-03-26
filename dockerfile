FROM debian:bookworm
COPY ./src/main /etc/appname/main
ENTRYPOINT ["/etc/appname/main"]