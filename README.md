
# Примеры запросов <a name="examples"></a>


* [Список городов: GET http://localhost:8000/cities]


* [Список с кратким предсказанием для выбранного города: GET http://localhost:8000/city/{city}/forecast]
```
Реальный запрос: http://localhost:8000/city/London/forecast
```


* [Детальная информация о погоде для конкретного города и конкретного
  времени: GET http://localhost:8000/city/{city}/weather/{datetime}]
```
Реальный запрос: http://localhost:8000/city/London/weather/2024-07-15T12:00:00
```

* [Регистрация пользователя: POST http://localhost:8000/register с телом]
```
{
    "username": "exampleUser3",
    "password": "examplePassword3"
}
```
  
* [Вход пользователя: POST http://localhost:8000/login с телом]
```
{
    "username": "exampleUser",
    "password": "examplePassword"
}
```


* [Добавление города в список избранных
  пользователя: POST http://localhost:8000/favorite с телом]
```
{
    "user_id": 4,
    "city_name": "London"
}
```



```
Процесс запуска: 

Для запуска БД: docker-compose up --build
Прогон миграций делается из: make migration-up (Нужно перед этим установить goose)
Для запуска приложения нужно перейти в директорию cmd/main и запустить команду go run main.go
Для запуска тестов нужно перейти в директорию internal/order/storage и запустить go test

```

# Документация находится в директории api