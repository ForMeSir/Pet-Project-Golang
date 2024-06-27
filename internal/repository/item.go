package repository

import (
	"database/sql"
	"fmt"
	"pet/internal/model"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemService struct {
	db *sqlx.DB
}

func NewItemService(db *sqlx.DB) *ItemService {
	return &ItemService{db: db}
}

func (i *ItemService) Create(item model.Item) (uuid.UUID, error) {
	var err error
	item.Id, err = uuid.NewUUID()
	if err != nil {
		fmt.Println(err)
		return item.Id, err
	}
	query := fmt.Sprintf("INSERT INTO %s (id,title,description,price,image) values($1,$2,$3,$4,$5)", itemTable)
	_, err = i.db.Exec(query, item.Id, item.Title, item.Description, item.Price, item.Image)
	if err != nil {
		fmt.Println(err)
		return item.Id, err
	}

	return item.Id, err
}

func (i *ItemService) Delete(id uuid.UUID) (err error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", itemTable)
	_, err = i.db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (i *ItemService) FindByTitle(title string, limit int, offset int) (items []model.Item, err error) {
	var rows *sql.Rows
	query := fmt.Sprint("SELECT * FROM ", itemTable, " WHERE title LIKE '%", title, "%' LIMIT $1 OFFSET $2")

	if limit == 0 {
		query := fmt.Sprint("SELECT * FROM ", itemTable, " WHERE title LIKE '%", title, "%' LIMIT ALL OFFSET $1")
		rows, err = i.db.Query(query, offset)
	} else {
		rows, err = i.db.Query(query, limit, offset)
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var item model.Item
		if err = rows.Scan(&item.Id, &item.Title, &item.Description, &item.Price, &item.Image); err != nil {
			fmt.Println(err)
		}
		items = append(items, item)
	}
	return
}

func (i *ItemService) Update(id uuid.UUID, item model.UpdateItem) (err error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *item.Title)
		argId++
	}

	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	if item.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *item.Price)
		argId++
	}
	if item.Image != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ",")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", itemTable, setQuery, id)

	args = append(args, id)

	_, err = i.db.Exec(query, args)

	return err
}
