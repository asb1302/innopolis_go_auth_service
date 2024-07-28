# innopolis_go_auth_service

## task

Добавить телеграмм бота для входа в сервис через ввод кода, отправленного сервисом авторизации в телеграмм-бота.
Для этого вам понадобятся эндпоинты:

* привязки телеграмм бота к авторизованному аккаунту;
* входа с логином и отправкой кода в телеграмм-чат;
* подтверждения кода из телеграмм чата. После успешной авторизации, код должен быть сброшен.

# Инструкция по работе с приложением

## Актуальные данные

* Актуальный адрес прода: https://innopolisgoauthservice-production.up.railway.app/ping
* Актуальный адрес телеграм бота: [@innopolisGoAuthServiceAsb1302Bot](https://t.me/innopolisGoAuthServiceAsb1302Bot)

## Регистрация нового пользователя

```http request
POST https://innopolisgoauthservice-production.up.railway.app/sign_up
Content-Type: application/json

{
  "login": "<укажите логин пользователя>",
  "password": "<укажите пароль>"
}
```

```curl
curl -X POST https://innopolisgoauthservice-production.up.railway.app/sign_up \
  -H "Content-Type: application/json" \
  -d '{
  "login": "<укажите логин пользователя>",
  "password": "<укажите пароль>"
}'
```

## Вход в систему

```http request
POST https://innopolisgoauthservice-production.up.railway.app/sign_in
Content-Type: application/json

{
  "login": "<укажите логин пользователя из ответа>",
  "password": "<укажите пароль>"
}
```

```curl
curl -X POST https://innopolisgoauthservice-production.up.railway.app/sign_in \
  -H "Content-Type: application/json" \
  -d '{
  "login": "<укажите логин пользователя из ответа>",
  "password": "<укажите пароль>"
}'
```

## Привязка Телеграм-бота к аккаунту

```http request
POST https://innopolisgoauthservice-production.up.railway.app/bind_telegram
User-ID: <укажите user id из ответа эндпоинта /sign_in>
Authorization: <укажите токен эндпоинта /sign_in>
Content-Type: application/json

{
  "telegram_username": "<укажите логин телеграм без @>"
}
```
```curl
curl -X POST https://innopolisgoauthservice-production.up.railway.app/bind_telegram \
  -H "User-ID: <укажите user id из ответа эндпоинта /sign_in>" \
  -H "Authorization: <укажите токен эндпоинта /sign_in>" \
  -H "Content-Type: application/json" \
  -d '{
  "telegram_username": "<укажите логин телеграм без @>"
}'
```

## Вход с логином и отправкой кода в Телеграм-чат

```http request
POST https://innopolisgoauthservice-production.up.railway.app/login_with_telegram
Content-Type: application/json

{
  "telegram_username": "<укажите логин телеграм без @>"
}
```

```curl
curl -X POST https://innopolisgoauthservice-production.up.railway.app/login_with_telegram \
  -H "Content-Type: application/json" \
  -d '{
  "telegram_username": "<укажите логин телеграм без @>"
}'
```

## Вход с кодом из телеграм чата

```http request
POST https://innopolisgoauthservice-production.up.railway.app/confirm_telegram_code
Content-Type: application/json

{
  "telegram_username": "<укажите логин телеграм без @>",
  "code": "<код из сообщения бота>"
}
```
```curl
curl -X POST https://innopolisgoauthservice-production.up.railway.app/confirm_telegram_code \
  -H "Content-Type: application/json" \
  -d '{
  "telegram_username": "<укажите логин телеграм без @>",
  "code": "<код из сообщения бота>"
}'
```

## Примечание

В процессе работы всегда можно получить актуальную информацию о пользователе с помощью эндпоинта:

```http request
GET https://innopolisgoauthservice-production.up.railway.app/get_user_info
Content-Type: application/json
User-ID: <укажите user id эндпоинта /sign_in>
Authorization: <укажите токен эндпоинта /sign_in>
```

```curl
curl -X GET https://innopolisgoauthservice-production.up.railway.app/get_user_info \
  -H "Content-Type: application/json" \
  -H "User-ID: <укажите user id эндпоинта /sign_in>" \
  -H "Authorization: <укажите токен эндпоинта /sign_in>"
```

# Локальная работа с проектом

## Требования

- Docker
- Docker Compose
- Make

## Основные команды Makefile

### Сборка и запуск проекта

Для сборки и запуска проекта выполните следующую команду:

```sh
make
```

### Очистка проекта

Для остановки и удаления контейнеров, образов, томов и зависших контейнеров выполните следующую команду:

```sh
make clean
```

### Проверка доступности сервиса

Для проверки доступности сервиса выполните следующую команду:

```sh
make check
```

### Примечание

Убедитесь, что вы настроили переменные окружения(.env) перед запуском команд.