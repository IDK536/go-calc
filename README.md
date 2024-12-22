# go-calc

go-calc - это калькулятор, вычисляющий простые математические выражения.

## API

Проект предоставляет конечную точку API для вычисления выражений. Вы можете отправить POST-запрос на `/api/v1/calculate` с JSON-пакетом, содержащим выражение для вычисления. Программа возвращает результат решения выражение при отсутствии ошибок.

## Использование

### Запуск сервера

Для запуска срвера на стндартном(8080) порте ввудите:
```
go run cmd/main.go
```

### Json файл

```json
{
  "expression": "2+2*2"
}
```
### Пример запроса

```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
### Пример ответа

Возвращает json файл с ответом и статус 200(ok)

```json
{
  "result": 6
}
```

## Ошибки

При некоректном вооде вы получите статус 422 (Unprocessable Entity) и ошибку "invalid expression".

Запрос:
```json
{
  "expression": "+2+2*2"
}
```

Ответ:
```json
{
"err": "invalid expression"
}
```

При вводе пустой строки вы получите статус 422 (Unprocessable Entity) и ошибку "empty expression".

Запрос:
```json
{
  "expression": ""
}
```

Ответ:
```json
{
"err": "empty expression"
}
```

При делении на 0 во время вычисления получите статус 422 (Unprocessable Entity) и ошибку "division by zero".

Запрос:
```json
{
  "expression": "1/0"
}
```

Ответ:
```json
{
"err": "division by zero"
}
```

При возникновении ошибки на сревере. Например если запрос будет пустым. Вы получите статус 500 (Internal Server Error) и ошибку "internal server error".

Запрос:
```json

```

Ответ:
```json
{
"err": "internal server error"
}
```
