-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.mtr_metrics (
	id varchar NOT NULL,
	mtype varchar NOT NULL,
	delta bigint NULL,
	value numeric NULL,
	CONSTRAINT mtr_metrics_pk PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS  mtr_metrics_id_idx ON public.mtr_metrics (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.mtr_metrics;
-- +goose StatementEnd
