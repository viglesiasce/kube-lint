FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go build -o kube-lint .
FROM alpine
RUN adduser -S -D -H -h /app kube-lint
USER kube-lint
COPY --from=builder /build/kube-lint /app/
WORKDIR /app
CMD ["./kube-lint"]
