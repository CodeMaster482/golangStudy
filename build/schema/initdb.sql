CREATE OR REPLACE TABLE user (
    id         SERIAL PRIMARY KEY,
    login      TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,
    balance    NUMERIC DEFAULT 0,
    is_moderator BOOLEAN,
    CHECK(LENGTH(login) <= 255),
    CHECK(LENGTH(password) <= 255)
);

-- Таблица операций
CREATE OR REPLACE TABLE operation (
    id         SERIAL PRIMARY KEY,
    user_id    REFERENCES user(id) NOT NULL,
    is_income  BOOLEAN,
    status     TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    formation_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    completion_at TIMESTAMPTZ,
    CHECK(LENGTH(status) <= 5)
);

-- Таблица банкнот
CREATE OR REPLACE TABLE banknote (
    id          SERIAL PRIMARY KEY,
    nominal     NUMERIC,
    image_url   TEXT CHECK(LENGTH(status) <= 500) NOT NULL,
    status      TEXT CHECK(LENGTH(status) <= 30)
);

-- Таблица связи между операцией и банкнотой
CREATE OR REPLACE TABLE operation_banknote (
    id              SERIAL PRIMARY KEY,
    operation_id    REFERENCES operation(id) NOT NULL,
    banknote_id     REFERENCES banknote(id) NOT NULL,
    quantity        INT DEFAULT 1 CHECK(quantity >= 1),
);
