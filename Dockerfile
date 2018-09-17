FROM golang:1.11-alpine as builder
COPY ./app /go/src/go-ex/app
WORKDIR /go/src/go-ex/app
RUN apk add --update git curl \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    && dep ensure \ 
    && go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates \
    && apk add --update tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && apk del tzdata
WORKDIR /root/
COPY --from=builder /go/src/go-ex/app/main /root
COPY --from=builder /go/src/go-ex/app/config/ /root/config
EXPOSE 1234
CMD ["./main"]