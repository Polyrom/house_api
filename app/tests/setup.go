package tests

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/Polyrom/houses_api/internal/config"
	"github.com/Polyrom/houses_api/internal/flat"
	"github.com/Polyrom/houses_api/internal/house"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/internal/modstatus"
	"github.com/Polyrom/houses_api/internal/server"
	"github.com/Polyrom/houses_api/internal/user"
	"github.com/Polyrom/houses_api/pkg/client/postgres"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

type testContext struct {
	Server         *server.Server
	ModeratorToken middleware.Token
	ClientToken    middleware.Token
	Houses         map[int]house.House
}

func (ctx *testContext) setup() {
	userRepo := user.NewRepository(ctx.Server.DB, &MockLogger{})
	testModerID, err := createTestModerator(userRepo)
	if err != nil {
		log.Fatal(err)
	}
	testClientID, err := createTestClient(userRepo)
	if err != nil {
		log.Fatal(err)
	}
	moderToken, err := createUserToken(userRepo, testModerID)
	if err != nil {
		log.Fatal(err)
	}
	clientToken, err := createUserToken(userRepo, testClientID)
	if err != nil {
		log.Fatal(err)
	}
	ctx.ModeratorToken = moderToken
	ctx.ClientToken = clientToken
	houseRepo := house.NewRepository(ctx.Server.DB, &MockLogger{})
	h, err := createHouse(houseRepo)
	if err != nil {
		log.Fatal(err)
	}
	ctx.Houses[h.ID] = h
}

func (ctx *testContext) cleanup() {
	err := deleteFlats(ctx.Server.DB)
	if err != nil {
		log.Fatal(err)
	}
	err = deleteHouses(ctx.Server.DB)
	if err != nil {
		log.Fatal(err)
	}
	err = deleteTokens(ctx.Server.DB)
	if err != nil {
		log.Fatal(err)
	}
	err = deleteUsers(ctx.Server.DB)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteFlats(db *pgxpool.Pool) error {
	dq := `DELETE FROM flats;`
	_, err := db.Exec(context.Background(), dq)
	if err != nil {
		return err
	}
	sq := `ALTER SEQUENCE flats_id_seq RESTART WITH 1;`
	_, err = db.Exec(context.Background(), sq)
	if err != nil {
		return err
	}
	return nil
}

func deleteHouses(db *pgxpool.Pool) error {
	dq := `DELETE FROM houses;`
	_, err := db.Exec(context.Background(), dq)
	if err != nil {
		return err
	}
	sq := `ALTER SEQUENCE houses_id_seq RESTART WITH 1;`
	_, err = db.Exec(context.Background(), sq)
	if err != nil {
		return err
	}
	return nil
}

func deleteUsers(db *pgxpool.Pool) error {
	dq := `DELETE FROM users;`
	_, err := db.Exec(context.Background(), dq)
	if err != nil {
		return err
	}
	return nil
}

func deleteTokens(db *pgxpool.Pool) error {
	dq := `DELETE FROM tokens;`
	_, err := db.Exec(context.Background(), dq)
	if err != nil {
		return err
	}
	return nil
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

var testStorageCfg = config.StorageConfig{
	Username:    "testuser",
	Password:    "testpassword",
	Host:        "localhost",
	Port:        "5433",
	Database:    "testdb",
	MaxAttempts: 5,
}
var testCfg = config.Config{
	Debug: new(bool),
	Listen: struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}{
		Host: "localhost",
		Port: "8080",
	},
	Storage: testStorageCfg,
}

func newTestServer() *server.Server {
	pg, err := postgres.NewClient(context.Background(), testStorageCfg)
	if err != nil {
		os.Exit(1)
	}
	router := mux.NewRouter()
	server := server.New(&testCfg, &MockLogger{}, router, pg)
	server.ConfigureRouter()
	return server
}

func createTestModerator(r user.Repository) (user.UserID, error) {
	uid := uuid.New().String()
	u := user.User{
		ID:       user.UserID(uid),
		Email:    "moderator@haha.foo",
		Password: "testpassword",
		Role:     string(middleware.Moderator),
	}
	return r.Create(context.Background(), u)
}

func createTestClient(r user.Repository) (user.UserID, error) {
	uid := uuid.New().String()
	u := user.User{
		ID:       user.UserID(uid),
		Email:    "client@haha.foo",
		Password: "testpassword",
		Role:     string(middleware.Client),
	}
	return r.Create(context.Background(), u)
}

func createUserToken(r user.Repository, uid user.UserID) (middleware.Token, error) {
	testToken := uuid.New().String()
	err := r.AddToken(context.Background(), uid, user.Token(testToken))
	if err != nil {
		return middleware.Token(""), err
	}
	return middleware.Token(testToken), nil
}

func createHouse(hr house.Repository) (house.House, error) {
	hdto := house.CreateHouseDTO{
		Address:   "somewhere",
		Year:      1898,
		Developer: "someone",
	}
	h, err := hr.Create(context.Background(), hdto)
	if err != nil {
		return house.House{}, err
	}
	return h, nil
}

func createTestFlats(fr flat.Repository) ([]flat.FlatDTO, error) {
	createFlats := []flat.CreateFlatDTO{
		{
			HouseID: 1,
			Price:   14_000_000,
			Rooms:   12,
		},
		{
			HouseID: 1,
			Price:   12_000_000,
			Rooms:   18,
		},
		{
			HouseID: 1,
			Price:   4_500_000,
			Rooms:   88,
		},
		{
			HouseID: 1,
			Price:   4_800_000,
			Rooms:   2,
		},
	}
	flats := make([]flat.FlatDTO, 0)
	for _, cfdto := range createFlats {
		f, err := fr.Create(context.Background(), cfdto)
		if err != nil {
			return flats, err
		}
		flats = append(flats, f)
	}
	flatStatuses := map[int]modstatus.ModerationStatus{
		1: modstatus.Created,
		2: modstatus.OnModeration,
		3: modstatus.Approved,
		4: modstatus.Declined,
	}
	createdFlats := make([]flat.FlatDTO, 0)
	for _, fl := range flats {
		var status modstatus.ModerationStatus
		status, ok := flatStatuses[fl.ID]
		if !ok {
			return flats, errors.New("error getting status")
		}
		fudto := flat.UpdateFlatStatusDTO{
			ID:      fl.ID,
			HouseID: fl.HouseID,
			Status:  status.String(),
		}
		cf, err := fr.Update(context.Background(), fudto)
		if err != nil {
			return flats, errors.New("error updating test flats")
		}
		createdFlats = append(createdFlats, cf)
	}
	return createdFlats, nil
}
