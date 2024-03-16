CREATE EXTENSION IF NOT EXISTS moddatetime
    WITH SCHEMA public
    CASCADE;

CREATE TABLE "user"
(
    id         SERIAL PRIMARY KEY,
    email      TEXT  NOT NULL UNIQUE,
    password   BYTEA NOT NULL UNIQUE,
    role       INT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime(updated_at);


CREATE TABLE film
(
    id                 SERIAL PRIMARY KEY,
    title              VARCHAR(150)          NOT NULL,
    description        VARCHAR(1000),
    release_date       DATE
        CONSTRAINT release_date_range
            CHECK (release_date >= '1800-01-01'
                AND release_date <= CURRENT_DATE),
    rating             FLOAT(2)
        CONSTRAINT rating_range
            CHECK (rating BETWEEN 0 AND 10),
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_film_updated_at
    BEFORE UPDATE
    ON film
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE TABLE actor
(
    id                 SERIAL PRIMARY KEY,
    name              TEXT          NOT NULL,
    sex        TEXT,
    birthdate       DATE
        CONSTRAINT birthdate_range
            CHECK (birthdate >= '1800-01-01'
                AND birthdate <= CURRENT_DATE),
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER modify_actor_updated_at
    BEFORE UPDATE
    ON actor
    FOR EACH ROW
    EXECUTE PROCEDURE public.moddatetime(updated_at);

CREATE TABLE film_actor
(
    film_id INTEGER REFERENCES film (id),
    actor_id   INTEGER REFERENCES actor (id),
    PRIMARY KEY (film_id, actor_id)
);
