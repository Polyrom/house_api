package flat

import (
	"context"
	"reflect"
	"testing"

	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/pkg/logging"
)

var moderFlatDTOList = []FlatDTO{{ID: 1, HouseID: 1, Price: 1, Rooms: 1, Moderator: "", Status: "created"}}
var clientFlatDTOList = []FlatDTO{{ID: 2, HouseID: 2, Price: 2, Rooms: 2, Moderator: "moder", Status: "approved"}}
var mockCreateFlatDTO = CreateFlatDTO{HouseID: 1, Price: 12_000_000, Rooms: 4}
var mockFlatDTO = FlatDTO{ID: 1, HouseID: 1, Price: 12_000_000, Rooms: 4, Moderator: "", Status: "created"}

func setUpRoleCtx(ctx context.Context, role middleware.Role) context.Context {
	ctx = context.WithValue(ctx, middleware.UserRole, role)
	return ctx
}

type MockLogger struct{}

func (ml *MockLogger) Trace(args ...interface{})                   {}
func (ml *MockLogger) Debug(args ...interface{})                   {}
func (ml *MockLogger) Info(args ...interface{})                    {}
func (ml *MockLogger) Warn(args ...interface{})                    {}
func (ml *MockLogger) Warning(args ...interface{})                 {}
func (ml *MockLogger) Error(args ...interface{})                   {}
func (ml *MockLogger) Fatal(args ...interface{})                   {}
func (ml *MockLogger) Tracef(format string, args ...interface{})   {}
func (ml *MockLogger) Debugf(format string, args ...interface{})   {}
func (ml *MockLogger) Infof(format string, args ...interface{})    {}
func (ml *MockLogger) Warnf(format string, args ...interface{})    {}
func (ml *MockLogger) Warningf(format string, args ...interface{}) {}
func (ml *MockLogger) Errorf(format string, args ...interface{})   {}
func (ml *MockLogger) Fatalf(format string, args ...interface{})   {}
func (ml *MockLogger) Panicf(format string, args ...interface{})   {}

type MockFlatRepo struct{}

func (mfr *MockFlatRepo) GetByHouseIDModerator(ctx context.Context, fl FlatID) ([]FlatDTO, error) {
	return moderFlatDTOList, nil
}
func (mfr *MockFlatRepo) GetByHouseIDClient(ctx context.Context, fl FlatID) ([]FlatDTO, error) {
	return clientFlatDTOList, nil
}
func (mfr *MockFlatRepo) GetByID(ctx context.Context, fl GetFlatByIDDTO) (FlatDTO, error) {
	return FlatDTO{}, nil
}
func (mfr *MockFlatRepo) Create(ctx context.Context, fl CreateFlatDTO) (FlatDTO, error) {
	return mockFlatDTO, nil
}
func (mfr *MockFlatRepo) Update(ctx context.Context, fl UpdateFlatStatusDTO) (FlatDTO, error) {
	return FlatDTO{}, nil
}
func (mfr *MockFlatRepo) UpdateWithNewMod(ctx context.Context, uid string, fl UpdateFlatStatusDTO) (FlatDTO, error) {
	return FlatDTO{}, nil
}

func TestService_GetByHouseID(t *testing.T) {
	type fields struct {
		repo   Repository
		logger logging.Logger
	}
	type args struct {
		ctx context.Context
		f   FlatID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []FlatDTO
		wantErr bool
	}{
		{name: "test client list flats", fields: fields{&MockFlatRepo{}, &MockLogger{}}, args: args{setUpRoleCtx(context.Background(), middleware.Client), 1}, want: clientFlatDTOList, wantErr: false},
		{name: "test moder list flats", fields: fields{&MockFlatRepo{}, &MockLogger{}}, args: args{setUpRoleCtx(context.Background(), middleware.Moderator), 1}, want: moderFlatDTOList, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			got, err := s.GetByHouseID(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetByHouseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetByHouseID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	type fields struct {
		repo   Repository
		logger logging.Logger
	}
	type args struct {
		ctx context.Context
		f   CreateFlatDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    FlatDTO
		wantErr bool
	}{
		{name: "test client list flats", fields: fields{&MockFlatRepo{}, &MockLogger{}}, args: args{setUpRoleCtx(context.Background(), middleware.Client), mockCreateFlatDTO}, want: mockFlatDTO, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
