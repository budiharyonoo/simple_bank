package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/budiharyonoo/simple_bank/db/mock"
	db "github.com/budiharyonoo/simple_bank/db/sqlc"
	"github.com/budiharyonoo/simple_bank/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccountById(t *testing.T) {
	account := randomAccount()

	// Create the test cases scenario for each HTTP Response Status Code
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().

					// Run Test for GetAccount query from Mock DB
					// The args should be the same with the original GetAccount method
					GetAccount(gomock.Any(), account.ID).

					// Run the test for 1 time (Connect to Mock DB)
					Times(1).

					// Tell gomock that the return values should be account,
					// The args should be the same with the return value of original GetAccount method
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the HTTP Response Status Code
				require.Equal(t, http.StatusOK, recorder.Code)

				requireMatchBodyResponse(t, recorder.Body, account)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().

					// Run Test for GetAccount query from Mock DB
					// The args should be the same with the original GetAccount method
					GetAccount(gomock.Any(), gomock.Any()).

					// No need to communicate with Mock DB
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the HTTP Response Status Code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().

					// Run Test for GetAccount query from Mock DB
					// The args should be the same with the original GetAccount method
					GetAccount(gomock.Any(), account.ID).

					// No need to communicate with Mock DB
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the HTTP Response Status Code
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.
					EXPECT().

					// Run Test for GetAccount query from Mock DB
					// The args should be the same with the original GetAccount method
					GetAccount(gomock.Any(), account.ID).

					// No need to communicate with Mock DB
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Check the HTTP Response Status Code
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	// Run each testCases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init gomock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			// Build stubs
			tc.buildStubs(store)

			// Start Mock HTTP Server
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Create new HTTP Request to the Mock Server
			url := fmt.Sprintf("/v1/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			// Check error
			require.NoError(t, err)

			// The Response saved in the recorder var,
			// and The Request saved in request var
			server.router.ServeHTTP(recorder, request)

			// Run the test
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Account {
	return db.Account{
		ID:        utils.RandomInt(1, 1000),
		Owner:     utils.RandomOwner(),
		Balance:   utils.RandomMoney(),
		Currency:  utils.RandomCurrency(),
		CreatedAt: time.Time{},
	}
}

func requireMatchBodyResponse(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	if err != nil {
		log.Fatalln("Error read body from response:", err)
		return
	}

	var respAcc db.Account
	err = json.Unmarshal(data, &respAcc)
	if err != nil {
		log.Fatalln("Error parse response body to JSON:", err)
		return
	}

	// === Test the data ===
	require.NotEmpty(t, respAcc)
	require.Equal(t, account.ID, respAcc.ID)
	require.Equal(t, account.Owner, respAcc.Owner)
	require.Equal(t, account.Balance, respAcc.Balance)
	require.Equal(t, account.Currency, respAcc.Currency)
}
