FROM golang:alpine AS builder

LABEL maintainers="Anggit M Ginanjar <anggit@isi.co.id>"
# RUN apk update && apk --no-cache add git
WORKDIR $GOPATH/src/
COPY . .

RUN ls -alh
# RUN go get -v -d
RUN go mod download
RUN GOOS=linux go build -o /go/bin/ory-kratos-workshop main.go
COPY views /go/bin/views
RUN ls -alh /go/bin

FROM alpine
COPY --from=builder /go/bin/ /srv/
WORKDIR /srv/
RUN ls -alh
CMD ./ory-kratos-workshop
