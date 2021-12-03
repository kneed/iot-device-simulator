FROM golang:1.16.5-alpine3.13 AS build

RUN mkdir /app
WORKDIR /app
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download
ADD . /app
ENV GOPROXY=https://goproxy.cn,direct
RUN go build -o main .

#CMD ["/app/main"]
###
FROM golang:1.16.5-alpine3.13 as final
COPY --from=build /app .
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["./main"]