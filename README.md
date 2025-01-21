# Test task

## Технологии
- Go
- PostgreSQL
- Docker

## Использованное API:
- [CryptoCompare](https://min-api.cryptocompare.com/documentation)

## Использование
### Добавление криптовалюты:
**POST** `/currency/add`

```json
{
  "coin": "BTC"
}
```

### Для удаления криптовалюты:
**DELETE** `/currency/remove`

```json
{
  "coin": "ETH"
}
```

### Требования:
- [Task](https://taskfile.dev/)
- [Docker](https://www.docker.com/)

### Запуск
```sh
task run
```

### Остановка
```sh
task stop
```