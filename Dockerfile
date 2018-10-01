FROM golang:1.10.3-alpine3.7 as builder

COPY . /go/src/2018_2_codeloft
#COPY ../github.com /go/src/github
#COPY ../golang.org /go/src/golang.org

RUN apk update && apk upgrade && \
    apk add --no-cache bash git
RUN go get github.com/icrowley/fake
RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/vk
#RUN go install 2018_2_codeloft
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .
RUN cd /go/src/2018_2_codeloft && go build -o goapp


FROM alpine
WORKDIR /app
COPY --from=builder /go/src/2018_2_codeloft/goapp /app/
EXPOSE 8080
ENTRYPOINT ./goapp

#CMD ["2018_2_codeloft"]

#FROM scratch
#COPY --from=builder /app ./
#EXPOSE 8080
#ENTRYPOINT ["./app"]


#CMD ["2018_2_codeloft"]

