FROM alpine AS build
RUN addgroup -g 10001 app && adduser --disabled-password -u 10001 -G app -h /app app -s /bin/scan_health
RUN apk --no-cache add ca-certificates

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY /dist/scan_health-linux-amd64 /bin/scan_health
USER app
WORKDIR /app
ENTRYPOINT ["/bin/scan_health"]