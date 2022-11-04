FROM golang:1.19.2-alpine

WORKDIR /crypto_quant

COPY . /crypto_quant

RUN cd /crypto_quant && go build

EXPOSE 8080

ENTRYPOINT ./cryptoQuant