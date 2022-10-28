# Тестовое задание на позицию стажёра-бэкендера

**Инструкция по запуску**
1. 
```shell
    git clone https://github.com/krm-shrftdnv/avito-internship.git
    cd avito-internship
    cp ./.template.env .env
```
2. Меняем значения переменных среды при желании
3. 
```shell
    docker-compose up --build
```

**Примеры запросов ответов**

[Postman](https://www.postman.com/krm-shrftdnv/workspace/avito-internship-backend-krm-shrftdnv/collection/24117006-5ce35bb5-3aa2-431a-9487-cdca44836372?action=share&creator=24117006&ctx=documentation)

[Swagger](https://app.swaggerhub.com/apis/krm-shrftdnv/avito-internship/1.0)


**Комментарии по ТЗ**

1. Не был описан механизм внесения данных об услугах, поэтому при резервировании создаётся услуга при получении нового service_id.
2. По доп. заданию 2: пропустил часть с комментариями к транзакциям, т.к. непонятно откуда брать, например, информацию о пополнении баланса.