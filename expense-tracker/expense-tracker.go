package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

type Expense struct {
	Id          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
}
type Expenses map[string]Expense

var Exp Expenses
var dataKeys []int
var amountPtr, descriptionPtr, categoryPrt string
var fs *flag.FlagSet

/*
Check if exists our data.json file, if wasn't, create a new one to store the expense's data
*/
func initialize() error {
	//check if the data.json file exists
	dir, err := os.ReadDir("./")
	if err != nil {
		return err
	}
	find := false
	for _, d := range dir {
		if d.Name() == "data.json" {
			find = true
			break
		}
	}
	if !find {
		f, err := os.Create("data.json")
		if err != nil {
			return err
		}
		defer f.Close()
	} else {
		dat, err := os.ReadFile("data.json")
		if err != nil {
			return err
		}
		err = json.Unmarshal(dat, &Exp)
		if err != nil {
			return err
		}
	}
	Exp.sortKeys()
	return nil
}

/*
Write data into data.json file
*/
func (e Expenses) writeFile() error {
	by, err := json.Marshal(e)
	if err != nil {
		return err
	}
	err = os.WriteFile("data.json", by, 0664)
	if err != nil {
		return err
	}
	return nil
}

/*
Organize the ID expense's key
*/
func (e Expenses) sortKeys() {
	for k := range e {
		i, err := strconv.Atoi(k)
		if err != nil {
			fmt.Println("Fail while sorting keys. ", err.Error())
			os.Exit(1)
		}
		dataKeys = append(dataKeys, i)
	}
	slices.Sort(dataKeys)
	fmt.Printf("Keys: %v\n", dataKeys)
}

/*
Validate if a Expense exists
*/
func (e Expenses) exists(id string) (Expense, bool) {
	v, bol := e[id]
	return v, bol
}

/*
Add a new Expense
*/
func (e Expenses) add(description, amount, category string) error {
	var idExpense int
	if len(dataKeys) == 0 {
		idExpense = 1
	} else {
		idExpense = dataKeys[len(dataKeys)-1] + 1
	}
	dataKeys = append(dataKeys, idExpense)
	t := time.Now()

	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return err
	}
	ex := Expense{
		Id:          idExpense,
		Description: description,
		Amount:      f,
		Category:    category,
		Date:        fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()),
	}

	Exp[strconv.Itoa(idExpense)] = ex
	Exp.writeFile()
	return nil
}

func regFlag(v *string, name, value, usage string) {
	if fs.Lookup(name) == nil {
		fs.StringVar(v, name, value, usage)
	}
}

func getFlag(name string) string {
	return fs.Lookup(name).Value.String()
}

func init() {
	fs = flag.NewFlagSet("", flag.ContinueOnError)
	initFlags()
	Exp = make(Expenses)
}

func initFlags() {
	//defining expense program flags
	regFlag(&amountPtr, "amount", "", "(mandatory) expense's amount, example 32.99")
	regFlag(&descriptionPtr, "description", "", "(mandatory) expense's description, example \"Groceries\"")
	regFlag(&categoryPrt, "category", "general", "(optional) expense's category, example \"Market\"")
	amountPtr = getFlag("amount")
	descriptionPtr = getFlag("description")
	categoryPrt = getFlag("category")
}

func parseFlags(args []string) {
	fmt.Println("Args:", args)
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error starting flagset: ", err.Error())
	}
}

func main() {
	//definig the command's action
	args := os.Args[1:]
	action := args[0]
	//Custom flags
	parseFlags(args[1:])
	//initialize
	initialize()
	fmt.Println("Description: ", descriptionPtr)
	fmt.Println("Amount: ", amountPtr)
	fmt.Println("Category: ", categoryPrt)

	//checking if the mandatory flags are empty
	if amountPtr == "" {
		fmt.Println("Amount not specified, plese run ./expense-tracker -h for more details.")
		os.Exit(1)
	}
	if descriptionPtr == "" {
		fmt.Println("Description not specified, plese run ./expense-tracker -h for more details.")
		os.Exit(1)
	}

	//switch between the commands
	switch action {
	case "add":
		Exp.add(descriptionPtr, amountPtr, categoryPrt)
	default:
		fmt.Printf("Command %s not found\n", action)
	}
}
