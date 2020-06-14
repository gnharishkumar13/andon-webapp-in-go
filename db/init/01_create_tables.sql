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