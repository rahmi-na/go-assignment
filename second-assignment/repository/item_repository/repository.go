package item_repository

import (
	"second-assignment/entity"
	"second-assignment/pkg/errs"
)

type Repository interface {
	GetItemsByCodes(itemCodes []string) ([]entity.Item, errs.Error)
}
