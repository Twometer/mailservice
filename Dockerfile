FROM golang:latest
WORKDIR /app

COPY . .
RUN cd mailservice && go build -o ../mailservice

CMD ["./mailservice"]
