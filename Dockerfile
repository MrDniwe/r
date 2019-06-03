FROM golang:1.12
RUN apt-get install curl -y && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /var/server
COPY . .
RUN dep ensure
RUN go build ./cmd/serverd/serverd.go
EXPOSE 3000
CMD ['serverd']
