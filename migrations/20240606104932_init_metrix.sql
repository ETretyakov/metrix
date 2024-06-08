-- +goose Up
-- +goose StatementBegin
CREATE IF NOT EXISTS TABLE public.mtr_metrics (
	id varchar NOT NULL,
	mtype varchar NOT NULL,
	delta bigint NULL,
	value numeric NULL,
	CONSTRAINT mtr_metrics_pk PRIMARY KEY (id)
);
CREATE IF NOT EXISTS INDEX mtr_metrics_id_idx ON public.mtr_metrics (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.mtr_metrics;
-- +goose StatementEnd
