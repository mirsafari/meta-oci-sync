FROM golang:1.21 AS build-stage

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /controller

# Second stage - minimal image
FROM scratch
WORKDIR /
COPY --from=build-stage /controller /controller
ENTRYPOINT ["/controller"]