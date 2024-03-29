FROM golang:1.14.1
ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn"
WORKDIR $GOPATH/src/github.com/asynccnu/ele_service_v2
COPY . $GOPATH/src/github.com/asynccnu/ele_service_v2

RUN apt-get -qq update && apt-get -qq install -y --no-install-recommends ca-certificates curl # to solve x509(https)

RUN make
EXPOSE 8080
CMD ["./main", "-c", "conf/config.yaml"]
