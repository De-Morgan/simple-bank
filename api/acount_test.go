package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	mock_db "github.com/morgan/simplebank/db/mock"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestGetAccountById(t *testing.T) {
	account := randomAccount()
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	scenerios := []testCase{
		{
			name: "Get Account",
			id:   account.ID,
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireAccountEquatCheck(t, recorder, account)
			},
		},
		{
			name: "Not Found",
			id:   account.ID,
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrNoRows)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{
			name: "Internal Server Error",
			id:   account.ID,
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{
			name: "Bad Request",
			id:   0,
			setup: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}
	for _, test := range scenerios {
		tN := test.name
		t.Run(tN, func(t *testing.T) {
			mockStore := mock_db.NewMockStore(cntrl)
			server := NewServer(mockStore)
			recorder := httptest.NewRecorder()
			test.setup(mockStore)

			url := fmt.Sprintf("/accounts/%d", test.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			test.checkResponse(t, recorder)

		})
	}

}

func requireAccountEquatCheck(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
	body := recorder.Body
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var result db.Account
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	require.Equal(t, result, account)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    utils.RandomName(),
		Balance:  utils.RandomMoney(),
		Currency: "NGN",
	}
}
