CREATE TABLE text (
       id serial PRIMARY KEY,
       text text UNIQUE NOT NULL,
       origin cidr NOT NULL,
       norm_id integer REFERENCES norm (id),
       created timestamptz NOT NULL,
       attempts integer NOT NULL,
       seen integer NOT NULL,
       edited integer NOT NULL
);
