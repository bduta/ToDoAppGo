package engine

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"newtodoapp/models"
	"os"
	"slices"
	"strconv"
	"strings"
)

var ToDoListFileName string = "ToDoList.txt"

func createTheToDoListFileIfNeeded() (bool, error) {
	creationRequired := false
	_, err := os.Stat(ToDoListFileName)
	if os.IsNotExist(err) {
		f, err := os.Create(ToDoListFileName)
		defer f.Close()
		creationRequired = true
		if err != nil {
			return true, errors.New("file does not exist and could not be created")
		}
	}
	return creationRequired, nil
}

func readExistingList() (list []models.ToDoItem, err error) {

	_, fileErr := os.Stat(ToDoListFileName)
	if os.IsNotExist(fileErr) {
		return []models.ToDoItem{}, errors.New("ToDo file does not exist")
	}

	file, fileErr := os.Open(ToDoListFileName)
	if fileErr != nil {
		return []models.ToDoItem{}, errors.New("file could not be opened")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var toDos []models.ToDoItem
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return []models.ToDoItem{}, errors.New("Line has incorrect format: " + scanner.Text())
		}

		toDoId, err := strconv.Atoi(parts[0])
		if err != nil {
			return []models.ToDoItem{}, errors.New("Id could not be converted to int: " + scanner.Text())
		}

		toDo := models.ToDoItem{
			Id:          toDoId,
			Name:        parts[1],
			Description: parts[2],
		}
		toDos = append(toDos, toDo)
	}

	slices.SortFunc(toDos, func(i, j models.ToDoItem) int {
		return i.Id - j.Id
	})

	return toDos, nil
}

func generateItemId(fileCreationRequired bool) (id int, err error) {
	if !fileCreationRequired {
		toDos, err := readExistingList()
		if err != nil {
			return -1, err
		}

		if len(toDos) > 0 {
			lastToDo := toDos[len(toDos)-1]
			return lastToDo.Id + 1, nil
		} else {
			return 1, nil
		}

	} else {
		return 1, nil
	}
}

func getIndexBasedOnId(toDos []models.ToDoItem, id int) (index int, err error) {

	flagIdIndex := -1
	for index, item := range toDos {
		if item.Id == id {
			flagIdIndex = index
		}
	}

	if flagIdIndex == -1 {
		return -1, errors.New("flag Id could not be found")
	}

	return flagIdIndex, nil
}

func writeItemToFile(item models.ToDoItem) error {

	_, err := os.Stat(ToDoListFileName)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("ToDoList file does not exist")
	} else {

		f, err := os.OpenFile(ToDoListFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if err != nil {
			return errors.New("could not open the ToDoList file")
		}

		if _, err := f.WriteString(item.ToFileFormat()); err != nil {
			return errors.New("could not append the new toDo item to the ToDoList file: " + item.ToFileFormat())
		}
	}

	return nil
}

func writeItemsToFile(items []models.ToDoItem) error {
	f, err := os.Create(ToDoListFileName)
	if err != nil {
		return errors.New("could not open the ToDoList file")
	}
	defer f.Close()

	for _, item := range items {
		if _, err := f.WriteString(item.ToFileFormat()); err != nil {
			return errors.New("could not append the toDo item to the ToDoList file: " + item.ToFileFormat())
		}
	}

	return nil
}

func ExecuteCommand(arguments []string) error {
	if len(arguments) == 0 {
		toDos, err := readExistingList()
		if err != nil {
			return errors.New("Error reading existing list: " + err.Error())
		}

		for _, item := range toDos {
			fmt.Printf("Id:%d, ToDo:%s, Description:%s\n", item.Id, item.Name, item.Description)
		}

		return nil
	}

	op := flag.String("op", "", "operation to be performed")
	id := flag.Int("id", -1, "to do id")
	name := flag.String("name", "", "to do name")
	description := flag.String("description", "", "to do description")

	flag.Parse()

	switch strings.ToLower(*op) {
	case "add":
		if *name == "" {
			return errors.New("missing flag name")
		}

		if *description == "" {
			return errors.New("missing flag description")
		}

		err := CreateItem(*name, *description)
		if err != nil {
			return errors.New("Error creating item: " + err.Error())
		}

		fmt.Println("Item added successfully")
	case "update":

		if *id == -1 {
			return errors.New("invalid flag id")
		}

		if *description == "" {
			return errors.New("missing flag description")
		}

		err := UpdateItem(*id, *description)
		if err != nil {
			return errors.New("Error updating item: " + err.Error())
		}

		fmt.Println("Item updated successfully")

	case "delete":
		if *id == -1 {
			return errors.New("invalid flag id")
		}

		err := DeleteItem(*id)
		if err != nil {
			return errors.New("Error deleting item: " + err.Error())
		}

		fmt.Println("Item deleted successfully")

	default:
		fmt.Println("The operation flag entered is not valid.")
		fmt.Println("To add a flag: -op=add -name=<FlagName> -description=<FlagDescription>")
		fmt.Println("To update a flag: -op=update -id=<FlagId> -description=<FlagDescription>")
		fmt.Println("To delete a flag: -op=delete -id=<FlagId>")
	}
	return nil
}

func GetItems() ([]models.ToDoItem, error) {
	toDos, err := readExistingList()
	if err != nil {
		return nil, errors.New("Error reading existing list: " + err.Error())
	}
	return toDos, nil
}

func CreateItem(name string, description string) error {
	fileCreationRequired, err := createTheToDoListFileIfNeeded()
	if err != nil {
		return errors.New("Error creating ToDo list file: " + err.Error())
	}

	newItem := models.ToDoItem{
		Name:        name,
		Description: description,
	}

	id, err := generateItemId(fileCreationRequired)
	if err != nil {
		return errors.New("Error generating item ID: " + err.Error())
	}
	newItem.Id = id

	writeItemError := writeItemToFile(newItem)
	if writeItemError != nil {
		return errors.New("Error writing item to file: " + writeItemError.Error())
	}

	return nil
}

func UpdateItem(id int, description string) error {
	toDos, err := readExistingList()
	if err != nil {
		return errors.New("Error reading existing list: " + err.Error())
	}

	index, err := getIndexBasedOnId(toDos, id)
	if err != nil {
		return errors.New("Error finding item by ID: " + err.Error())
	}

	toDos[index].Description = description

	overwritingFileErr := writeItemsToFile(toDos)
	if overwritingFileErr != nil {
		return errors.New("Error overwriting file: " + overwritingFileErr.Error())
	}

	return nil
}

func DeleteItem(id int) error {
	toDos, err := readExistingList()
	if err != nil {
		return errors.New("Error reading existing list: " + err.Error())
	}

	index, err := getIndexBasedOnId(toDos, id)
	if err != nil {
		return errors.New("Error finding item by ID: " + err.Error())
	}

	toDos = append(toDos[:index], toDos[index+1:]...)

	overwritingFileErr := writeItemsToFile(toDos)
	if overwritingFileErr != nil {
		return errors.New("Error overwriting file: " + overwritingFileErr.Error())
	}

	return nil
}
