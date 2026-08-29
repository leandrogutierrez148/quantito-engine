package interfaces

type IDBHandler interface {
	QueryScalar(query string, dest any, params ...interface{}) error
	QueryOne(query string, dest any, params ...interface{}) error
	QueryAll(query string, dest any, params ...interface{}) error
	QueryJson(query string, obj any) ([]byte, error)
}
