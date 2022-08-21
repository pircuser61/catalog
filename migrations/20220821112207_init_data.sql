-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
INSERT INTO public.country (name) VALUES ('Россия') ON CONFLICT (name) DO NOTHING;
INSERT INTO public.country (name) VALUES ('Китай') ON CONFLICT (name) DO NOTHING;
INSERT INTO public.unit_of_measure (name) VALUES ('шт') ON CONFLICT (name) DO NOTHING;
INSERT INTO public.unit_of_measure (name) VALUES ('кг') ON CONFLICT (name) DO NOTHING;
INSERT INTO public.unit_of_measure (name) VALUES ('л') ON CONFLICT (name) DO NOTHING;
INSERT INTO public.unit_of_measure (name) VALUES ('м.кв') ON CONFLICT (name) DO NOTHING;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
