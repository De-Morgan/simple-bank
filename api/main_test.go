package api

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	mock_db "github.com/morgan/simplebank/db/mock"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())

}

type testCase struct {
	name          string
	id            interface{}
	setup         func(store *mock_db.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}
