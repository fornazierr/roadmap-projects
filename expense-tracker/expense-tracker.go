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
var amountArg, descriptionArg, categoryArg, monthArg, yearArg string
var fs *flag.FlagSet

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
	t := time.Now()
	regFlag(&amountArg, "amount", "", "(mandatory) expense's amount, example 32.99")
	regFlag(&descriptionArg, "description", "", "(mandatory) expense's description, example \"Groceries\"")
	regFlag(&categoryArg, "category", "general", "(optional) expense's category, example \"Market\"")
	regFlag(&monthArg, "month", "", "(optional) month's summary filter, must be between 1 and 12, example 8")
	regFlag(&yearArg, "year", strconv.Itoa(t.Year()), "(optional) year's summary filter, example 2012")

	amountArg = getFlag("amount")
	descriptionArg = getFlag("description")
	categoryArg = getFlag("category")
	monthArg = getFlag("month")
	yearArg = getFlag("year")
}

func parseFlags(args []string) {
	if err := fs.Parse(args); err != nil {
		fmt.Println("Error starting flagset: ", err.Error())
	}
}

func checkFlags() {
	//checking if the mandatory flags are empty
	if amountArg == "" {
		fmt.Println("Amount not specified, plese run ./expense-tracker -h for more details.")
		os.Exit(1)
	}
	if descriptionArg == "" {
		fmt.Println("Description not specified, plese run ./expense-tracker -h for more details.")
		os.Exit(1)
	}
}

func isHelp(arg string) bool {
	switch arg {
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
func (e Expenses) add(description, amount, category string) error {
	checkFlags()
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
	expense := Expense{
		Id:          idExpense,
		Description: description,
		Amount:      f,
		Category:    category,
		Date:        fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day()),
	}

	e[strconv.Itoa(idExpense)] = expense
	e.writeFile()
	fmt.Println("ID:", idExpense)
	return nil
}

func (e Expenses) update(idExpense, description, amount, category string) {
	checkFlags()
	var err error
	expense, bol := e.exists(idExpense)
	if !bol {
		fmt.Printf("Expense not found, ID: %s", idExpense)
		os.Exit(1)
	}
	expense.Amount, err = strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println("Error converting amount, please use a correct value. Ex.: 2.99", err.Error())
		os.Exit(1)
	}
	expense.Description = description
	//prevent the default category override custom category
	if category != "general" {
		expense.Category = category
	}
	e[idExpense] = expense
	e.writeFile()
	fmt.Printf("Expense updated, ID %s", idExpense)
}

func (e Expenses) delete(idExpense string) {
	_, exists := e.exists(idExpense)
	if !exists {
		fmt.Printf("Expense not found, ID: %s", idExpense)
		os.Exit(1)
	}
	delete(e, idExpense)
	e.writeFile()
	fmt.Printf("Expense deleted. ID: %s\n", idExpense)
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
		if validateMonth(monthArg, m[1]) {
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
		if yearArg == y[0] {
			ret[k] = v
		}
	}
	return ret
}

/*
Print the summary
*/
func (e Expenses) printSummary() {
	var sum float64
	var aux Expenses
	if yearArg != "" {
		aux = e.filterByYear()
	}
	if monthArg != "" {
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
	//definig the command's action
	args := os.Args[1:]
	action := args[0]
	//Load data from data.json file
	dataLoader()
	// fmt.Println("Description: ", descriptionPtr)
	// fmt.Println("Amount: ", amountPtr)
	// fmt.Println("Category: ", categoryPrt)

	//switch between the commands
	switch action {
	case "add":
		//Custom flags
		parseFlags(args[1:])
		Exp.add(descriptionArg, amountArg, categoryArg)
	case "update":
		idExpense := args[1]
		if idExpense == "" {
			fmt.Println("ID from Expense are empty. Please informa a valid ID, use the command list to find the desired Expense.")
			os.Exit(1)
		}
		//Custom flags
		parseFlags(args[2:])
		Exp.update(idExpense, descriptionArg, amountArg, categoryArg)
	case "delete":
		//no flags for this guy only the expense ID
		idExpense := args[1]
		if idExpense == "" {
			fmt.Println("ID from Expense are empty. Please informa a valid ID, use the command list to find the desired Expense.")
			os.Exit(1)
		}
		Exp.delete(idExpense)
	case "summary":
		// fmt.Println(args)
		if len(args) > 1 {
			parseFlags(args[1:])
		}
		Exp.printSummary()
		os.Exit(1)
	case "list":
		Exp.list()
		os.Exit(1)
	default:
		if !isHelp(action) {
			fmt.Printf("Command %s not found.\n", action)
		}
		parseFlags(args)
	}
}
