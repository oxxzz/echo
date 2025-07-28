package types

import "time"

type Entity struct {
	ID        int64     `db:"id"`
	UUID      string    `db:"uuid"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
