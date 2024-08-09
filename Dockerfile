FROM golang:alpine

WORKDIR /app

COPY /app/go.mod /app/go.sum ./
RUN go mod download

COPY /app/. ./

RUN cd cmd/app && go build -o server

EXPOSE 8080

RUN chmod +x cmd/app/server

CMD ["cmd/app/server"]
