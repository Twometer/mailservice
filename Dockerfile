FROM golang:latest
WORKDIR /app

COPY . .
RUN cd mailservice && go build -o ../svcmain

CMD "./svcmain"
