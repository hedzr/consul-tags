FROM golang:alpine
MAINTAINER hedzr

ADD . /go/src/github.com/hedzr/consul-tags
#RUN apk --no-cache add bash go bzr git mercurial subversion openssh-client ca-certificates
RUN cd /go/src/github.com/hedzr/consul-tags \
 && ls -la \
 && apk --update add --virtual .build-dependencies bash curl git ca-certificates \
 && git checkout devel \
 && go get -u \
 && go build -v . \
 && mv ./consul-tags /go/bin/ \
 && apk del .build-dependencies \
 && rm -rf /var/cache/apk/*

#ONBUILD COPY . /go/src/app
#ONBUILD RUN go-wrapper download
#ONBUILD RUN go-wrapper install

CMD ["-h"]
ENTRYPOINT ["consul-tags"]
