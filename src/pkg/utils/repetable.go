package repeatable

import (
	"github.com/google/uuid"
	"time"
)

type Pagination struct {
	Limit, Offset int
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--
			continue
		}
		return nil
	}
	return
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
