package server

var Items = make(map[string]string)

func GetItems() map[string]string {
	return Items
}

func UpdateItem(name string, description string) {
	Items[name] = description
}

func DeleteItem(name string) {
	delete(Items, name)
}
