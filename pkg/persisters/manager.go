package persisters

//go:generate sqlc -f ../../sqlc.yaml generate

import (
	"database/sql"

	"github.com/pojntfx/aeolius/pkg/migrations"
	"github.com/pojntfx/aeolius/pkg/models"
	"github.com/pressly/goose/v3"

	_ "github.com/lib/pq"
)

type ManagerPersister struct {
	pgaddr  string
	queries *models.Queries
	db      *sql.DB
}

func NewManagerPersister(pgaddr string) *ManagerPersister {
	return &ManagerPersister{
		pgaddr: pgaddr,
	}
}

func (p *ManagerPersister) Open() error {
	var err error
	p.db, err = sql.Open("postgres", p.pgaddr)
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(p.db, "."); err != nil {
		return err
	}

	p.queries = models.New(p.db)

	return nil
}

func (p *ManagerPersister) Close() error {
	if p.db != nil {
		_ = p.db.Close()
	}

	return nil
}
