FROM golang:1.20.2 as builder
WORKDIR /app
ADD . .
ENV GOPROXY="https://goproxy.cn"
ENV CGO_ENBLED 1
ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH}:/app/lib
RUN go mod tidy && go build -v -a -ldflags="-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}" -installsuffix cgo -o msg_audit

MAINTAINER go-laoji
FROM alpine:3.17.2
WORKDIR /app
COPY --from=builder /app/lib/ ./lib/
COPY --from=builder /app/msg_audit /app/.env  ./
ENV LD_LIBRARY_PATH ${LD_LIBRARY_PATH}:/app/lib
ENV TZ Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
  apk --no-cache add tzdata ca-certificates libc6-compat gcompat libgcc libstdc++  &&\
  ln -s /lib/ld-musl-x86_64.so.1 /lib/libresolv.so.2 &&\
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone
VOLUME /app/logs
VOLUME /app/cache.db
EXPOSE 8080
CMD ["/app/msg_audit"]