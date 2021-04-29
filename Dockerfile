FROM golang:1.14 AS artifact

# Building artifact
RUN apt install ca-certificates git
WORKDIR /app
COPY . ./
RUN ls
RUN go mod download
RUN CGO_ENABLED=0 go build -o hefesto github.com/ambientis-org/hefesto/cmd/server

# Running server
FROM alpine
WORKDIR /app
COPY --from=artifact /app/hefesto .
EXPOSE 8080

RUN chmod +x ./hefesto

CMD ["./hefesto"]