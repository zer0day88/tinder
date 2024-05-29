
FROM golang:1.21 AS builder
RUN mkdir -p /app
ENV TZ=Asia/Jakarta
RUN GOCACHE=OFF
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY . /app
RUN cd /app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/cmd ./cmd/main.go


FROM alpine:latest
LABEL maintaner="andrekurniawan0891@gmail.com"
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /go/bin/cmd .
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY config/docker-config.yml config/app-config.yml
#COPY .env .
ENV TZ=Asia/Jakarta

ENTRYPOINT ["/root/cmd"]
EXPOSE 3000