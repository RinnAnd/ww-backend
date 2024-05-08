FROM golang:latest
LABEL maintainer="Andrés Casas <andre.ksas@gmail.com>"
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
ENTRYPOINT [ "/app/main" ]

# Create a docker-compose.yml file for this to be able to run the application