package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
var amountFlag, descriptionFlag, categoryFlag, monthFlag, yearFlag, idFlag string
var fs *flag.FlagSet
var commandMessage string
var command string
var subCom string

/*
Check if exists our data.json file, if wasn't, create a new one to store the expense's data
*/
func dataLoader() error {
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

func parseFlags(command string, args []string) {
	fs = flag.NewFlagSet("flags", flag.ContinueOnError)

	initFlags(command)
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error starting flagset: ", err.Error())
	}
}

func initFlags(command string) {
	//defining expense program flags
	t := time.Now()
	if command == "add" || command == "update" {
		regFlag(&amountFlag, "amount", "", "(mandatory) expense's amount, example 32.99")
		regFlag(&descriptionFlag, "description", "", "(mandatory) expense's description, example \"Groceries\"")
		regFlag(&categoryFlag, "category", "general", "(optional) expense's category, example \"Market\"")
		regFlag(&idFlag, "id", "", "(mandatory) expense's ID, must be grater than 0. Use list command to gather this value.")

		amountFlag = getFlag("amount")
		descriptionFlag = getFlag("description")
		categoryFlag = getFlag("category")
		idFlag = getFlag("id")
	}

	if command == "summary" {
		regFlag(&monthFlag, "month", "", "(optional) month's summary filter, must be between 1 and 12, example 8")
		regFlag(&yearFlag, "year", strconv.Itoa(t.Year()), "(optional) year's summary filter, example 2012")

		monthFlag = getFlag("month")
		yearFlag = getFlag("year")
	}

	if command == "delete" {
		regFlag(&idFlag, "id", "", "(mandatory) expense's ID, must be grater than 0. Use list command to gather this value.")
		idFlag = getFlag("id")
	}
}

/*
Registra a FLAG
*/
func regFlag(v *string, name, value, usage string) {
	if fs.Lookup(name) == nil {
		fs.StringVar(v, name, value, usage)
	}
}

/*
Return the flag value
*/
func getFlag(name string) string {
	return fs.Lookup(name).Value.String()
}

/*
Checking if the mandatory flags are empty
*/
func checkFlags(command string) {
	if command == "add" || command == "update" {
		if amountFlag == "" {
			fmt.Println("amount ", amountFlag)
			fmt.Println(commandMessage)
			os.Exit(1)
		}
		if descriptionFlag == "" {
			fmt.Println("description ", descriptionFlag)
			fmt.Println(commandMessage)
			os.Exit(1)
		}
	}
	if command == "update" || command == "delete" {
		if idFlag == "" {
			fmt.Println("id ", idFlag)
			fmt.Println(commandMessage)
			os.Exit(1)
		}
	}
}

/*
Verify if a flag is a help flag
*/
func isHelp() bool {
	switch subCom {
	case "-h":
		return true
	case "--h":
		return true
	case "-help":
		return true
	case "--help":
		return true
	}
	return false
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
func (e Expenses) add(args []string) error {
	parseFlags("add", args)
	if isHelp() {
		fs.PrintDefaults()
		os.Exit(1)
	}
	checkFlags("add")
	var idExpense int
	if len(dataKeys) == 0 {
		idExpense = 1
	} else {
		idExpense = dataKeys[len(dataKeys)-1] + 1
	}
	dataKeys = append(dataKeys, idExpense)
	t := time.Now()

	f, err := strconv.ParseFloat(amountFlag, 64)
	if err != nil {
		return err
	}
	expense := Expense{
		Id:          idExpense,
		Description: descriptionFlag,
		Amount:      f,
		Category:    categoryFlag,
		Date:        fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()),
	}
	fmt.Println("idExpense: ", idExpense)
	e[strconv.Itoa(idExpense)] = expense
	e.writeFile()
	fmt.Println("ID:", idExpense)
	return nil
}

func (e Expenses) update(args []string) {
	parseFlags("update", args)
	if isHelp() {
		fs.PrintDefaults()
		os.Exit(1)
	}
	checkFlags("update")
	var err error
	expense, bol := e.exists(idFlag)
	if !bol {
		fmt.Printf("Expense not found, ID: %s", idFlag)
		os.Exit(1)
	}
	expense.Amount, err = strconv.ParseFloat(amountFlag, 64)
	if err != nil {
		fmt.Println("Error converting amount, please use a correct value. Ex.: 2.99", err.Error())
		os.Exit(1)
	}
	expense.Description = descriptionFlag
	//prevent the default category override custom category
	if categoryFlag != "general" {
		expense.Category = categoryFlag
	}
	e[idFlag] = expense
	e.writeFile()
	fmt.Printf("Expense updated, ID %s", idFlag)
}

/*
Delete a expense bby your ID
*/
func (e Expenses) delete(args []string) {
	parseFlags("delete", args)
	if isHelp() {
		fs.PrintDefaults()
		os.Exit(1)
	}
	checkFlags("delete")
	_, exists := e.exists(idFlag)
	if !exists {
		fmt.Printf("Expense not found, ID: %s", idFlag)
		os.Exit(1)
	}
	delete(e, idFlag)
	e.writeFile()
	fmt.Printf("Expense deleted. ID: %s\n", idFlag)
}

/*
Run some validation on month argument
*/
func validateMonth(monthArg, monthExpense string) bool {
	a, err := strconv.Atoi(monthArg)
	if err != nil {
		fmt.Printf("Fail filtering by month, the vale must be between 1 and 12, please check the value %s\n", monthArg)
		os.Exit(1)
	}
	e, err := strconv.Atoi(monthExpense)
	if err != nil {
		fmt.Println("Fail filtering by month, please check the expensive date.")
		os.Exit(1)
	}
	if a < 1 || a > 12 {
		fmt.Printf("The month must be between 1 and 12. Invalid value %s", monthArg)
		os.Exit(1)
	}
	if a != e {
		return false
	}
	return true
}

/*
Filter expenses by month
*/
func (e *Expenses) filterByMonth() Expenses {
	ret := make(Expenses, 0)
	for k, v := range *e {
		m := strings.Split(v.Date, "-")
		if validateMonth(monthFlag, m[1]) {
			ret[k] = v
		}
	}
	return ret
}

/*
Filter expenses by year, default current year.
*/
func (e *Expenses) filterByYear() Expenses {
	ret := make(Expenses, 0)
	for k, v := range *e {
		y := strings.Split(v.Date, "-")
		if yearFlag == y[0] {
			ret[k] = v
		}
	}
	return ret
}

/*
Print the summary
*/
func (e Expenses) printSummary(args []string) {
	parseFlags("summary", args)
	if isHelp() {
		fs.PrintDefaults()
		os.Exit(1)
	}
	var sum float64
	var aux Expenses
	if yearFlag != "" {
		aux = e.filterByYear()
	}
	if monthFlag != "" {
		aux = aux.filterByMonth()
	}
	for _, v := range aux {
		sum += v.Amount
	}
	fmt.Println("Total expenses: ", sum)
}

func (e Expenses) list() {
	ptrn := "|%-6s|%-30s|%-10s|%-10s|%-15s|\n"
	fmt.Printf(ptrn, "ID", "EXPENSE", "AMOUNT", "DATE", "CATEGORY")
	for k, v := range e {
		amount := strconv.FormatFloat(v.Amount, 'f', 4, 32)
		fmt.Printf(ptrn, k, v.Description, amount, v.Date, v.Category)
	}
}

func main() {
	commandMessage = "Parameter not specified, plese run ./expense-tracker command -h for more details."
	//definig the command's action
	args := os.Args[1:]
	command = args[0]
	if len(args) > 1 {
		subCom = args[1]
	}

	//start map
	Exp = make(Expenses)
	//Load data from data.json file
	dataLoader()

	//switch between the commands
	switch command {
	case "add":
		//Custom flags
		Exp.add(args[1:])
		os.Exit(1)
	case "update":
		//Custom flags
		Exp.update(args[1:])
		os.Exit(1)
	case "delete":
		Exp.delete(args[1:])
		os.Exit(1)
	case "summary":
		Exp.printSummary(args[1:])
		os.Exit(1)
	case "list":
		Exp.list()
		os.Exit(1)
	default:
		fmt.Println("hi")
	}
}
