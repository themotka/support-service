CREATE TABLE platforms
(
    platform_id SERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE
);


CREATE TABLE users
(
    user_id     SERIAL PRIMARY KEY,
    platform_id INTEGER REFERENCES platforms (platform_id)
);


CREATE TABLE managers
(
    manager_id SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE
);


CREATE TABLE conversations
(
    conversation_id SERIAL PRIMARY KEY,
    user_id         INTEGER REFERENCES users (user_id),
    manager_id      INTEGER REFERENCES managers (manager_id),
    reserved_to     TIMESTAMP DEFAULT current_timestamp
);


CREATE TABLE messages
(
    message_id      SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations (conversation_id),
    is_manager      BOOLEAN DEFAULT FALSE,
    content         TEXT    NOT NULL,
    is_sent         BOOLEAN   DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO platforms (name)
VALUES ('telegram');
INSERT INTO platforms (name)
VALUES ('email');


-- Вставка данных в таблицу Manager
INSERT INTO managers (name)
VALUES ('Manager 1');
INSERT INTO managers (name)
VALUES ('Manager 2');
