package infraestructures

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"net/url"
	"quantito-engine/internal/config"
	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQueryJson(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		query := "MYQUERY"
		input := struct{ Name string }{"TestObject"}
		inputJson := []byte(`{"Name":"TestObject"}`)

		mockConnection.On("Get", mock.Anything, query, []any{inputJson}).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]byte)
			mockedOutput := struct{ Output string }{":)"}
			jsonBytes, _ := json.Marshal(mockedOutput)
			*arg = jsonBytes
		})

		expectedResult := []byte(`{"Output":":)"}`)

		result, err := db.QueryJson(query, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	t.Run("Fails if error generating JSON from input", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		result, err := db.QueryJson("MYQUERY", math.Inf(1)) // math.Inf(1) causes json.Marshal to return an error

		assert.Error(t, err)
		assert.Nil(t, result)
		mockConnection.AssertNotCalled(t, "Get")
	})

	t.Run("Fails if DB error executing query", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		mockConnection.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("ERROR"))
		result, err := db.QueryJson("MYQUERY", 5)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Err = nil if no rows in result", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		mockConnection.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
		result, err := db.QueryJson("MYQUERY", 5)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestQueryAll(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		query := "MYQUERY"
		dbResponse := []int{1, 2, 3}

		var dest []int

		mockConnection.EXPECT().
			Select(mock.AnythingOfType("*[]int"), query).
			Return(nil).
			Once().
			RunFn = utils.MockWriteToPtr(0, dbResponse)

		err := db.QueryAll(query, &dest)

		assert.NoError(t, err)
		assert.Equal(t, dbResponse, dest)
	})

	t.Run("Fails if error executing query", func(t *testing.T) {
		mockConnection := mocks.NewMockIDB(t)
		db := &mySQLHandler{mockConnection}

		query := "MYQUERY"

		var dest []int

		mockConnection.EXPECT().
			Select(mock.AnythingOfType("*[]int"), query).
			Return(errors.New("ERROR")).
			Once()

		err := db.QueryAll(query, &dest)

		assert.Error(t, err)
		assert.Equal(t, 0, len(dest))
	})
}

func TestBuildDSN(t *testing.T) {
	dbConfig := &config.DBConfig{
		Username: "user",
		Password: "123456",
		Schema:   "mydb",
		Host:     "dbhost:3308",
	}
	t.Run("Contains all configuration options", func(t *testing.T) {
		dsn := buildMysqlDSN(dbConfig)

		assert.Contains(t, dsn, "user:123456@tcp(dbhost:3308)/mydb")
		assert.Contains(t, dsn, "interpolateParams=true")
		assert.Contains(t, dsn, "parseTime=true")
	})

	t.Run("UTC timezone", func(t *testing.T) {
		time.Local = time.UTC
		dsn := buildMysqlDSN(dbConfig)
		assert.Contains(t, dsn, "time_zone="+url.QueryEscape(`'+00:00'`))

	})

	t.Run("Timezone < UTC", func(t *testing.T) {
		time.Local = time.FixedZone("ART", int((-3 * time.Hour).Seconds()))
		dsn := buildMysqlDSN(dbConfig)
		assert.Contains(t, dsn, "loc=ART")
		assert.Contains(t, dsn, "time_zone="+url.QueryEscape(`'-03:00'`))

	})

	t.Run("Timezone > UTC", func(t *testing.T) {
		time.Local = time.FixedZone("BT", int(time.Hour.Seconds()))
		dsn := buildMysqlDSN(dbConfig)
		assert.Contains(t, dsn, "loc=BT")
		assert.Contains(t, dsn, "time_zone="+url.QueryEscape(`'+01:00'`))
	})
}
