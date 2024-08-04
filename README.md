# Микросервис обработки сообщений
## Задание:
Разработать микросервис на Go, который будет принимать сообщения через HTTP API, сохранять их в PostgreSQL, а затем отправлять в Kafka для дальнейшей обработки. Обработанные сообщения должны помечаться. Сервис должен также предоставлять API для получения статистики по обработанным сообщениям.
Требования:
1.	Использовать Go 1.20+
2.	Использовать PostgreSQL для хранения сообщений
3.	Реализовать отправку и чтение сообщений в Kafka
4.	Предоставить возможность запустить проект в Docker

## Требования к результату:

Мы ожидаем, что тестовое задание будет запущено на сервере и доступно для тестирования через интернет

На выходе ожидаем получить:
Ссылку на проект развернутый на сервере
Инструкцию по подключению
Git репозиторий с кодом

# API

## SSO(jwt auth)

### {addr}:{port}/register
```JSON
{
  "user":"username",  
  "pass":"password"
}
```

### {addr}:{port}/auth
```JSON
{
  "user":"username",
  "pass":"password"
}
```
this method returns jwt auth in cookies

## Messager

```JSON
{
  "id": uint,
  "pal_id": uint,
  "user_id": uint,
  "created_at": uint,
  "message": "your msg"
}
```

## POST {addr}:{port}/send

```JSON
{
  "pal_id": uint,
  "message": "your msg"
}
```

## GET {addr}:{port}/chat?pal_id={uint}&offset={uint}&limit={uint}

```json
[
  {
    "id": uint,
    "pal_id": uint,
    "user_id": uint,
    "created_at": uint,
    "message": "your msg"
  },
  {
    "id": uint,
    "pal_id": uint,
    "user_id": uint,
    "created_at": uint,
    "message": "your msg"
  },
  ...
]
```

## GET {addr}:{port}/chats

```json
[
  {
    "id": uint,
    "pal_id": uint,
    "user_id": uint,
    "created_at": uint,
    "message": "last sent message from [pal_id]"
  },
  {
    "id": uint,
    "pal_id": uint,
    "user_id": uint,
    "created_at": uint,
    "message": "last sent message from [pal_id]"
  },
  ...
]
```

## Statistics

## GET {addr}:{port}/MessagesReceived?since={uint}&to={uint}

```json
{
  "messages": uint,
  "since": uint,
  "to": uint,
}
```

## GET {addr}:{port}/MessagesReceivedByUser?since={uint}&to={uint}

```json
{
  "user_id": uint,
  "messages": uint,
  "since": uint,
  "to": uint,
}
```

## GET {addr}:{port}/MessagesSentByUser?since={uint}&to={uint}

```json
{
  "user_id": uint,
  "messages": uint,
  "since": uint,
  "to": uint,
}
```