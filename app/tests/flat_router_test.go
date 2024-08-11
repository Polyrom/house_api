package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/Polyrom/houses_api/internal/flat"
	"github.com/Polyrom/houses_api/internal/house"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/internal/modstatus"
	"github.com/gorilla/mux"
)

const testFailedFlatID = -1

func executeRequest(r *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestCreateFlat(t *testing.T) {
	ctx := &testContext{
		Server:         newTestServer(),
		ModeratorToken: middleware.Token(""),
		ClientToken:    middleware.Token(""),
		Houses:         map[int]house.House{},
	}
	ctx.setup()
	var hs house.House
	hs, ok := ctx.Houses[1]
	if !ok {
		t.Errorf("failed to get test house")
	}
	var priceOne, roomsOne int = 41_000_000, 2
	flatReqBodyOne := flat.CreateFlatDTO{
		HouseID: hs.ID,
		Price:   priceOne,
		Rooms:   roomsOne,
	}
	expectedRespOne := flat.FlatDTO{
		ID:      1,
		HouseID: hs.ID,
		Price:   priceOne,
		Rooms:   roomsOne,
		Status:  modstatus.Created.String(),
	}
	var priceTwo, roomsTwo int = 4_000_000, 4
	flatReqBodyTwo := flat.CreateFlatDTO{
		HouseID: hs.ID,
		Price:   priceTwo,
		Rooms:   roomsTwo,
	}
	expectedRespTwo := flat.FlatDTO{
		ID:      2,
		HouseID: hs.ID,
		Price:   priceTwo,
		Rooms:   roomsTwo,
		Status:  modstatus.Created.String(),
	}
	priceThree := 4_000_000
	flatReqBodyBadReq := flat.CreateFlatDTO{
		HouseID: hs.ID,
		Price:   priceThree,
	}
	expectedRespBadReq := flat.FlatDTO{
		ID: testFailedFlatID,
	}
	var priceFour, roomsInvalid int = 4_000_000, -1
	flatReqBodyRoomsInvalid := flat.CreateFlatDTO{
		HouseID: hs.ID,
		Price:   priceFour,
		Rooms:   roomsInvalid,
	}
	var priceInvalid, roomsFour int = -1, 1
	flatReqBodyPriceInvalid := flat.CreateFlatDTO{
		HouseID: hs.ID,
		Price:   priceInvalid,
		Rooms:   roomsFour,
	}
	type args struct {
		authToken middleware.Token
		fdto      flat.CreateFlatDTO
	}
	type want struct {
		code int
		body flat.FlatDTO
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{name: "create flat client", args: args{ctx.ClientToken, flatReqBodyOne}, want: want{http.StatusOK, expectedRespOne}, wantErr: false},
		{name: "create flat moderator", args: args{ctx.ModeratorToken, flatReqBodyTwo}, want: want{http.StatusOK, expectedRespTwo}, wantErr: false},
		{name: "create flat bad request", args: args{ctx.ModeratorToken, flatReqBodyBadReq}, want: want{http.StatusBadRequest, expectedRespBadReq}, wantErr: false},
		{name: "create flat invalid rooms number", args: args{ctx.ClientToken, flatReqBodyRoomsInvalid}, want: want{http.StatusBadRequest, expectedRespBadReq}, wantErr: false},
		{name: "create flat invalid price", args: args{ctx.ClientToken, flatReqBodyPriceInvalid}, want: want{http.StatusBadRequest, expectedRespBadReq}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.args.fdto)
			if err != nil {
				t.Errorf("failed to marshal CreateFlatDTO: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, "/flat/create", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Errorf("failed to create Create Flat request: %v", err)
			}
			req.Header.Set("Authorization", string(tt.args.authToken))
			resp := executeRequest(ctx.Server.Router, req)
			if resp.Code != tt.want.code {
				t.Errorf("expected response code %d. Got %d\n", tt.want.code, resp.Code)
			}
			body := resp.Body.String()
			var got flat.FlatDTO
			err = json.Unmarshal([]byte(body), &got)
			if err != nil {
				t.Errorf("failed to unmarshal create flat response body: %v", err)
			}
			if !reflect.DeepEqual(expectedRespBadReq, tt.want.body) && !reflect.DeepEqual(got, tt.want.body) {
				t.Errorf("create flat = %v, want %v", got, tt.want.body)
			}
		})
	}
	t.Cleanup(func() {
		ctx.cleanup()
	})
}

func TestGetHouse(t *testing.T) {
	ctx := &testContext{
		Server:         newTestServer(),
		ModeratorToken: middleware.Token(""),
		ClientToken:    middleware.Token(""),
		Houses:         map[int]house.House{},
	}
	ctx.setup()
	var hs house.House
	hs, ok := ctx.Houses[1]
	if !ok {
		t.Errorf("failed to get test house")
	}
	fr := flat.NewRepository(ctx.Server.DB, &MockLogger{})
	expectedRespModer, err := createTestFlats(fr)
	if err != nil {
		t.Errorf("failed to create test flats: %v", err)
	}
	expectedRespClient := make([]flat.FlatDTO, 0)
	for _, f := range expectedRespModer {
		if f.Status == modstatus.Approved.String() {
			expectedRespClient = append(expectedRespClient, f)
		}
	}
	nonExistentID := 234
	type args struct {
		authToken middleware.Token
		hid       int
	}
	type want struct {
		code int
		body []flat.FlatDTO
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{name: "get house moderator", args: args{ctx.ModeratorToken, hs.ID}, want: want{http.StatusOK, expectedRespModer}, wantErr: false},
		{name: "get house client", args: args{ctx.ClientToken, hs.ID}, want: want{http.StatusOK, expectedRespClient}, wantErr: false},
		{name: "get house id not found", args: args{ctx.ClientToken, nonExistentID}, want: want{http.StatusOK, []flat.FlatDTO{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hidStr := strconv.Itoa(tt.args.hid)
			req, err := http.NewRequest(http.MethodGet, "/house/"+hidStr, nil)
			if err != nil {
				t.Errorf("failed to create get house request: %v", err)
			}
			req.Header.Set("Authorization", string(tt.args.authToken))
			resp := executeRequest(ctx.Server.Router, req)
			if resp.Code != tt.want.code {
				t.Errorf("expected response code %d. Got %d\n", tt.want.code, resp.Code)
			}
			body := resp.Body.String()
			var got []flat.FlatDTO
			err = json.Unmarshal([]byte(body), &got)
			if err != nil {
				t.Errorf("failed to unmarshal get house response body: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want.body) {
				t.Errorf("get house = %v, want %v", got, tt.want.body)
			}
		})
	}
	t.Cleanup(func() {
		ctx.cleanup()
	})
}
