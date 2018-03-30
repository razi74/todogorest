FROM golang:latest
WORKDIR /app
COPY bin/main /app/
EXPOSE 3000
CMD ["./main"]

