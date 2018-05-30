FROM golang:1.10.2-stretch as builder

ENV SRC_DIR /go/src/github.com/mirakl/http2back
WORKDIR $SRC_DIR

COPY cli/ ./cli
COPY notifier/ ./notifier
COPY provider/ ./provider
COPY server ./server
COPY Makefile ./
COPY glide.* ./
COPY *.go ./

RUN go get -u github.com/Masterminds/glide

RUN glide install

RUN make build

FROM centos:latest

COPY --from=builder /go/src/github.com/mirakl/http2back/bin/http2back /bin
RUN chmod +x /bin/http2back

EXPOSE 8080

USER nobody
ENTRYPOINT ["/bin/http2back"]
