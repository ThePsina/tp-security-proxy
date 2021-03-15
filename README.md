# Proxy server
Тип: Домашняя работа

Номер: 1

Учебная организация: Технопарк МГТУ им. Баумана

Учебный курс: Анализ защищенности

Студент: Трущелев Михаил

Учебная группа: АПО-31

# Необходимо
```
1) docker
2) docker-compose
```

# Запуск
```
- находимся в корневой директории
- docker-compose -f deploy/docker-compose.yml up -d -build
```

# HTTP
```
- curl -x http://localhost:8081/ http://student.bmstu.ru
```

# Повтор запроса
```
- переходим в браузере по адресу localhost:8082/requests
- находим id ревеста
- переходим по адресу localhost:8082/repeat/<id>
```

# HTTPS
```
- '-k' - скипаем ошибки сертификации курла
- curl -x http://localhost:8081/ https://mail.ru -k
```

# Scanner
```
- curl -x http://localhost:8081/ http://student.bmstu.ru/?iphone=true\&sandbox=false\&notOurMinerParam\=4
- переходим в браузере по адресу localhost:8082/requests
- находим id ревеста
- переходим по адресу localhost:8082/scan/<id>
```

# Зайти в базу
```
- docker-compose -f deploy/docker-compose.yml exec /bin/bash
- su postgres
- psql -U thepsina -d db
```