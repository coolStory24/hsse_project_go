# Booking service

В корневой папке (`booking_service`) должен быть файл `.env`, содержащий переменные: 
- `DB_URL`
- `DB_USERNAME`
- `DB_PASSWORD`
- 'hotel_service_url' - url для gRPC с hotel_service, например `localhost:50051`

Если указана переменная среды `GO_ENV=dev`, то данные будут браться из файла `.env.dev`.

