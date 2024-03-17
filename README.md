### REST API сервис "Фильмотека" для управления базой данных фильмов

1. Запрос на авторизацию пользователя /auth
+ В теле запроса передаются e-mail и пароль пользователя (для упрощения работаем с нехешированным паролем)
+ Поскольку не требовалось создания метода регистрации, два базовых пользователя создаются в БД
+ По умолчанию время жизни токена - 1 час, токен возвращается в теле ответа

2. Запрос на создание актера /actor/create
+ Добавление информации об актере, поле name не должно быть пустым
+ В данной реализации принято допущение, что имена актеров уникальные

3. Запрос на редактирование актера /actor/update/{id}
+ Изменение информации об актере, возможное как полное, так и частичное редактирование

4. Запрос на удаление актера actor/delete/{id}
+ Удаление информации об актере

5. Запрос на получение списков актеров actor/get-list
+ Получение полного списка актеров и фильмов, где они снимались

6. Запрос на создание фильма /movie/create
+ Добавление информации о фильме, поле name не должно быть пустым
+ В данной реализации принято допущение, что названия фильмов уникальные
+ Добавлена дополнительная валидация, согласно ТЗ

7. Запрос на редактирование актера /actor/update/{id}
+ Изменение информации о фильме, возможное как полное, так и частичное редактирование

8. Запрос на удаление фильма movie/delete/{id}
+ Удаление информации о фильме

9. Запрос на получение списка фильмов movie/get-list?sortby={}&line={}
+ Получение списка фильмов с возможностью сортировки по названию, по рейтингу, по дате выпуска. По умолчанию используется сортировка по рейтингу(по убыванию)

10. Запрос на получение списка фильмов movie/get-list?actor={}&movie={}
+ Получение списка фильмов по фрагменту названия и/или по фрагменту имени актера

##### Сервис разбит на 3 основных слоя:

- слой обработки запросов: ./api
- слой основной логики: ./internal/service
- слой работы с бд: ./internal/db

##### В качестве основной БД выбран PostgreSQL
Для хранения созданы 4 таблицы:

- actors - для хранения данных актеров
- movies - для хранения данных фильмов
- movie_actors - для хранения связей фильмов и актеров, снимавшихся в них
- users - для хранения пользователей

### Swagger
- по умолчанию документация swagger доступна по адресу http://localhost:8080/swagger
- для возможности работы с документацией, необходимо произвести авторизацию с помощью метода /auth (для упрощения имеются два пользователя в БД)
- после авторизации добавить токен в поле авторизации swagger, без Bearer и ковычек

### Сборка и запуск
Сервис, а так же база данных собирается в docker:  
Для сборки и запуска используется makefile  
По команде make test будут выполнены юнит тесты (тесты слоя работы с БД, выполняются на тестовой базе данных, запускаемой в отдельном контейнере)


