create database db_name;
create user db_user with encrypted password 'db_pwd';
grant all privileges on database db_name to db_user;

\c db_name;

CREATE TABLE items(
   id serial PRIMARY KEY,
   name VARCHAR (255) UNIQUE NOT NULL,
   price INTEGER
);

grant all privileges ON table items TO db_user; -- Grant Query & Delete
grant all on sequence items_id_seq to db_user;  -- Grant Insert

INSERT INTO items(name, price) VALUES('apple',  10);
INSERT INTO items(name, price) VALUES('banana', 20);
INSERT INTO items(name, price) VALUES('orange', 30);