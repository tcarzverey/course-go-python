package handler

import (
	"context"
	stderrors "errors"

	herr "github.com/tcarzverey/course-go-python/homeworks/hw1/handler/errors"
)

type Handler struct {
	db UsersDB
}

func NewHandler(db UsersDB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) UpdateUserBalance(ctx context.Context, userID, balance int64) error {
	err := h.db.UpdateBalance(ctx, userID, balance)
	if err == nil {
		return nil
	}

	var rerr *herr.RetryableError
	if stderrors.As(err, &rerr) {
		attempts := rerr.RetryCount()
		for i := 0; i < attempts; i++ {
			if e := h.db.UpdateBalance(ctx, userID, balance); e == nil {
				return nil
			} else {
				err = e
			}
		}
		return err
	}

	var addErr *herr.AdditionalMessageError
	if stderrors.As(err, &addErr) {
		return err
	}

	var nf *herr.NotFoundError
	if stderrors.As(err, &nf) {
		return herr.NewAdditionalMessageError(err, "not found")
	}

	panic(err)
}
