# Backend

## Установка

Для установки приложения вам потребуется:

- [Docker engine](https://docs.docker.com/engine/) (версия 24.x. или выше)
- [Docker compose](https://www.npmjs.com/) (версия 2.x. или выше)

```bash
git clone https://github.com/Fi44er/sdmed.git
cd ./sdmed
   ```
## Запуск

Запуск postgres и redis контейнеров

```bash
docker compose up
```

Запуск golang приложения

```bash
go run cmd/main.go
```

## Swagger

Доступ к документации можно получить по адресу: *http://127.0.0.1:8080/swagger/index.html*
