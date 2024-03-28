CREATE TABLE norm (
       id serial PRIMARY KEY,
       norm text UNIQUE NOT NULL,
       created timestamptz NOT NULL,
       attempts integer NOT NULL
);
GRANT SELECT, INSERT, UPDATE ON norm TO indr;
GRANT SELECT, USAGE, UPDATE ON norm_id_seq TO indr;
