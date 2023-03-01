FROM alpine
ENV HOME=/
COPY /dist/scan_health-linux-amd64 /bin/scan_health
RUN adduser -D app
USER app
ENTRYPOINT ["/bin/scan_health"]