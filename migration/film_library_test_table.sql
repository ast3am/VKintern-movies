CREATE TABLE users(
    uuid     uuid primary key,
    email    varchar unique not null,
    password varchar        not null,
    role     varchar        not null
);

CREATE TABLE actors(
    uuid       uuid primary key,
    name       varchar unique not null,
    gender     varchar(6),
    birth_date timestamp
);

CREATE TABLE movies(
    uuid         uuid primary key not null,
    name         varchar unique   not null,
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

INSERT INTO  users
VALUES ('2300a1f6-b2aa-4f5b-b6ca-8f495582e255', 'testuser@mail.com',
        'userPassword', 'user');

INSERT INTO  users
VALUES ('482d6f53-b2ee-4684-887e-2588ae6c9d48',
        'admin@vk.ru', 'adminPassword#1', 'admin');

INSERT INTO actors (uuid, name, gender, birth_date)
VALUES
    ('b0482c7a-1a4c-4a3c-9463-35f0036a0d60', 'Tom Hanks', 'Male', '1956-07-09'),
    ('7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416', 'Julia Roberts', 'Female', '1967-10-28'),
    ('eb1b5f32-82c3-4dfb-aad2-9432908d12b7', 'Leonardo DiCaprio', 'Male', '1974-11-11'),
    ('d86d5541-bda1-4b89-9256-5cfbba11dc89', 'Cate Blanchett', 'Female', '1969-05-14'),
    ('40d882f7-b027-4a07-85da-76e0f7d9b6e3', 'Brad Pitt', 'Male', '1963-12-18');


INSERT INTO movies (uuid, name, description, release_date, rating)
VALUES
    ('20648636-b14e-4f88-a02a-8c4e3c2f534d', 'Forrest Gump', 'Description 1', '1994-07-06', 8.8),
    ('f44d4a8f-7f16-4c1d-836b-02e0b8de4a99', 'The Devil Wears Prada', 'Description 2', '2006-06-30', 6.9),
    ('a5e85bc9-5e75-4e5f-b890-0b156e15e40f', 'Training Day', 'Description 3', '2001-10-05', 7.7),
    ('43a4b2df-8ff8-4d32-9203-0367c2714431', 'Pretty Woman', 'Description 4', '1990-03-23', 7.0),
    ('b19069b7-7296-4e67-a5f2-80c6e01a32d0', 'Inception', 'Description 5', '2010-07-16', 8.8);

INSERT INTO movie_actors (movie_uuid, actor_uuid, actor_name)
VALUES
    ('20648636-b14e-4f88-a02a-8c4e3c2f534d', 'b0482c7a-1a4c-4a3c-9463-35f0036a0d60', 'Tom Hanks'),
    ('f44d4a8f-7f16-4c1d-836b-02e0b8de4a99', '7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416', 'Julia Roberts'),
    ('a5e85bc9-5e75-4e5f-b890-0b156e15e40f', 'eb1b5f32-82c3-4dfb-aad2-9432908d12b7', 'Leonardo DiCaprio'),
    ('43a4b2df-8ff8-4d32-9203-0367c2714431', 'd86d5541-bda1-4b89-9256-5cfbba11dc89', 'Cate Blanchett'),
    ('b19069b7-7296-4e67-a5f2-80c6e01a32d0', '40d882f7-b027-4a07-85da-76e0f7d9b6e3', 'Brad Pitt'),
    ('20648636-b14e-4f88-a02a-8c4e3c2f534d', '7cfd7e7f-3a27-46c5-9c3e-3b3d3b1c1416', 'Julia Roberts'),
    ('a5e85bc9-5e75-4e5f-b890-0b156e15e40f', 'b0482c7a-1a4c-4a3c-9463-35f0036a0d60', 'Tom Hanks'),
    ('43a4b2df-8ff8-4d32-9203-0367c2714431', 'b0482c7a-1a4c-4a3c-9463-35f0036a0d60', 'Tom Hanks'),
    ('f44d4a8f-7f16-4c1d-836b-02e0b8de4a99', '40d882f7-b027-4a07-85da-76e0f7d9b6e3', 'Brad Pitt'),
    ('b19069b7-7296-4e67-a5f2-80c6e01a32d0', 'd86d5541-bda1-4b89-9256-5cfbba11dc89', 'Cate Blanchett');

