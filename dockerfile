FROM debian:bookworm
ARG TARGETARCH
COPY "./src/main_$TARGETARCH/main_$TARGETARCH" /etc/appname/main
RUN ["chmod", "+x", "/etc/appname/main"]
ENTRYPOINT ["/etc/appname/main"]