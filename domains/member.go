package domains

import "time"

type Member struct {
	Name      string
	Link      string
	Email     string
	CreatedAt *time.Time
}

type Repo struct {
	Name      string
	Link      string
	CreatedAt *time.Time
}
