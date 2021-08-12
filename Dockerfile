FROM golang:alpine as builder
LABEL maintainer="candraalims@gmail.com"

ENV GOBIN /go/bin
ENV GOPATH /app
ENV PATH=$GOPATH/bin:$PATH

RUN apk update && \
apk upgrade && \
apk add --no-cache git tzdata && \
cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
echo "Asia/Jakarta" > /etc/timezone

RUN mkdir -p /app/referral_service
ADD . /app/referral_service
WORKDIR /app/referral_service/cmd/app

RUN go get . && \
go mod tidy && \
go mod download && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o referral

FROM scratch

COPY --from=builder /app/referral_service/config.json .
COPY --from=builder /app/referral_service/cmd/app/referral .
COPY --from=builder /etc/localtime /etc/localtime
COPY --from=builder /etc/timezone /etc/timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 8080

CMD ["/referral"]

