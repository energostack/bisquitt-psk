FROM golang:1.22-alpine as builder

WORKDIR /app

ARG TEST
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go test -o e2e_test -c /app/tests

FROM alpine
RUN apk --no-cache add curl
COPY --from=builder /app/api /
COPY --from=builder /app/e2e_test /
COPY --from=builder /app/docs/swagger.json /docs/swagger.json

CMD if [ "$TEST" = "1" ]; then \
        /e2e_test; \
        # Check the return code
        if [ $? -eq 0 ]; then \
            echo "Tests passed"; \
            exit 0; \
        else \
            echo "Tests failed"; \
            exit 1; \
        fi; \
    else \
        /api; \
    fi
