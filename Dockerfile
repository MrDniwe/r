FROM golang:stretch AS builder
RUN apt-get install curl -y && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
WORKDIR /go/src/github.com/mrdniwe/r
ENV PATH="/go/src/github.com/mrdniwe/r:${PATH}"
RUN go get -v github.com/rubenv/sql-migrate/...
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build serverd.go

FROM scratch 
WORKDIR /go/bin
COPY --from=builder /go/bin/sql-migrate .
COPY --from=builder /go/src/github.com/mrdniwe/r/migrations ./migrations
COPY --from=builder /go/src/github.com/mrdniwe/r/serverd .
COPY --from=builder /go/src/github.com/mrdniwe/r/default-entrypoint.sh .
COPY --from=builder /go/src/github.com/mrdniwe/r/template ./template
ENV PATH="/go/bin:${PATH}"
ENTRYPOINT [ "/go/bin/serverd" ]
