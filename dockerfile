FROM golang

##buat folder APP
RUN mkdir /app

##set direktori utama
WORKDIR /app

##copy seluruh file ke app
ADD . /app

##buat executeable
RUN go build -o cmd/api/main.go

EXPOSE 8080

##jalankan executeable
CMD ["/app/main"]