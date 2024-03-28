CREATE TABLE text (
       id serial PRIMARY KEY,
       text text UNIQUE NOT NULL,
       origin cidr NOT NULL,
       norm_id integer REFERENCES norm (id),
       created timestamptz NOT NULL,
       attempts integer NOT NULL

);
GRANT SELECT, INSERT, UPDATE ON text TO indr;
GRANT SELECT, USAGE, UPDATE ON text_id_seq TO indr;
