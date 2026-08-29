package infraestructures

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"quantito-engine/internal/config"
	"quantito-engine/internal/interfaces"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

const driverName = "mysql"

type mySQLHandler struct {
	db interfaces.IDB
}

func NewMysqlDB(dbConfig *config.DBConfig) (interfaces.IDBHandler, error) {
	dsn := buildMysqlDSN(dbConfig)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConn)
	db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxTTLConn) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleTimeConn) * time.Second)

	// Verify connection credentials upfront to avoid checking during the first query
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Override default sqlx mapping to avoid having to add tags when using UpperCamelCase in the database.
	sqlx.NameMapper = strings.Clone

	logger := log.Logger.With().Str("component", "mysql").Logger()
	loggerAdapter := zerologadapter.New(logger)
	db = sqldblogger.OpenDriver(
		dsn,
		db.Driver(),
		loggerAdapter,
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
	)

	sqlxDb := sqlx.NewDb(db, driverName).Unsafe()

	return &mySQLHandler{sqlxDb}, nil
}

// QueryScalar executes a query expecting a scalar return value assigned to `dest`.
// Placeholders are replaced by the corresponding parameters.
//
// Usage example:
//
//	db.QueryScalar("SELECT Value FROM Company WHERE Param = ?", &result, "MYPARAM")
func (h *mySQLHandler) QueryScalar(query string, dest any, params ...any) error {
	return h.QueryOne(query, dest, params...)
}

// QueryOne executes a query expecting a single row and binds it to a struct pointed to by `dest`.
// If the query returns multiple rows, only the first one is used.
// The `db` tag can be used to specify the corresponding column name.
// Placeholders are replaced by the corresponding parameters.
//
// Usage example:
//
//	db.QueryOne("SELECT * FROM Company WHERE Param = ?", &result, "MYPARAM")
func (h *mySQLHandler) QueryOne(query string, dest any, params ...any) error {
	return h.db.Get(dest, query, params...)
}

// QueryAll executes a query and binds all rows to a slice pointer `dest`.
// The `db` tag can be used to specify the corresponding column name.
// Placeholders are replaced by the corresponding parameters.
//
// Usage example:
//
//	db.QueryAll("SELECT * FROM Company WHERE Param LIKE ?", &result, "W%")
func (h *mySQLHandler) QueryAll(query string, dest any, params ...any) error {
	return h.db.Select(dest, query, params...)
}

// QueryJson executes queries taking a struct input and returning a JSON payload.
// `obj` is any struct that can be marshaled with `json.Marshal`. The output is a `[]byte`
// that can be unmarshaled using `utils.UnmarshalJson`. In case of no rows, returns nil without error.
func (h *mySQLHandler) QueryJson(query string, obj any) ([]byte, error) {
	input, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var output []byte
	err = h.QueryScalar(query, &output, input)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

func buildMysqlDSN(dbConfig *config.DBConfig) string {
	config := dbConfig.ToMysqlConfig()
	config.InterpolateParams = true
	config.ParseTime = true

	config.Loc = time.Local
	_, offset := time.Now().Zone()
	config.Params = map[string]string{"time_zone": fmt.Sprintf("'%+03d:00'", offset/60/60)}

	return config.FormatDSN()
}
