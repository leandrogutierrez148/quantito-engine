package interfaces

import "time"

type IClock interface {
	Now() time.Time
}
