# Medods Authorization Service

Сервис который отвечает за аутентификацию пользователя.

* Rest Api with GIN
* Docker
* Mail service
* Swagger
* Migrations

## Installation

Сначала нужно скопировать репозиторий из Github:
```
git clone https://github.com/islombay/medods.git
```

Далее создаём файл `.env` и заполняем его как показано в тестовом файле `.env.dev`
```js
SERVER_HOST=0.0.0.0
SERVER_PORT=8095
SERVER_URL=0.0.0.0:8095

DB_HOST=postgres
DB_PORT=5432
DB_NAME=medods
DB_SSL_MODE=disable
DB_MIGRATIONS_PATH=migrations
DB_PWD=password
DB_USER=postgres

TOKEN_SECRET_KEY='secret key for generating tokens'
TOKEN_ACCESS_DURATION_MINUTES=5

MAIL_HOST='smtp.gmail.com'
MAIL_PORT=587
MAIL_PWD='gmail app password'
MAIL_SENDER='your email'
```

### **Важно:** если изменили хост, порт, и тд. то также измените настройки docker контейнера!


## Run
Можно использовать команду которая создает и билдит image.
```
docker compose up --build
```
а вне которых случаях команда такая:
```
docker-compose up --build
```

## Swagger

Документация swagger доступна по ссылке: https://0.0.0.0:8095/sw/index.html, если конфигурация в .env файле не был изменен.