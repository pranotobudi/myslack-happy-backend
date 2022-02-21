# --- Section 1 Start
FROM golang:1.16-alpine as build-dev

RUN mkdir /app
COPY . /app
WORKDIR /app

# should this line enabled?? yes, I've tested it. 
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
# --- Section 1 End

# --- Section 2 Start
# this line: no shell, no static image, minimal size as possible
# FROM gcr.io/distroless/static

# COPY --from=build-dev /app/main .
# --- Section 2 End

EXPOSE 8080
CMD ["./main"]