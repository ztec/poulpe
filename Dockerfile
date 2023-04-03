FROM golang:1.20 as builder

ARG BUILD_MODE=prod

WORKDIR /app
RUN adduser --disabled-password --gecos "" appuser

ADD ./ /app

RUN MODE=${BUILD_MODE} bin/build.sh


FROM debian:stable-slim as server
RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/poulpe /app/

# allow to use appuser to execute the app
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 8080

USER appuser
CMD ["/app/poulpe", "server"]

FROM server as server-preloaded
RUN /app/poulpe search love