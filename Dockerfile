FROM alpine:3.20

WORKDIR /app

COPY food-app .

RUN chmod +x food-app

EXPOSE 8080

ENTRYPOINT ["./food-app"]
