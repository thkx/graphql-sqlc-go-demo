package integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth"
	"github.com/thkx/graphql-sqlc-go-dome/internal/graph/model"
	"github.com/thkx/graphql-sqlc-go-dome/internal/repository"
	"github.com/thkx/graphql-sqlc-go-dome/internal/usecase"
)

type indicator struct {
	repo repository.IndicatorRepository
}

func NewUsecase(repo repository.IndicatorRepository) usecase.IndicatorUsecase {
	return &indicator{repo: repo}
}

func (i *indicator) CreateIndicator(ctx context.Context, input model.IndicatorInput) error {
	_, err := auth.RequireAuth(ctx)
	if err != nil {
		return err
	}
	ind := &repository.Indicator{
		ID:            uuid.New(),
		Indicator:     input.Indicator,
		IndicatorType: input.IndicatorType,
		MetaSource:    input.MetaSource,
	}
	return i.repo.CreateIndicator(ctx, ind)
}

func (i *indicator) ListIndicators(ctx context.Context) ([]*model.Indicator, error) {
	params := &repository.Indicator{}
	pagination := &repository.Pagination{
		Limit:  0,
		Offset: 10,
	}
	rows, err := i.repo.ListIndicators(ctx, params, pagination)
	if err != nil {
		return nil, err
	}

	var result []*model.Indicator
	for _, v := range rows {
		ind := &model.Indicator{
			ID:            v.ID.String(),
			Indicator:     v.Indicator,
			IndicatorType: v.IndicatorType,
			MetaSource:    v.MetaSource,
		}

		result = append(result, ind)
	}
	return result, nil
}
