FROM golang:latest
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download && \
    go-wrapper install
CMD ["go-wrapper", "run", "-webhook"]
EXPOSE 8443