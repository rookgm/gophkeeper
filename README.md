# Менеджер паролей GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

### Типы хранимой информации
- пары логин/пароль;
- произвольные текстовые данные;
- произвольные бинарные данные;
- данные банковских карт.


## Особенности
1. JWT аутентификация.
2. Шифрование данных AES-256-GSM и формирование ключа PBKDF2 SHA-256.
3. Поддержка HTTPS соединения.
4. Данные шифруются на стороне клиента.

## Установка

### Требования
- Go 1.22 и выше
- PostgreSQL 16 и выше
- Утилита make

## Сборка
1. Клонирование репозитория
```bash
git clone https://github.com/rookgm/gophkeeper
cd gophkeeper
```
2. Установка зависимостей
```bash
go mod tidy
```

3. Создание базы данных
```bash
create database gophkeeper;

create user gophkeeper with encrypted password 'you_password';

grant all privileges on database gophkeeper to gophkeeper;

GRANT CREATE ON SCHEMA public TO gophkeeper;

alter database gophkeeper owner to gophkeeper;
```
4. Сборка сервера и клиента
```bash
make
```
Сборка для ОС Linux, Windows, MacOS
```bash
make all
```


## Запуск

### Сервер
```bash
Общий вид команды сервера

gophkeeper-server [-a server_addr] [-d dsn] [-l log_level] [-c config_name]

Описание параметров
          -a 
              адрес сервера
          -d
              имя источника данных
          -l
              уровень логирования
          -c 
              имя файла конфигурации
```
Переменные окружения:

    SERVER_ADDRESS - адрес сервера
    SERVER_DATABASE_DSN - имя источника данных
    SERVER_LOG_LEVEL - уровень логирования
    SERVER_AUTH_TOKEN_KEY - ключ для подписи JWT

### Клиент
```bash
Общий вид команды клиента

gophkeeper-client [-a server_addr] [-l log_level] [-f config_dir]

Описание параметров
          -a 
              адрес сервера
          -l
             уровень логирования
          -f
              задает путь до директории конфигурации
```
Переменные окружения:

    GOPHKEEPER_SERVER_ADDRESS - адрес сервера
    CLIENT_LOG_LEVEL - уровень логирования
    CLIENT_CONFIG_DIR - путь до директории конфигурации


### Описание команд

| Команда                                               | Описание                                         |
|-------------------------------------------------------|--------------------------------------------------|
| `version`                                             | Получение информации о версии приложения клиента |
| `register -u <username>`                              | Регистрация пользователя                         |
| `login -u <username>`                                 | Аутентификация пользователя                      |
| `secret add credentials -n <name> -l <login> -e note` | Добавление учетных данных                        |
| `secret add text -n <name> -c <content> -e note`      | Добавление текстовых данных                      |
| `secret add binary -n <name> -p <file_name> -e note`  | Добавление бинарных данных                       |
| `secret add card -n <name> -e <note>`                 | Добавление банковской карты                      |
| `secret get <secretid>`                               | Просмотр секрета                                 |
| `secret delete <secretid>`                            | Удаление секрета                                 |

### Примеры использования

#### Аутентификация пользователя
```bash
$ ./gophkeeper login -u gopher
password for gopher: *****
login successfully
```

#### Добавление учетных данных
```bash
$ ./gophkeeper secret add credentials
Enter master password: *****
Name: github
Login: username
Note: dev
password for username: *****
Successfully added credentials, ID: 25e5c83f-5238-470b-952a-9202b226b824
```

#### Добавление банковской карты
```bash
$ ./gophkeeper secret add card
Enter master password: ******
Name: multicard
Note: visa
Card number: 12345789101112
Card expiration month: 05
Card expiration year: 26
Card holder name: john
Cardholders billing address: city
Card type: debit
Issue name: bank
CCV: ***
Successfully added bank card, ID: 36264f0f-f09f-412d-9fd1-cd7a30454055
```

#### Просмотр секрета
```bash
$ ./gophkeeper secret get 25e5c83f-5238-470b-952a-9202b226b824
Enter master password: *****
===CREDENTIALS DETAILS===
ID: 25e5c83f-5238-470b-952a-9202b226b824
Name: github
Note: dev
Login: username
Password: 123
Created: 2025-10-26 20:48:29
Updated: 2025-10-26 20:48:29
```

## Сводное HTTP API

Менеджер паролей GophKeeper предоставляет следующие HTTP-хендлеры:

- POST /api/user/register — регистрация пользователя;
- POST /api/user/login — аутентификация пользователя;
- POST /api/user/secrets — создание секрета;
- GET /api/user/secrets/{id} — получение секрета;
- PUT /api/user/secrets/{id} — обновление секрета;
- DELETE /api/user/secrets/{id} — удаление секрета;

### Регистрация пользователя

Хендлер: **POST /api/user/register**.

Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
Для передачи аутентификационных данных используется механизм _JWT_.

Формат запроса:

    POST /api/user/register HTTP/1.1
    Content-Type: application/json
    ...
    
    {
    "login": "<login>",
    "password": "<password>"
    }

Коды ответа:

- 200 — пользователь успешно зарегистрирован;
- 400 — неверный формат запроса;
- 409 — логин уже занят;
- 500 — внутренняя ошибка сервера.

### Аутентификация пользователя

Хендлер: **POST /api/user/login**.

Аутентификация производится по паре логин/пароль.
Для передачи аутентификационных данных используется механизм _JWT_.

Формат запроса:

    POST /api/user/login HTTP/1.1
    Content-Type: application/json
    ...
    
    {
    "login": "<login>",
    "password": "<password>"
    }

Коды ответа:

- 200 — пользователь успешно аутентифицирован;
- 400 — неверный формат запроса;
- 401 — неверная пара логин/пароль;
- 500 — внутренняя ошибка сервера.

### Создание секрета
- POST /api/user/secrets

Формат запроса:

    POST /api/user/secrets HTTP/1.1
    Content-Type: application/json
    ...    

    {
        "name":"github",
        "type":0,
        "note":"dev",
        "data":"encrypted_data"
    }

Формат ответа:
    
    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "id":"f5f10806-5784-4e39-a592-d45abd41bdf3",
        "name":"github",
        "type":0,
        "note":"dev",
        "data":""encrypted data",
        createdAt":"2025-10-26T23:36:35.863665+03:00",
        "updatedAt":"2025-10-26T23:36:35.863665+03:00"
    } 

Коды ответа:
- 201 - секрет успешно создан;
- 400 - неверный формат запроса;
- 401 - пользователь не авторизован;
- 405 - метод запрещен;
- 500 - внутренняя ошибка сервера.


### Получение секрета
- GET /api/user/secrets/{id}

Формат запроса:

    GET /api/user/secrets/f5f10806-5784-4e39-a592-d45abd41bdf3 HTTP/1.1
    Content-Type: application/json
    ...

Формат ответа:

    200 OK HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "id":"f5f10806-5784-4e39-a592-d45abd41bdf3",
        "name":"github",
        "type":0,
        "note":"dev",
        "data":""encrypted data",
        createdAt":"2025-10-26T23:36:35.863665+03:00",
        "updatedAt":"2025-10-26T23:36:35.863665+03:00"
    }

Коды ответа:
- 201 - секрет успешно получен;
- 400 - неверный формат запроса;
- 401 - пользователь не авторизован;
- 404 - секрет не найден;
- 405 - метод запрещен;
- 500 - внутренняя ошибка сервера.

### Обновление серета
- PUT /api/user/secrets/{id}

Формат запроса:

    PUT /api/user/secrets/f5f10806-5784-4e39-a592-d45abd41bdf3 HTTP/1.1
    Content-Type: application/json
     ...

    {
        "name":"github",
        "type":0,
        "note":"prod,updated",
        "data":"encrypted_data"
    }

Коды ответа:
- 200 - секрет успешно обновлен;
- 400 - неверный формат запроса;
- 401 - пользователь не авторизован;
- 404 - секрет не найден;
- 405 - метод запрещен;
- 500 - внутренняя ошибка сервера.

### Удаление серета
- DELETE /api/user/secrets/{id}
Формат запроса:

      PUT /api/user/secrets/f5f10806-5784-4e39-a592-d45abd41bdf3 HTTP/1.1
      Content-Type: application/json
      ...

Коды ответа:
- 200 - секрет успешно удален;
- 400 - неверный формат запроса;
- 401 - пользователь не авторизован;
- 404 - секрет не найден;
- 405 - метод запрещен;
- 500 - внутренняя ошибка сервера.