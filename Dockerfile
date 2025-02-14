FROM golang:alpine AS builder

# 构建可执行文件
ENV CGO_ENABLED=0
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN go build -o main

FROM alpine
WORKDIR /app
COPY --from=builder /build/main /app
RUN apk add tzdata
CMD ["./main"]
