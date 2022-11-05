FROM golang:1.19.2-alpine

WORKDIR /crypto_quant_v2

COPY . /crypto_quant_v2

RUN cd /crypto_quant_v2 && go build

EXPOSE 8080

ENTRYPOINT ./CryptoQuant-v2