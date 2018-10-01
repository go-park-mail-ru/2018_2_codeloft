FROM golang:1.10.3-alpine

COPY . /go/src/2018_2_codeloft

RUN go install 2018_2_codeloft

EXPOSE 8080

CMD ["2018_2_codeloft"]

