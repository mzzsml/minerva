#FROM golang:latest
FROM golang:1.22
WORKDIR /app
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download
COPY *.go ./
COPY asset/ ./asset/
COPY module/ ./module/
COPY ui/ ./ui/
RUN CGO_ENABLED=0 GOOS=linux go build -o /minerva
EXPOSE 9999
CMD ["/minerva"]