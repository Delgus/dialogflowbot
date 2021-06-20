FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY dialogflowbot . 
EXPOSE 80
CMD ["./dialogflowbot"]