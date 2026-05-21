package dbops

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresRepositoryImpl struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PostgresRepository {
	return &PostgresRepositoryImpl{
		db: db,
	}
}

func (p *PostgresRepositoryImpl) GetStatus() error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf(
			"failed to ping postgres: %w",
			err,
		)
	}

	return nil
}

///////////////////////////////////////////////////////////
//////////////////////// INSERT ///////////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) Insert(query string, args ...any) (*sql.Rows, error) {

	stmt, err := p.db.PrepareContext(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare insert statement: %w",
			err,
		)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(
		context.Background(),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute insert query: %w",
			err,
		)
	}

	return rows, nil
}

///////////////////////////////////////////////////////////
//////////////////////// UPDATE ///////////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) Update(query string, args ...any) (*sql.Rows, error) {

	stmt, err := p.db.PrepareContext(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare update statement: %w",
			err,
		)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(
		context.Background(),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute update query: %w",
			err,
		)
	}

	return rows, nil
}

///////////////////////////////////////////////////////////
//////////////////////// FETCH ////////////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) Fetch(query string, args ...any) (*sql.Rows, error) {

	stmt, err := p.db.PrepareContext(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare fetch statement: %w",
			err,
		)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(
		context.Background(),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute fetch query: %w",
			err,
		)
	}

	return rows, nil
}

///////////////////////////////////////////////////////////
//////////////////////// DELETE ///////////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) Delete(query string, args ...any) (*sql.Rows, error) {

	stmt, err := p.db.PrepareContext(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to prepare delete statement: %w",
			err,
		)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(
		context.Background(),
		args...,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to execute delete query: %w",
			err,
		)
	}

	return rows, nil
}

///////////////////////////////////////////////////////////
/////////////////////// BEGIN TX //////////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) BeginTx() (*sql.Tx, error) {

	tx, err := p.db.BeginTx(
		context.Background(),
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	return tx, nil
}

///////////////////////////////////////////////////////////
//////////////// BEGIN TX WITH OPTIONS ////////////////////
///////////////////////////////////////////////////////////

func (p *PostgresRepositoryImpl) BeginTxWithOptions(opts *sql.TxOptions) (*sql.Tx, error) {

	tx, err := p.db.BeginTx(
		context.Background(),
		opts,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin transaction with options: %w",
			err,
		)
	}

	return tx, nil
}