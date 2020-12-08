FROM golang:buster

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

ENV MMDBINSPECT_VERSION 0.1.1
ENV LESSCHARSET utf-8

COPY . /project/

WORKDIR /tmp

RUN curl --location https://github.com/maxmind/mmdbinspect/releases/download/v0.1.1/mmdbinspect_${MMDBINSPECT_VERSION}_linux_amd64.deb -o mmdbinspect.deb && dpkg -i /tmp/mmdbinspect.deb

RUN apt-get update && apt-get install less

WORKDIR /project

RUN go build
RUN ./mmdb-from-go-blogpost
