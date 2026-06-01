# pg-stresstest

Предназначен для отправки 20 потоками 10000 записей в бд

Необходимо для работы:
+ go
+ Переменная в env для подключения к бд

## пример в env:
```env
DB_CONN_STRING=postgres://USER:PASS@IP:PORT/DB_NAME
```