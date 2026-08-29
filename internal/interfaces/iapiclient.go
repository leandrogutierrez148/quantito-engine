package interfaces

import (
	"quantito-engine/internal/models"
)

type IAPIClient interface {
	Get(path string, result any, options ...models.APIClientOptionsFunc) error
	Post(path string, result any, options ...models.APIClientOptionsFunc) error
	Put(path string, result any, options ...models.APIClientOptionsFunc) error
	Patch(path string, result any, options ...models.APIClientOptionsFunc) error
	Delete(path string, result any, options ...models.APIClientOptionsFunc) error
	SetLoginFn(func() (string, error))
}
