FROM golang:1.11.2 as builder

RUN go get -u github.com/golang/dep/cmd/dep

RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

WORKDIR /go/src/github.com/vishwanathj/golangexercises/go_web_server

ADD ./ /go/src/github.com/vishwanathj/protovnfdparser/

WORKDIR /go/src/github.com/vishwanathj/protovnfdparser/cmd/app

RUN dep ensure

RUN dep status
#RUN go get -d -v
#RUN go get gopkg.in/tomb.v2 github.com/stretchr/testify golang.org/x/lint/golint

# Compile Delve
#RUN apk add --no-cache git
#RUN go get github.com/derekparker/delve/cmd/dlv

#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -a -installsuffix cgo -o /go/bin/vnfdservice .


WORKDIR /go/src/github.com/vishwanathj/protovnfdparser/
#https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go
ADD ./pkg/schema/* /go/src/github.com/vishwanathj/protovnfdparser/pkg/schema/* /usr/share/vnfdservice/schema/
RUN go test ./... -coverprofile=coverage.out -tags=unit

FROM golang:1.11.2
#RUN apk --no-cache add ca-certificates
# Port 8080 belongs to our application, 40000 belongs to Delve
#EXPOSE 8080 40000

# Allow delve to run on Alpine based containers.
#RUN apk add --no-cache libc6-compat

#WORKDIR /go/src/github.com/vishwanathj/protovnfdparser/pkg/schema
COPY --from=builder /go/bin/* /go/bin/
#COPY --from=builder /go/pkg/* /go/pkg/
#COPY --from=builder /go/src/github.com/vishwanathj/protovnfdparser/pkg/schema/* /go/src/github.com/vishwanathj/protovnfdparser/pkg/schema/
COPY --from=builder /go/src/github.com/vishwanathj/protovnfdparser/pkg/schema/* /usr/share/vnfdservice/schema/
COPY --from=builder /go/src/github.com/vishwanathj/protovnfdparser/test/* /usr/share/vnfdservice/test/

RUN rm -rf /usr/share/vnfdservice/test/jmeter/
WORKDIR /usr/share/vnfdservice/test/
RUN apt-get update && apt-get install -y netcat dnsutils
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/ssl/certs/vnfdsvc.key -out /etc/ssl/certs/vnfdsvc.crt -subj /CN=vnfdsvc/O=vnfdsvc

# Run delve
#CMD ["dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/server"]