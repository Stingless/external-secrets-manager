FROM golang:1.21 as golang
WORKDIR /app
COPY ./es-manager .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /es-manager .
FROM gcr.io/distroless/static-debian11
COPY --from=golang /es-manager .
COPY --from=golang /bin/bash .
COPY --from=golang /bin/curl .
COPY --from=golang /bin/sleep .
EXPOSE 6443
EXPOSE 8200
CMD ["/es-manager"]
