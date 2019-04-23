FROM golang:1.12 AS build

RUN go get github.com/google/uuid && \
    go get github.com/gorilla/mux

ADD src /go/src

RUN go install buildnumber

FROM alpine:3.9 AS runtime
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    addgroup -S app && adduser -S -G app app

WORKDIR /app
ENV HOME /app
EXPOSE 8080

COPY --from=build /go/bin/buildnumber ./

USER app

ENTRYPOINT /app/buildnumber
