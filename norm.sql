CREATE TABLE norm (
       id serial PRIMARY KEY,
       norm text UNIQUE NOT NULL,
       created timestamptz NOT NULL,
       attempts integer NOT NULL
);
