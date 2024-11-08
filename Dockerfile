FROM alpine AS build
RUN adduser --disabled-password -u 10001 app
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY /dist/scan_health-linux-amd64 /bin/scan_health
USER app
ENV HOME=/
ENTRYPOINT ["/bin/scan_health"]