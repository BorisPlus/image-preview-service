
# Fat image with sources
FROM golang:1.21 as development

ARG MICROSERVICE
ENV SRC_DIR /go/src

# All source-code files
WORKDIR ${SRC_DIR}
COPY cmd cmd
COPY core core
COPY docker.integration-test docker.integration-test
COPY go.mod .

# Go requirements
RUN go mod tidy
RUN go mod download -x

# Special compile with out additional info
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /opt/service/service ./cmd/${MICROSERVICE}/

# Slim image with binary file
FROM alpine:3.9 as production

ARG MICROSERVICE

LABEL ORGANIZATION="OTUS Online Education"
LABEL MODULE="${MICROSERVICE}"
LABEL MAINTAINERS="admin@localhost"

COPY --from=development /opt/service/service /opt/service/service
COPY docker.integration-test/config/${MICROSERVICE}.json "/etc/service/config.json"

# Common log directory
RUN mkdir -p /var/log/previewer
RUN chown 1000:1000 /var/log/previewer

# One imege - one process
ENTRYPOINT [ "/opt/service/service", "-config", "/etc/service/config.json" ]

# docker build --no-cache --build-arg MICROSERVICE=http_service --progress plain -f ./docker/Dockerfile --tag=int-test-previewer-http_service .
# docker build --no-cache --build-arg MICROSERVICE=storage_service --progress plain -f ./docker/Dockerfile --tag=int-test-previewer-storage_service .
