# Запуск Postgres в контейнере

Для запуска Postgres в контейнере выполните:

```bash
make pg
```

Скрипты инициализации лежат в `db/init`. Файлы БД лежат в `db/data`.

Для того, чтобы убедиться, что БД была запущена корректно, можно посмотреть ее логи:

```bash
docker logs praktikum-webinar-db
```

# Создание таблиц и их наполнение

Код для создания таблиц и их наполнения случайными данными лежит в `app/cmd/datagen`.

Для запуска `datagen` выполните следующую комманду:

```bash
make build-datagen

./app/bin/datagen -dsn postgresql://gopher:gopher@localhost:5432/gopher_corp -emp-count 100000
```

На что стоит обратить внимание:
- работа с БД при помощи `database/sql`
- работа с транзакциями
- выполнение операции батчем
- получение ID вставленных строк
- обработка ошибок и совместимость версий `pgx`

# Работа с данными

Код приложения, работающего с БД, лежит в `app/cmd/employees`.

Для запуска сервера выполните следующую комманду:

```bash
make build-app

./app/bin/employees -dsn postgresql://gopher:gopher@localhost:5432/gopher_corp
```

Для получения результатов перейдите в бразуере, например, по ссылке:

`http://localhost:8080/employees/ann?limit=50&last-id=46`

На что стоит обратить внимание:
- PGX pool и его конфигурация
- Пагинация: https://use-the-index-luke.com/no-offset

## Индексы

Лучший источник информации об индексах: https://postgrespro.ru/docs/postgrespro/16/indexes

Подключимся к БД:

```bash
psql -h localhost -p 5432 -U gopher -d gopher_corp
```

Посмотрим на стоимость запроса поиска по фамилии без использования индексов:

```sql
EXPLAIN (ANALYZE, VERBOSE)
SELECT id, first_name, last_name, salary, position, email
FROM employees
WHERE
    id > 542 AND lower(last_name) LIKE 's%'
ORDER BY ID asc
LIMIT 1000
;
```

Создадим индекс для т.н. `fuzzy`-поиска:

```sql
CREATE INDEX
ON employees using btree(id, lower(last_name) text_pattern_ops);
```

Про класс операторов `text_pattern_ops` можно почитать здесь: https://postgrespro.ru/docs/postgrespro/16/indexes-opclass

Выполним анализ стоимости запроса после добавления индекса:

```sql
EXPLAIN (ANALYZE, VERBOSE)
SELECT id, first_name, last_name, salary, position, email
FROM employees
WHERE
    id > 542 AND lower(last_name) LIKE 's%'
ORDER BY ID asc
LIMIT 1000
;
```

Чтобы просмотреть мета-данные таблицы `employees` используйте `\d+ employees` (если вы используете `psql` -- это метакоманда, детали здесь: https://postgrespro.ru/docs/postgresql/9.6/app-psql). Таким образом вы сможете найти имя созданного индекса.

Чтобы удалить индекс используйте:

```sql
DROP index <index-name>
```