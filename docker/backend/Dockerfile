FROM golang:1.13.6

WORKDIR /src

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o /bin/api main.go

WORKDIR /
RUN rm -rf /src

CMD ["/bin/api"]
EXPOSE 8091