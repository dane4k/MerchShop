## Для запуска проекта нужно:
### 1) Заполнить .env:
```
SERVER_PORT=8080              # Порт проекта
DB_HOST=db                    # Адрес сервера БД
DB_PORT=5432                  # Порт PostgreSQL
DB_USER=postgres              # Пользователь PostgreSQL
DB_PASSWORD=admin             # Пароль пользователя
DB_NAME=merch_shop            # Название базы данных
JWT_SECRET=secret             # JWT secret key
LOGGER_FILE=MerchShop.log     # Файл для логов
```
### 2) Выполнить команду:
```
docker-compose up --build
```

Для разделения ответственности надо использовать отдельные сущности для БД и для хендлеров, но мне показалось, 
что для такого маленького приложения большое количество преобразований при таком сценарии будет избыточным

### Нагрузочное тестирование 1000 RPS:
```
wrk2 -t10 -c100 -d60s -R1000 --latency http://localhost:8080/api/buy/umbrella
```
- /api/buy/{item}:
```
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.56ms
 75.000%    1.93ms
 90.000%    2.48ms
 99.000%    9.96ms
 99.900%   22.03ms
 99.990%   31.89ms
 99.999%   41.95ms
100.000%   41.95ms
```
- /api/info:
```
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.63ms
 75.000%    2.06ms
 90.000%    2.69ms
 99.000%    7.05ms
 99.900%   17.52ms
 99.990%   25.39ms
 99.999%   28.93ms
100.000%   28.93ms
```

- /api/sendCoin:
```
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.51ms
 75.000%    1.86ms
 90.000%    2.26ms
 99.000%    5.51ms
 99.900%   19.07ms
 99.990%   37.89ms
 99.999%   50.37ms
100.000%   50.37ms
```

- /api/auth:
```
  Latency Distribution (HdrHistogram - Recorded Latency)
 50.000%    1.47ms
 75.000%    1.78ms
 90.000%    2.14ms
 99.000%    4.11ms
 99.900%   17.66ms
 99.990%   25.81ms
 99.999%   34.43ms
100.000%   34.43ms
```
