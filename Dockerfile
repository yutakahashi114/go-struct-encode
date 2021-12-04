FROM golang:1.17

RUN go get -u github.com/gogo/protobuf/proto && \
    go install github.com/gogo/protobuf/protoc-gen-gogofaster@latest && \
    go get -u github.com/gogo/protobuf/gogoproto && \
    go install github.com/gogo/protobuf/protoc-gen-gofast@latest

RUN apt-get update && apt-get install -y unzip

RUN mkdir -p /tmp/protoc && \
    curl -sSL "https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip" > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip  && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    chmod go+rx /usr/local/bin/protoc && \
    cp -r /tmp/protoc/include/* /usr/local/include/ && \
    cd /tmp && \
    rm -r /tmp/protoc

# RUN apt-get update && \
#     apt-get install autoconf automake libtool && \
#     curl -L -o /tmp/protobuf.tar.gz https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protobuf-cpp-3.19.1.tar.gz && \
#     cd /tmp && \
#     tar xvzf protobuf.tar.gz && \
#     cd protobuf-3.19.1/ && \
#     ./autogen.sh && \
#     ./configure && \
#     make -j 3 && \
