FROM golang:1.16-alpine
WORKDIR /app
COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o main .
CMD [ "/app/main" ]