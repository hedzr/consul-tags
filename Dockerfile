FROM golang:alpine
MAINTAINER hedzr

#ARG branch
#ARG BUILD_DATE
#ARG VCS_REF
#LABEL com.hedzr.build-date=$BUILD_DATE \
#      com.hedzr.vcs-ref=$VCS_REF \
#      com.hedzr.branch=$branch

ADD . /go/src/github.com/hedzr/consul-tags
#RUN apk --no-cache add bash go bzr git mercurial subversion openssh-client ca-certificates
RUN cd /go/src/github.com/hedzr/consul-tags \
 && ( ls -la; echo; echo "SOURCE_BRANCH = $SOURCE_BRANCH"; echo; env|sort; ) \
 && apk --update add --virtual .build-dependencies bash curl git ca-certificates \
 && ( [ "$SOURCE_BRANCH" != "master" ] && git checkout $SOURCE_BRANCH; ) \
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
