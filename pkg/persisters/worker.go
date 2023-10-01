package persisters

import (
	"database/sql"

	"github.com/pojntfx/skysweeper/pkg/models"

	_ "github.com/lib/pq"
)

type WorkerPersister struct {
	pgaddr  string
	queries *models.Queries
	db      *sql.DB
}

func NewWorkerPersister(pgaddr string) *WorkerPersister {
	return &WorkerPersister{
		pgaddr: pgaddr,
	}
}

func (p *WorkerPersister) Open() error {
	var err error
	p.db, err = sql.Open("postgres", p.pgaddr)
	if err != nil {
		return err
	}

	p.queries = models.New(p.db)

	return nil
}

func (p *WorkerPersister) Close() error {
	if p.db != nil {
		_ = p.db.Close()
	}

	return nil
}
