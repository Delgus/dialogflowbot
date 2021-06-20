FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY dialogflowbot . 
COPY wsapp/web wsapp/web
EXPOSE 80
CMD ["./dialogflowbot"]