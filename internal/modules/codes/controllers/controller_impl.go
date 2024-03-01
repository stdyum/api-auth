package controllers

import (
	"context"
	"github.com/google/uuid"
	"math/rand"
)

func (c *controller) CreateAndStoreCode(ctx context.Context, userId uuid.UUID) (string, error) {
	code := c.generateCode()
	return code, c.StoreCode(ctx, userId, code)
}

func (c *controller) DeleteCode(ctx context.Context, code string) error {
	return c.repository.DeleteCode(ctx, code)
}

func (c *controller) GetUserIdByCode(ctx context.Context, code string) (uuid.UUID, error) {
	id, err := c.repository.GetUserIdByCode(ctx, code)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}

func (c *controller) StoreCode(ctx context.Context, userId uuid.UUID, code string) error {
	return c.repository.StoreCode(ctx, userId.String(), code)
}

func (c *controller) generateCode() string {
	code := ""

	for i := 0; i < 7; i++ {
		if i == 3 {
			code += "."
			continue
		}

		randNum := rand.Intn(36)
		if randNum < 10 {
			code += string(rune('0' + randNum))
		} else {
			code += string(rune('A' + randNum - 10))
		}
	}

	return code
}
