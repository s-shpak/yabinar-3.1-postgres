# Запуск Postgres в контейнере

Для запуска Postgres в контейнере выполните:

```bash
docker compose up
```

Если вы находитесь в РФ или РБ и не используете VPN, то здесь можно прочитать как использовать Docker без Docker Hub: https://habr.com/ru/articles/818527/ .

Скрипты инициализации лежат в `db/init`.

Для того, чтобы убедиться, что БД была запущена корректно, можно посмотреть ее логи:

```bash
docker logs praktikum-webinar-db
```

Для остановки и полного удаления данных и БД выполните:

```bash
docker compose down --volumes
```

# Создание таблиц и их наполнение

Код для создания таблиц и их наполнения случайными данными лежит в `app/cmd/datagen`.

Для запуска `datagen` выполните следующую комманду:

```bash
make build-datagen

./bin/datagen -dsn postgresql://gopher:gopher@localhost:5432/gopher_corp -emp-count 10000000
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

`http://localhost:8080/employees/ann?limit=50&last-id=1024`

На что стоит обратить внимание:
- PGX pool и его конфигурация
- Пагинация: https://use-the-index-luke.com/no-offset
- Query tracer

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
    id > 542 AND lower(last_name) LIKE 'ann%'
ORDER BY ID asc
;
```

Создадим индекс для т.н. `fuzzy`-поиска:

```sql
CREATE INDEX employees_id_lower_last_name_idx
ON employees using btree(lower(last_name) text_pattern_ops, id);
```

Про класс операторов `text_pattern_ops` можно почитать здесь: https://postgrespro.ru/docs/postgrespro/16/indexes-opclass

Выполним анализ стоимости запроса после добавления индекса:

```sql
EXPLAIN (ANALYZE, VERBOSE)
SELECT id, first_name, last_name, salary, position, email
FROM employees
WHERE
    id > 542 AND lower(last_name) LIKE 'ann%'
ORDER BY ID asc
;
```

Чтобы просмотреть мета-данные таблицы `employees` используйте `\d+ employees` (если вы используете `psql` -- это метакоманда, детали здесь: https://postgrespro.ru/docs/postgresql/9.6/app-psql). Таким образом вы сможете найти имя созданного индекса.

Чтобы удалить индекс используйте:

```sql
DROP index employees_id_lower_last_name_idx
```