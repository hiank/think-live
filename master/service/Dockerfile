FROM alpine:latest
ADD master-service /
EXPOSE 8022 8026 8028
ENTRYPOINT ["/master-service", "--logtostderr"]