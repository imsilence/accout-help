package domains

import "time"

type Member struct {
	ID        int64
	Name      string
	Link      string
	Email     string
	CreatedAt *time.Time
}

type Repo struct {
	ID        int64
	Name      string
	Link      string
	CreatedAt *time.Time
}
