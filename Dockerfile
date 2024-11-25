FROM golang:1.23 as buildbase

COPY . .

RUN CGO_ENABLED=0 go build -o /client .

FROM gcr.io/distroless/static-debian12

COPY --from=buildbase /client .

ENTRYPOINT ["./client"]