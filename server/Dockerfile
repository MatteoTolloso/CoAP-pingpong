FROM golang:1.16-alpine
WORKDIR /server
COPY server.go ./
RUN go build -o server server.go 
EXPOSE 56789/udp
#CMD /server/server 56789