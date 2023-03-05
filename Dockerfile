FROM golang:1.19
WORKDIR /src
COPY . . 
RUN go build -o /bin/app .
ENTRYPOINT ["/bin/app"]