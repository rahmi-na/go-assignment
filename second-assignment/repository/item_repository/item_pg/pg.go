package item_pg

import (
	"database/sql"
	"second-assignment/entity"
	"second-assignment/pkg/errs"
	"second-assignment/repository/item_repository"
)

type itemPG struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) item_repository.Repository {
	return &itemPG{
		db: db,
	}
}

func (itemPG *itemPG) GetItemsByCodes(itemCodes []string) ([]entity.Item, errs.Error) {
	query := itemPG.generateGetItemsByCodesQuery(len(itemCodes))

	itemCodeDbArg := []any{}

	for _, itemCode := range itemCodes {
		itemCodeDbArg = append(itemCodeDbArg, itemCode)
	}

	rows, err := itemPG.db.Query(query, itemCodeDbArg...)

	if err != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	items := []entity.Item{}

	for rows.Next() {

		item := entity.Item{}

		err = rows.Scan(&item.ItemId, &item.ItemCode, &item.Quantity, &item.Description, &item.OrderId, &item.CreatedAt)

		if err != nil {
			return nil, errs.NewInternalServerError("something went wrong")
		}

		items = append(items, item)
	}

	return items, nil
}
