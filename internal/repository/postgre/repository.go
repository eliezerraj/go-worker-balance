package db_postgre

import (
	"context"
	"time"
	"errors"
	"database/sql"
	
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/go-worker-balance/internal/core"
	"github.com/go-worker-balance/internal/erro"
	"github.com/aws/aws-xray-sdk-go/xray"

)

var childLogger = log.With().Str("repository", "WorkerRepository").Logger()

type WorkerRepository struct {
	databaseHelper DatabaseHelper
}

func NewWorkerRepository(databaseHelper DatabaseHelper) WorkerRepository {
	childLogger.Debug().Msg("NewWorkerRepository")
	return WorkerRepository{
		databaseHelper: databaseHelper,
	}
}

func (w WorkerRepository) StartTx(ctx context.Context) (*sql.Tx, error) {
	childLogger.Debug().Msg("StartTx")

	client := w.databaseHelper.GetConnection()

	tx, err := client.BeginTx(ctx, &sql.TxOptions{})
    if err != nil {
        return nil, errors.New(err.Error())
    }

	return tx, nil
}

func (w WorkerRepository) Ping(ctx context.Context) (bool, error) {
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")
	childLogger.Debug().Msg("Ping")
	childLogger.Debug().Msg("++++++++++++++++++++++++++++++++")

	client := w.databaseHelper.GetConnection()

	err := client.PingContext(ctx)
	if err != nil {
		return false, errors.New(err.Error())
	}

	return true, nil
}

func (w WorkerRepository) Add(ctx context.Context, balance core.Balance) (*core.Balance, error){
	childLogger.Debug().Msg("Add")

	_, root := xray.BeginSubsegment(ctx, "SQL.Add-Balance-CDC")
	defer func() {
		root.Close(nil)
	}()

	client := w.databaseHelper.GetConnection()

	userLastUpdate := "CDC"

	stmt, err := client.Prepare(`INSERT INTO balance_cdc ( 	account_id, 
															person_id, 
															currency,
															amount,
															create_at,
															tenant_id,
															user_last_update) 
									VALUES($1, $2, $3, $4, $5, $6, $7) `)
	if err != nil {
		childLogger.Error().Err(err).Msg("INSERT statement")
		return nil, errors.New(err.Error())
	}
	
	_, err = stmt.ExecContext(	ctx,	
								balance.AccountID, 
								balance.PersonID,
								balance.Currency,
								balance.Amount,
								time.Now(),
								balance.TenantID,
								userLastUpdate)
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return nil, errors.New(err.Error())
	}
	defer stmt.Close()

	return &balance , nil
}

func (w WorkerRepository) Get(ctx context.Context, balance core.Balance) (*core.Balance, error){
	childLogger.Debug().Msg("Get")

	_, root := xray.BeginSubsegment(ctx, "SQL.Get-Balance-CDC")
	defer func() {
		root.Close(nil)
	}()

	client := w.databaseHelper.GetConnection()

	result_query := core.Balance{}
	rows, err := client.QueryContext(ctx, `SELECT id, account_id, person_id, currency, amount, create_at, update_at, tenant_id, user_last_update FROM balance_cdc WHERE account_id =$1`, balance.AccountID)
	if err != nil {
		childLogger.Error().Err(err).Msg("Query statement")
		return nil, errors.New(err.Error())
	}

	for rows.Next() {
		err := rows.Scan( &result_query.ID, 
							&result_query.AccountID, 
							&result_query.PersonID, 
							&result_query.Currency,
							&result_query.Amount,
							&result_query.CreateAt,
							&result_query.UpdateAt,
							&result_query.TenantID,
							&result_query.UserLastUpdate,
							)
		if err != nil {
			childLogger.Error().Err(err).Msg("Scan statement")
			return nil, errors.New(err.Error())
        }
		return &result_query, nil
	}
	defer rows.Close()

	return nil, erro.ErrNotFound
}

func (w WorkerRepository) Update(ctx context.Context, balance core.Balance) (bool, error){
	childLogger.Debug().Msg("Update...")

	_, root := xray.BeginSubsegment(ctx, "SQL.Update-Balance-CDC")
	defer func() {
		root.Close(nil)
	}()

	client := w.databaseHelper.GetConnection()

	userLastUpdate := "CDC"
	
	stmt, err := client.Prepare(`Update balance_cdc
									set account_id = $1, 
										person_id = $2, 
										currency = $3, 
										amount = $4, 
										update_at = $5,
										user_last_update =$7
								where id = $6 `)
	if err != nil {
		childLogger.Error().Err(err).Msg("UPDATE statement")
		return false, errors.New(err.Error())
	}

	result, err := stmt.ExecContext(ctx,	
									balance.AccountID, 
									balance.PersonID,
									balance.Currency,
									balance.Amount,
									time.Now(),
									balance.ID,
									userLastUpdate,
								)
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return false, errors.New(err.Error())
	}
	defer stmt.Close()

	rowsAffected, _ := result.RowsAffected()
	childLogger.Debug().Int("rowsAffected : ",int(rowsAffected)).Msg("")

	return true , nil
}

func (w WorkerRepository) Delete(ctx context.Context, balance core.Balance) (bool, error){
	childLogger.Debug().Msg("Delete")

	_, root := xray.BeginSubsegment(ctx, "SQL.Delete-Balance-CDC")
	defer func() {
		root.Close(nil)
	}()
	
	client := w.databaseHelper.GetConnection()

	stmt, err := client.Prepare(`Delete from balance_cdc where id = $1 `)
	if err != nil {
		childLogger.Error().Err(err).Msg("DELETE statement")
		return false, errors.New(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,balance.ID )
	if err != nil {
		childLogger.Error().Err(err).Msg("Exec statement")
		return false, errors.New(err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	childLogger.Debug().Int("rowsAffected : ",int(rowsAffected)).Msg("")

	return true , nil
}
