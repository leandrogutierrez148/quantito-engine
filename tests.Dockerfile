# docker build -f Dockerfile --tag quantito-engine:tests .
FROM golang:1.21.0-alpine3.18

WORKDIR /build

RUN apk update && \
    apk add gcc musl-dev git

ENV GOPATH=/usr/local/.go CGO_ENABLED=1

RUN go install github.com/jstemmer/go-junit-report/v2@latest
RUN go install github.com/boumenot/gocover-cobertura@latest
RUN go install github.com/vektra/mockery/v2@v2.33.2

ENV GOPATH=/build/.go CGO_ENABLED=1

CMD /usr/local/.go/bin/mockery && go test -race -v 2>&1  ./... -tags musl -coverprofile=coverage.out -json | tee test-report.out && /usr/local/.go/bin/gocover-cobertura < coverage.out > coverage.xml && /usr/local/.go/bin/go-junit-report -in test-report.out -set-exit-code > report.xml 
