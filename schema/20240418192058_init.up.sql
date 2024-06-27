CREATE TABLE users
(
  id uuid not null unique,
  name varchar(255) not null,
  username varchar(255) not null unique,
  password_hash varchar(255) not null,
  user_role varchar(255) not null
);

CREATE TABLE sessions
(
  id uuid not null unique,
  user_id uuid references users (id) on delete cascade not null,
  refreshtoken varchar not null unique,
  is_blocked boolean not null,
  created_at timestamptz not null,
  expirated_at timestamptz not null
);

CREATE TABLE items
(
id uuid not null unique,
title varchar(255) not null,
description varchar not null,
price int not null,
image varchar not null
);

-- INSERT INTO users 
-- (
-- id,
-- ) НЕ ЗАБЫТЬ СОЗДАТЬ АДМИНА
-- Изменить тип данных у price на money или decimal