FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY ./main .

COPY .env .

EXPOSE 8080

CMD ["./main"]

# docker build -t grs2080wdock/worker-knotew:latest . ; docker push grs2080wdock/worker-knotew:latest