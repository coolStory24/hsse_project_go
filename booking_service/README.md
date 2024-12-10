# Booking service

В корневой папке (`booking_service`) должен быть файл `.env`, содержащий переменные: 
- `DB_URL`
- `DB_USERNAME`
- `DB_PASSWORD`
- 'hotel_service_url' - url для gRPC с hotel_service, например `localhost:50051`
- 'user_service_url' - url для gRPC с user_service (для получения информации о клиентах для отправки им уведомлений), например `localhost:50052`

Если указана переменная среды `GO_ENV=dev`, то данные будут браться из файла `.env.dev`.

