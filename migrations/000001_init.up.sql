CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(30)      NOT NULL UNIQUE,
    password_hashed VARCHAR(60)      NOT NULL,
    coins           INT DEFAULT 1000 NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users (username);


CREATE TABLE IF NOT EXISTS merch
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(100) NOT NULL UNIQUE,
    price INT          NOT NULL
);


CREATE TABLE IF NOT EXISTS inventory
(
    id       SERIAL PRIMARY KEY,
    user_id  INT REFERENCES users (id),
    merch_id INT REFERENCES merch (id),
    quantity INT NOT NULL,
    date     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, merch_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_user_id ON inventory (user_id);


CREATE TABLE IF NOT EXISTS transactions
(
    id          SERIAL PRIMARY KEY,
    amount      INT NOT NULL,
    receiver_id INT NOT NULL REFERENCES users (id),
    sender_id   INT NOT NULL REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_transactions_receiver_id ON transactions (receiver_id);
CREATE INDEX IF NOT EXISTS idx_transactions_sender_id ON transactions (sender_id);

INSERT INTO merch (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);


