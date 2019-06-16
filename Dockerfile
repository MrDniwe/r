FROM golang:stretch AS builder
RUN apt-get install curl -y && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/mrdniwe/r
ENV PATH="/go/src/github.com/mrdniwe/r:${PATH}"
COPY . .
RUN dep ensure
RUN go build serverd.go

FROM alpine 
WORKDIR /go/bin
COPY --from=builder /go/src/github.com/mrdniwe/r/serverd .
COPY --from=builder /go/src/github.com/mrdniwe/r/template .
ENTRYPOINT [ "/go/bin/serverd" ]
