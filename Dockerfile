FROM golang:1.13.13
ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn"
WORKDIR $GOPATH/src/github.com/asynccnu/ele_service_v2
COPY . $GOPATH/src/github.com/asynccnu/ele_service_v2
RUN make
EXPOSE 8080
CMD ["./main", "-c", "conf/config.yaml"]
