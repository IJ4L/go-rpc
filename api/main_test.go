package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "simplebank.com/db/sqlgen"
	"simplebank.com/utils"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymetricKey:    utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	if err != nil {
		t.Fatal("cannot create server:", err)
	}

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
