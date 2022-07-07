FROM golang

RUN apt-get update
RUN apt-get install -y npm shellcheck chromium

ENV GO111MODULE=on
RUN go install golang.org/x/tools/cmd/goimports
RUN go install mvdan.cc/sh/v3/cmd/shfmt
RUN go install golang.org/x/tools/cmd/stringer
RUN go install golang.org/x/lint/golint

RUN npm --unsafe-perm=true install -g prettier
