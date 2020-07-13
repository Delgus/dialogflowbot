FROM alpine:3.12
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY bot bot
EXPOSE 80
CMD ["./bot"]