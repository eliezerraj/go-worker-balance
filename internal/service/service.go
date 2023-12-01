package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/go-worker-balance/internal/repository/postgre"
	"github.com/go-worker-balance/internal/core"
	"github.com/go-worker-balance/internal/erro"
	"github.com/aws/aws-xray-sdk-go/xray"

)

var childLogger = log.With().Str("service", "service").Logger()

type WorkerService struct {
	workerRepository 		*db_postgre.WorkerRepository
}

func NewWorkerService(workerRepository *db_postgre.WorkerRepository) *WorkerService{
	childLogger.Debug().Msg("NewWorkerService")

	return &WorkerService{
		workerRepository:	workerRepository,
	}
}

func (s WorkerService) Add(ctx context.Context, balance core.Balance) (error){
	childLogger.Debug().Msg("Add")

	_, root := xray.BeginSubsegment(ctx, "Service.Add")
	defer root.Close(nil)

	tx, err := s.workerRepository.StartTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = s.workerRepository.Get(ctx, balance)
	if err == erro.ErrNotFound {
		_, err := s.workerRepository.Add(ctx, balance)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		_, err := s.workerRepository.Update(ctx, balance)
		if err != nil {
			return err
		}
	}
	return nil
}
