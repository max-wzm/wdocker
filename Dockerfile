FROM ubuntu:14.04.4
RUN apt-get update && apt-get install -y \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
		curl \
		git \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.20
ENV GOLANG_DOWNLOAD_URL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -f --insecure "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH /usr/local/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

COPY go-wrapper /usr/local/bin/
RUN chmod +x /usr/local/bin/go-wrapper

WORKDIR /go/src
COPY . .

RUN go-wrapper download