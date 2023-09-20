package persisters

import (
	"context"

	"github.com/pojntfx/aeolius/pkg/models"
)

func (p *ManagerPersister) UpsertConfiguration(
	ctx context.Context,
	did string,
	service string,
	refreshJWT string,
	enabled bool,
	postTtl int32,
) (models.Configuration, error) {
	return p.queries.UpsertConfiguration(ctx, models.UpsertConfigurationParams{
		Did:        did,
		Service:    service,
		RefreshJwt: refreshJWT,
		Enabled:    enabled,
		PostTtl:    postTtl,
	})
}

func (p *WorkerPersister) GetConfigurations(
	ctx context.Context,
) ([]models.Configuration, error) {
	return p.queries.GetConfiguration(ctx)
}

func (p *WorkerPersister) UpdateRefreshToken(
	ctx context.Context,
	did string,
	refreshJWT string,
) error {
	return p.queries.UpdateConfigurationRefreshJWT(ctx, models.UpdateConfigurationRefreshJWTParams{
		RefreshJwt: refreshJWT,
		Did:        did,
	})
}
