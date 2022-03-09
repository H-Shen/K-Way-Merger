FROM golang:latest
MAINTAINER Haohu Shen
ADD . /KWayMerger
WORKDIR /KWayMerger
ENTRYPOINT ["go", "test", "./test", "-v"]