package persisters

import (
	"context"

	"github.com/pojntfx/skysweeper/pkg/models"
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

func (p *ManagerPersister) GetConfiguration(
	ctx context.Context,
	did string,
) (models.Configuration, error) {
	return p.queries.GetConfiguration(ctx, did)
}

func (p *ManagerPersister) DeleteConfiguration(
	ctx context.Context,
	did string,
) error {
	return p.queries.DeleteConfiguration(ctx, did)
}

func (p *WorkerPersister) GetEnabledConfigurations(
	ctx context.Context,
) ([]models.Configuration, error) {
	return p.queries.GetEnabledConfigurations(ctx)
}

func (p *WorkerPersister) UpdateRefreshTokenAndCursor(
	ctx context.Context,
	did string,
	cursor string,
	refreshJWT string,
) error {
	return p.queries.UpdateConfigurationRefreshJWTAndCursor(ctx, models.UpdateConfigurationRefreshJWTAndCursorParams{
		RefreshJwt: refreshJWT,
		Cursor:     cursor,
		Did:        did,
	})
}
