// Module "repositry" that holds the functionality that is related to database communication.
package repository

import (
	"context"
	"metrix/internal/model"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestMetricRepositoryImpl_Create(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	rows := mock.
		NewRows([]string{"id", "mtype", "delta", "value"}).
		AddRow("Metric 1", "gauge", nil, func() *float64 { i := float64(300); return &i }())

	expectedQuery := `INSERT INTO +.?`
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	type fields struct {
		gr *Group
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
			name:   "Test 1: Success create",
			fields: fields{gr: gr},
			args: args{
				&model.Metric{
					ID:    "Metric 1",
					MType: "gauge",
					Value: func() *float64 { i := float64(300); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "gauge",
				Value: func() *float64 { i := float64(300); return &i }(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.Create(ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricRepositoryImpl.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_Read(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	rows := mock.
		NewRows([]string{"id", "mtype", "delta", "value"}).
		AddRow("Metric 1", "gauge", nil, func() *float64 { i := float64(300); return &i }())

	expectedQuery := `SELECT +.?`
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	type fields struct {
		gr *Group
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
			name:   "Test 1: Success read",
			fields: fields{gr: gr},
			args: args{
				metricID: "Metric 1",
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "gauge",
				Value: func() *float64 { i := float64(300); return &i }(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.Read(ctx, tt.args.metricID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricRepositoryImpl.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_ReadIDs(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	rows := mock.
		NewRows([]string{"id"}).
		AddRow("Metric 1")

	expectedQuery := `SELECT +.?`
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	type fields struct {
		gr *Group
	}
	type args struct{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]string
		wantErr bool
	}{
		{
			name:    "Test 1: Success read many",
			fields:  fields{gr: gr},
			args:    args{},
			want:    &[]string{"Metric 1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.ReadIDs(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.ReadIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricRepositoryImpl.ReadIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_ReadMany(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	rows := mock.
		NewRows([]string{"id", "mtype", "delta", "value"}).
		AddRow("Metric 1", "gauge", nil, func() *float64 { i := float64(300); return &i }())

	expectedQuery := `SELECT +.?`
	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	type fields struct {
		gr *Group
	}
	type args struct {
		metricIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]model.Metric
		wantErr bool
	}{
		{
			name:   "Test 1: Success read many",
			fields: fields{gr: gr},
			args: args{
				metricIDs: []string{"Metric 1"},
			},
			want: &[]model.Metric{
				{
					ID:    "Metric 1",
					MType: "gauge",
					Value: func() *float64 { i := float64(300); return &i }(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.ReadMany(ctx, tt.args.metricIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.ReadMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricRepositoryImpl.ReadMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_Update(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE +.?`).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectQuery(`SELECT +.?`).WillReturnRows(
		mock.
			NewRows([]string{"id", "mtype", "delta", "value"}).
			AddRow("Metric 1", "gauge", nil, func() *float64 { i := float64(300); return &i }()),
	)
	mock.ExpectCommit()

	type fields struct {
		gr *Group
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
			name:   "Test 1: Success update",
			fields: fields{gr: gr},
			args: args{
				metric: &model.Metric{
					ID:    "Metric 1",
					MType: "gauge",
					Value: func() *float64 { i := float64(300); return &i }(),
				},
			},
			want: &model.Metric{
				ID:    "Metric 1",
				MType: "gauge",
				Value: func() *float64 { i := float64(300); return &i }(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.Update(ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricRepositoryImpl.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_UpsertMany(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT +.?`).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	mock.ExpectCommit()

	type fields struct {
		gr *Group
	}
	type args struct {
		metrics []model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "Test 1: Success upsert",
			fields: fields{gr: gr},
			args: args{
				metrics: []model.Metric{
					{
						ID:    "Metric 1",
						MType: "gauge",
						Value: func() *float64 { i := float64(300); return &i }(),
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			got, err := r.UpsertMany(ctx, tt.args.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.UpsertMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricRepositoryImpl.UpsertMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricRepositoryImpl_Delete(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.FailNow()
	}

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	gr := NewGroup(ctx, sqlxDB, "", 0, false)

	mock.ExpectExec(`DELETE +.?`).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	type fields struct {
		gr *Group
	}
	type args struct {
		metricID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1: Success upsert",
			fields:  fields{gr: gr},
			args:    args{metricID: "Metric 1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepositoryImpl{
				gr: tt.fields.gr,
			}
			if err := r.Delete(ctx, tt.args.metricID); (err != nil) != tt.wantErr {
				t.Errorf("MetricRepositoryImpl.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
