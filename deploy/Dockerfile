FROM alpine:latest
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD cronhubot /
CMD ["/cronhubot"]