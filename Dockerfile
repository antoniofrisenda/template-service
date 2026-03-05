FROM scratch

WORKDIR /app
COPY bin/app .

EXPOSE 3000
CMD ["./app"]