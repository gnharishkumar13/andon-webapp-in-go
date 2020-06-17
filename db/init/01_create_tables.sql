DROP TABLE IF EXISTS public.workcenters;

--- workcenters

CREATE TABLE public.workcenters
(
    id serial,
    wc_name text COLLATE pg_catalog."default",
    current_product text COLLATE pg_catalog."default",
    escalation_level smallint,
    status_set_at timestamp without time zone,
    wc_status smallint,
    PRIMARY KEY (id)
);

ALTER TABLE public.workcenters OWNER to postgres;


--- users

DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users
(
    id serial,
    username text COLLATE pg_catalog."default",
    password text COLLATE pg_catalog."default",
    PRIMARY KEY (id)
);

ALTER TABLE public.users OWNER to postgres;

CREATE INDEX username_pwd_idx_users
    ON public.users USING btree
    (username ASC NULLS LAST, password ASC NULLS LAST)
    TABLESPACE pg_default;

--- logon_tokens

DROP TABLE IF EXISTS public.logon_tokens;

CREATE TABLE public.logon_tokens
(
    id serial,
    token text COLLATE pg_catalog."default",
    user_id int,
    expiration timestamp without time zone,
    PRIMARY KEY (id)
);

ALTER TABLE public.logon_tokens OWNER to postgres;

CREATE INDEX token_idx_logon_tokens
    ON public.logon_tokens USING btree
    (token ASC NULLS LAST)
    TABLESPACE pg_default;

CREATE INDEX user_id_idx_logon_tokens
    ON public.logon_tokens USING btree
    (user_id ASC NULLS LAST)
    TABLESPACE pg_default;
