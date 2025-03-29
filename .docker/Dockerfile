FROM golang:1.24 AS build-stage

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o working_hours

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /
COPY --from=build-stage /app/working_hours /
COPY ./.env* /
EXPOSE 3333
USER nonroot:nonroot

CMD ["./app", ""]
