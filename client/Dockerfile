FROM golang:1.16-alpine
WORKDIR /client
COPY client.go ./
RUN go build -o client client.go 
#CMD /client/client server_cont 56789 100 1