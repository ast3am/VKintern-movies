CREATE TABLE users(
    uuid          uuid primary key,
    email         varchar unique not null,
    password varchar        not null,
    role    varchar not null
);

CREATE TABLE actors(
    uuid       uuid primary key,
    name       varchar unique not null,
    gender     varchar(6),
    birth_date timestamp
);

CREATE TABLE movies(
    uuid         uuid primary key not null,
    name         varchar unique not null,
    description  text,
    release_date timestamp,
    rating       float
);

CREATE TABLE movie_actors(
    movie_uuid uuid references movies (uuid) on delete cascade,
    actor_uuid uuid references actors (uuid) on delete cascade,
    actor_name varchar,
    unique (movie_uuid, actor_uuid)
);

insert into users
values ('2300a1f6-b2aa-4f5b-b6ca-8f495582e255', 'testuser@mail.com',
        'userPassword', 'user');

insert into users
values ('482d6f53-b2ee-4684-887e-2588ae6c9d48',
        'admin@vk.ru', 'adminPassword#1', 'admin') ;

