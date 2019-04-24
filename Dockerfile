# This file is a template, and might need editing before it works on your project.
FROM golang:1.11
#设置时区
ENV TZ Asia/Shanghai

RUN mkdir -p /go/src/go_cache /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /go

RUN sed -i -e 's/deb.debian.org/mirrors.163.com\//g' /etc/apt/sources.list
RUN apt update
RUN apt install net-tools


#复制项目
COPY go_cache .
RUN chmod +x go_cache
# RUN mkdir -p /go/data
CMD ["./go_cache"]
