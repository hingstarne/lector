FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build 
RUN go build -o lector .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/lector /app/
WORKDIR /app
EXPOSE 8000
CMD ["./lector"]