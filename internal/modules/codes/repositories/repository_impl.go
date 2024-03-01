package repositories

import "context"

func (r *repository) DeleteCode(ctx context.Context, code string) error {
	return r.database.Del(ctx, code).Err()
}

func (r *repository) GetUserIdByCode(ctx context.Context, code string) (string, error) {
	return r.database.Get(ctx, code).Result()
}

func (r *repository) StoreCode(ctx context.Context, userId string, code string) error {
	return r.database.Set(ctx, code, userId, 0).Err()
}
