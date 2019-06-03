FROM golang:1.12
WORKDIR /var/server
COPY . .
RUN dep ensure
RUN go build ./cmd/serverd/serverd.go
EXPOSE 3000
CMD ['serverd']
