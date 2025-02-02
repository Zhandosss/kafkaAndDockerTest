## Задание для тестирования kafka и docker-compose



### Эндпоинты

Swagger документация доступна по пути `/swagger/index.html`

- `POST api/messages` - отправка сообщения в очередь
- `GET api/messages` - получение сообщений из базы данных
- `GET api/messages/{id}` - получение сообщения по id
- `DELETE api/messages` - удаление всех сообщений из базы данных
- `DELETE api/messages/{id}` - удаление сообщения по id
- `GET api/statistic/days` - получение статистики по сообщениям


### Используемые библиотеки
- `github.com/labstack/echo/v4` - для роутинга
- `github.com/jmoiron/sqlx` - для работы с базой данных
- `github.com/IBM/sarama` - для работы с kafka
- `github.com/swaggo/echo-swagger` - для генерации swagger документации
-  `github.com/rs/zerolog/log` - для логирования
- `github.com/pressly/goose/v3` - для миграций базы данных

### Деплоймент
Для деплоя используется docker-compose. Для запуска необходимо выполнить команду:
```bash
docker-compose up -d
```

Для остановки контейнеров:
```bash
docker-compose down
```

В приложений используется следующие контейнеры:
- `consumer` - приложение, которое получает сообщения из очереди и помечает их как прочитанные
- `producer` - приложение, которое является api-gateway для отправки сообщений в очередь и REST API для отправки и получения информации о сообщениях, а также для получения статистики о сообщениях

Также в проекте используется `kafka` для обмена сообщениями между приложениями и `postgres` в качестве базы данных.
