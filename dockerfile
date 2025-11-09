FROM node:20-slim

WORKDIR /app

COPY ./main .

COPY .env .

EXPOSE 8080

CMD ["./main"]