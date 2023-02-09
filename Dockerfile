FROM alpine
ENV HOME=/
COPY /dist/scan_health-linux-amd64 /bin/scan_health
ENTRYPOINT ["/bin/scan_health"]