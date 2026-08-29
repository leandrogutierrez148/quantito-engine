FROM golang:1.21.0-alpine3.19

LABEL AUTHOR="Leandro Gutierrez <leandrogutierrez148@gmail.com>" \
    image="quantito-engine:deploy"

EXPOSE 8080

WORKDIR /app

COPY quantito-engine /app/quantito-engine

CMD ["/app/quantito-engine"]