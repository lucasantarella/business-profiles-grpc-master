# syntax=docker/dockerfile:experimental
FROM golang:latest as builder
WORKDIR /app
RUN mkdir -p -m 0600 ~/.ssh && \
    ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts && \
    ssh-keygen -F github.com -l -E sha256 \
        | grep -q "SHA256:nThbg6kXUpJWGl7E1IGOCspRomTxdCARLviKw6E5SY8"
COPY go.mod go.sum ./
RUN git config --system url."ssh://git@github.com/".insteadOf "https://github.com/"
RUN --mount=type=ssh mkdir -p /var/ssh && \
    GIT_SSH_COMMAND="ssh -o \"ControlMaster auto\" -o \"ControlPersist 300\" -o \"ControlPath /var/ssh/%r@%h:%p\"" \
    go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./gateway.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/gateway .
EXPOSE 50051
CMD ["./main"]