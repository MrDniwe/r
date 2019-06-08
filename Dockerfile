FROM golang:1.12
RUN apt-get install curl -y && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/mrdniwe/r
ENV PATH="/go/src/github.com/mrdniwe/r:${PATH}"
COPY . .
RUN dep ensure
RUN go build serverd.go
EXPOSE 3000
ENTRYPOINT serverd
