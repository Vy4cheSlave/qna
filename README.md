# qna
questions and answers backend api

# Требования
- Язык программирования: Golang
- net/http
- PostgreSQL
- GORM
- goose
- go.uber.org/zap
- log/slog
- httptest/testify
- github.com/joho/godotenv
- Docker, docker-compose

# Шаги запуска
1. Клонировать репозиторий:
```
git clone <repo>
```
2. Перейти в корневую директорию проекта:
```
cd <qna>
```
3. Создать `.env` файл в корневой директории, по примеру `.env.example`
```bash
cp .env.example .env
```
4. Запустить сервер:
```
docker-compose -f .\docker\docker-compose.yml --project-directory . up --build
```

# API
Вопросы (Questions):
- GET /questions/ — список всех вопросов
- POST /questions/ — создать новый вопрос
- GET /questions/{id} — получить вопрос и все ответы на него
- DELETE /questions/{id} — удалить вопрос (вместе с ответами)

Ответы (Answers):
- POST /questions/{id}/answers/ — добавить ответ к вопросу
- GET /answers/{id} — получить конкретный ответ
- DELETE /answers/{id} — удалить ответ

Пользователи (Users):
- POST /users/ - создать пользователя
- GET /users/ - получить всех пользователей
- DELETE /users/{id} - удалить пользователя

# Логика:
- Нельзя создать ответ к несуществующему вопросу/ несуществующим пользователем.
- Один и тот же пользователь может оставлять несколько ответов на один вопрос.
- При удалении вопроса должны удаляться все его ответы (каскадно).
- При удалении пользователя должны удаляться все его ответы (каскадно).

# Описание директорий
```
.
├───cmd                 # Точки входа: файлы main.go
├───docker              # Файлы инфраструктуры: Dockerfile и docker-compose.yml
├───internal            # Внутренний код приложения.
│   ├───config          # Загрузка и парсинг конфигурации из .env, флагов командной строки.
│   ├───domain          # (DDD): Сущности (Entities)
│   ├───infrastructure  # Слой инфраструктуры.
│   │   ├───db          # Реализация репозиториев для БД.
│   │   │   └───dto     # Структуры данных БД с GORM-тегами.
│   │   └───rest        # Реализация HTTP API.
│   │       ├───dto     # Объекты передачи данных для REST.
│   │       │   ├───request  # Структуры входящих JSON-запросов.
│   │       │   └───response # Структуры исходящих JSON-ответов.
│   │       ├───middleware # HTTP-промежуточное ПО: CORS, логирование.
│   │       └───mocks   # моки для Unit-тестирования ручек
│   ├───logpack         # Реализация логирования log/slog над zap.
│   └───usecase         # Слой сервисов приложения
└───migrations          # Скрипты миграции БД.
    └───postgres        # SQL-файлы для goose.
```

Все было реализовано, опираясь на принципы SOLID, DDD, clean architecture