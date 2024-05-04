package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type EmptyFieldErr struct {
	field string
}

func (e EmptyFieldErr) Error() string {
	return fmt.Sprintf("missing field: %v", e.field)
}

var ErrInvalidId = errors.New("invalid id")

func main() {
	d := json.NewDecoder(strings.NewReader(data))
	count := 0
	for d.More() {
		count++
		var emp Employee
		err := d.Decode(&emp)
		if err != nil {
			fmt.Printf("record %d: %v\n", count, err)
			continue
		}
		errs := ValidateEmployee(emp)
		fmt.Printf("record %d: %+v\n", count, emp)
		for _, err := range errs {
			if ok := errors.Is(err, ErrInvalidId); ok {
				fmt.Printf("invalid id: %+v \n", err)
				continue
			}

			var emptyFieldErr EmptyFieldErr
			if ok := errors.As(err, &emptyFieldErr); ok {
				fmt.Printf("%+v \n", emptyFieldErr)
				continue
			}

		}
	}
}

const data = `
{
	"id": "ABCD-123",
	"first_name": "Bob",
	"last_name": "Bobson",
	"title": "Senior Manager"
}
{
	"id": "XYZ-123",
	"first_name": "Mary",
	"last_name": "Maryson",
	"title": "Vice President"
}
{
	"id": "BOTX-263",
	"first_name": "",
	"last_name": "Garciason",
	"title": "Manager"
}
{
	"id": "HLXO-829",
	"first_name": "Pierre",
	"last_name": "",
	"title": "Intern"
}
{
	"id": "MOXW-821",
	"first_name": "Franklin",
	"last_name": "Watanabe",
	"title": ""
}
{
	"id": "",
	"first_name": "Shelly",
	"last_name": "Shellson",
	"title": "CEO"
}
{
	"id": "YDOD-324",
	"first_name": "",
	"last_name": "",
	"title": ""
}
`

type Employee struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Title     string `json:"title"`
}

var (
	validID = regexp.MustCompile(`\w{4}-\d{3}`)
)

func ValidateEmployee(e Employee) (errs []error) {
	if len(e.ID) == 0 {
		errs = append(errs, EmptyFieldErr{
			field: "ID",
		})
	}
	if !validID.MatchString(e.ID) {
		errs = append(errs, ErrInvalidId)
	}
	if len(e.FirstName) == 0 {
		errs = append(errs,
			EmptyFieldErr{
				field: "FirstName",
			},
		)
	}
	if len(e.LastName) == 0 {
		errs = append(errs,
			EmptyFieldErr{
				field: "LastName",
			},
		)
	}
	if len(e.Title) == 0 {
		errs = append(errs,
			EmptyFieldErr{
				field: "Title",
			},
		)
	}
	return errs
}
