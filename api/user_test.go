package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mock_db "github.com/morgan/simplebank/db/mock"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

type createUserTestCase struct {
	name          string
	body          CreateUserRequest
	setup         func(store *mock_db.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func TestCreateUser(t *testing.T) {

	cntrl := gomock.NewController(t)
	password := utils.RandomString(8)
	userName := utils.RandomName()
	fullName := utils.RandomName()
	email := utils.RandomEmail(6)
	hashPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       userName,
		FullName:       fullName,
		HashedPassword: hashPassword,
		Email:          email,
	}
	tests := []createUserTestCase{
		{
			name: "Create User",
			body: CreateUserRequest{
				Username: userName,
				Password: password,
				FullName: fullName,
				Email:    email,
			},
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), checkCreateUserParam(
					db.CreateUserParams{
						Username:       userName,
						FullName:       fullName,
						Email:          email,
						HashedPassword: hashPassword,
					}, password,
				)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var userRes CreateUserResponse
				require.Equal(t, http.StatusCreated, recorder.Code)
				body, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				err = json.Unmarshal(body, &userRes)
				require.NoError(t, err)
				require.NotEmpty(t, userRes)
				require.Equal(t, userRes.Email, email)
				require.Equal(t, userRes.FullName, fullName)

			},
		},
		{
			name: "Invalid Email",
			body: CreateUserRequest{
				Username: userName,
				Password: password,
				FullName: fullName,
				Email:    "michael",
			},
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), checkCreateUserParam(db.CreateUserParams{
					Username:       userName,
					FullName:       fullName,
					Email:          "michael",
					HashedPassword: hashPassword,
				}, password)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := mock_db.NewMockStore(cntrl)
			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()
			test.setup(mockStore)

			url := "/users"
			bodyByte, err := json.Marshal(test.body)
			require.NoError(t, err)
			r := bytes.NewReader(bodyByte)
			request, err := http.NewRequest(http.MethodPost, url, r)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			test.checkResponse(t, recorder)
		})
	}
}

func checkCreateUserParam(userParam db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamMater{userParm: userParam, password: password}
}

type eqCreateUserParamMater struct {
	userParm db.CreateUserParams
	password string
}

func (e eqCreateUserParamMater) Matches(x interface{}) bool {

	// Check if types assignable and convert them to common type
	if val, ok := x.(db.CreateUserParams); ok {
		err := utils.CheckPasswordCorrect(e.password, val.HashedPassword)
		if err != nil {
			return false
		}
		e.userParm.HashedPassword = val.HashedPassword
		return reflect.DeepEqual(e.userParm, val)
	}

	return false
}

func (e eqCreateUserParamMater) String() string {
	return fmt.Sprintf("is equal to %v (%T)", e.userParm, e.password)
}
