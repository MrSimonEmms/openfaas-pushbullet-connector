FROM golang AS builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /opt/openfaas-pushbullet-connector
ADD . .
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build \
  -ldflags \
  "-X github.com/MrSimonEmms/openfaas-pushbullet-connector/cmd.Version=$VERSION -X github.com/MrSimonEmms/openfaas-pushbullet-connector/cmd.GitCommit=$GIT_COMMIT"
ENTRYPOINT [ "/opt/openfaas-pushbullet-connector/openfaas-pushbullet-connector" ]

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /opt/openfaas-pushbullet-connector/openfaas-pushbullet-connector /app
ENTRYPOINT [ "/app/openfaas-pushbullet-connector" ]
