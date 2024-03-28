CREATE TABLE norm (
       id serial PRIMARY KEY,
       norm text UNIQUE NOT NULL,
       attempts integer NOT NULL
);
GRANT SELECT, INSERT, UPDATE ON norm TO indr;
