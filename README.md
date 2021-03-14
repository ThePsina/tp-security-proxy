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
- docker-compose -f deploy/docker-compose.yml up -d -build
```

# Зайти в базу
```
- docker-compose -f deploy/docker-compose.yml exec /bin/bash
- su postgres
- psql -U thepsina -d db
```