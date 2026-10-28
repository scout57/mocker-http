# scout55/mocker-http

Лёгкий HTTP-мок-сервер на Go. Читает YAML-шаблоны и отвечает заглушками на входящие запросы — удобно для локальной разработки и тестов.

## Возможности

- Ответы из YAML-шаблонов в `./templates`
- Точное совпадение путей и wildcard `*` (один сегмент)
- Приоритет точного совпадения над wildcard
- `response.status` (по умолчанию `200`)
- Тело ответа: строка, объект или массив (объекты/массивы сериализуются в JSON)
- `Content-Type: application/json` по умолчанию
- Без кэша — шаблоны перечитываются на каждый запрос
- Fallback `hello world`, если совпадений нет
- Только стандартная библиотека + `gopkg.in/yaml.v3`



## Формат шаблона

Все файлы `./templates/*.yaml` и `./templates/*.yml` объединяются. Каждый файл — массив объектов:

```yaml
- method: GET              # case-insensitive
  path: /users/*           # * — один сегмент
  response:
    status: 201            # опционально, default 200
    headers:               # опционально
      X-Foo: bar
    body: ...              # строка | объект | массив
```

## Использование в своих проектах
```yaml
  scout_mocker_http:
    image: scout55/mocker:2.0.0
    volumes:
      - ./templates:/go/app/templates:ro
    ports:
      - 8080:8080
```


## Локальный разработка

```bash
# Запуск окружения
docker compose up -d

# Зависимости
docker compose exec app sh -c ". install.sh"

# Hot-reload
docker compose exec app air

# Выключить
docker compose down
```