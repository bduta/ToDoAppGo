package models

import (
	"fmt"
)

type ToDoItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (item ToDoItem) ToFileFormat() string {
	return fmt.Sprintf("%d,%s,%s\n", item.Id, item.Name, item.Description)
}
