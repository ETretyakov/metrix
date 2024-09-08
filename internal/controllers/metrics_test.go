package controllers

import (
	"context"
	"metrix/internal/model"
	"metrix/internal/repository"
	"reflect"
	"testing"
)

func TestMetricControllerImpl_Set(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)

	type fields struct {
		repoGroup *repository.Group
	}
	type args struct {
		metric *model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr bool
	}{
		{
			name:   "Test 1: Create new Counter Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metric: &model.Metric{
					ID:    "Metric 1",
					MType: "counter",
					Delta: func() *int64 { i := int64(100); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "counter",
				Delta: func() *int64 { i := int64(100); return &i }(),
			},
			wantErr: false,
		},
		{
			name:   "Test 2: Increment Counter Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metric: &model.Metric{
					ID:    "Metric 1",
					MType: "counter",
					Delta: func() *int64 { i := int64(100); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "counter",
				Delta: func() *int64 { i := int64(200); return &i }(),
			},
			wantErr: false,
		},
		{
			name:   "Test 3: Create Gauge Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metric: &model.Metric{
					ID:    "Metric 2",
					MType: "gauge",
					Value: func() *float64 { i := float64(100); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 2",
				MType: "gauge",
				Value: func() *float64 { i := float64(100); return &i }(),
			},
			wantErr: false,
		},
		{
			name:   "Test 4: Update Gauge Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metric: &model.Metric{
					ID:    "Metric 2",
					MType: "gauge",
					Value: func() *float64 { i := float64(200); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 2",
				MType: "gauge",
				Value: func() *float64 { i := float64(200); return &i }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricControllerImpl{
				repoGroup: tt.fields.repoGroup,
			}
			got, err := m.Set(ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricControllerImpl.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricControllerImpl.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricControllerImpl_Get(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)

	_, err := repoGroup.MetricRepo.Create(
		ctx,
		&model.Metric{
			ID:    "Metric 1",
			MType: "counter",
			Delta: func() *int64 { i := int64(200); return &i }(),
		},
	)
	if err != nil {
		t.Errorf("MetricControllerImpl.Get() error = %v", err)
		return
	}

	_, err = repoGroup.MetricRepo.Create(
		ctx,
		&model.Metric{
			ID:    "Metric 2",
			MType: "gauge",
			Value: func() *float64 { i := float64(300); return &i }(),
		},
	)
	if err != nil {
		t.Errorf("MetricControllerImpl.Get() error = %v", err)
		return
	}

	type fields struct {
		repoGroup *repository.Group
	}
	type args struct {
		metricID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr bool
	}{
		{
			name:   "Test 1: Get Counter Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metricID: "Metric 1",
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "counter",
				Delta: func() *int64 { i := int64(200); return &i }(),
			},
			wantErr: false,
		},
		{
			name:   "Test 2: Get Gauge Metric record",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metricID: "Metric 2",
			},
			want: &model.Metric{
				ID:    "Metric 2",
				MType: "gauge",
				Value: func() *float64 { i := float64(300); return &i }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricControllerImpl{
				repoGroup: tt.fields.repoGroup,
			}
			got, err := m.Get(ctx, tt.args.metricID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricControllerImpl.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricControllerImpl.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricControllerImpl_SetMany(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)

	type fields struct {
		repoGroup *repository.Group
	}
	type args struct {
		metricsIn []*model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "Test 1: Set many metrics",
			fields: fields{repoGroup: repoGroup},
			args: args{
				metricsIn: []*model.Metric{
					{
						ID:    "Metric 1",
						MType: "gauge",
						Value: func() *float64 { i := float64(300); return &i }(),
					},
					{
						ID:    "Metric 2",
						MType: "counter",
						Delta: func() *int64 { i := int64(200); return &i }(),
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricControllerImpl{
				repoGroup: tt.fields.repoGroup,
			}
			got, err := m.SetMany(ctx, tt.args.metricsIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricControllerImpl.SetMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricControllerImpl.SetMany() = %v, want %v", got, tt.want)
			}
		})
	}
}
