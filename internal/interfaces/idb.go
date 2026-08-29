package interfaces

type IDB interface {
	Select(dest any, query string, args ...any) error
	Get(dest interface{}, query string, args ...interface{}) error
}
