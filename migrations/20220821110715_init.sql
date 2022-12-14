-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.country (
	country_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY,
	"name" varchar NOT NULL,
	CONSTRAINT country_pk PRIMARY KEY (country_id),
    CONSTRAINT country_un UNIQUE ("name")
);

CREATE TABLE IF NOT EXISTS public.unit_of_measure (
	unit_of_measure_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY,
	"name" varchar NOT NULL,
	CONSTRAINT unit_of_measure_pk PRIMARY KEY (unit_of_measure_id),
	CONSTRAINT unit_of_measure_un UNIQUE ("name")
);
 
CREATE TABLE IF NOT EXISTS public.good (
	"name" varchar(255) NOT NULL,
	code int8 NOT NULL GENERATED BY DEFAULT AS IDENTITY,
	unit_of_measure_id int4 NOT NULL,
	country_id int4 NOT NULL,
	CONSTRAINT good_pk PRIMARY KEY (code),
	CONSTRAINT country_fk FOREIGN KEY (country_id) REFERENCES public.country(country_id),
	CONSTRAINT unit_of_measure_fk FOREIGN KEY (unit_of_measure_id) REFERENCES public.unit_of_measure(unit_of_measure_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.good;
DROP TABLE IF EXISTS public.country;
DROP TABLE IF EXISTS public.unit_of_measure;
-- +goose StatementEnd
