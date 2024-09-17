# GitHub User Activity

Sample solution for the [expense-tracker](https://roadmap.sh/projects/expense-tracker) challenge from [roadmap.sh](https://roadmap.sh/).

## How to run

Clone the repository and run the following command:

```bash
git clone https://github.com/fornazierr/roadmap-projects
cd /expense-tracker
go build expense-tracker.go
```

Run the following command to build and run the project:

Add command
```bash
./expense-tracker add --amount 1.99 --description "Pão de queijo"
#Expense added ID: 1
```

Update command
```bash
./expense-tracker update --id=1 --amount 2.99 --description "2 pao de queijo" --category food
#Expense updated, ID 1
```

Summary command
```bash
#ypou could use the property -year and -month. By default summary the current year.
./expense-tracker summary
#Total expenses:  2.99
```

List command
```bash
./expense-tracker list
#Total expenses:  2.99
#|ID    |EXPENSE                       |AMOUNT    |DATE      |CATEGORY       |
#|1     |2 pao de queijo               |2.9900    |2024-09-16|food           |
```

Delete command
```bash
./expense-tracker delete --id=1
#Expense deleted, ID: 1
```