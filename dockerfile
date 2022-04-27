FROM golang:buster AS build 

WORKDIR /app

# use modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/kubeclean .


FROM gcr.io/distroless/base

EXPOSE 8080

COPY --from=build /go/bin/kubeclean /go/bin/kubeclean

ENTRYPOINT ["/go/bin/kubeclean","serve"]