CREATE DATABASE IF NOT EXISTS test;

USE test;

CREATE TABLE IF NOT EXISTS todos (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50),
    done TINYINT,
    created_at DATETIME,
    updated_at DATETIME
);

INSERT INTO todos(name, done, created_at, updated_at) VALUES('todo1', 0, '2020-09-14 13:02:12', '2020-09-14 13:02:12');
INSERT INTO todos(name, done, created_at, updated_at) VALUES('todo2', 1, '2020-09-14 18:36:56', '2020-09-14 20:41:05');

CREATE USER IF NOT EXISTS 'developer' IDENTIFIED BY 'development';

GRANT ALL ON test.* TO 'developer';