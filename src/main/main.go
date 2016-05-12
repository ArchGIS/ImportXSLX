package main

import (
	"errs"
	"fmt"
	"importer"
)

// #FIXME: некоторые строки могут пустыми

func number(data string) {

}

var scheme1 = importer.ParseScheme{
	"scheme1", []importer.ParseSchemeCell{
		{"Номер", number},
		{"Название", number},
		{"Эпоха", number},
		{"Описание", number},
		{"Библиографические ссылки", number},
		{"Страницы", number},
	},
}

// Защита от неожиданных ошибок и паник времени выполнения
func panicGuard() {
	if err := recover(); err != nil {
		if fatal, ok := err.(errs.FatalError); ok {
			fmt.Printf(`{"fatal":"%s"}`, fatal.Error())
		} else {
			fmt.Printf(`{"fatal":"unexpected error: %+v"}`, err)
		}
	}
}

func printErrorsJson(e []error) {
	print(`{"errors":[`)
	for i := 0; i < len(e)-2; i++ {
		fmt.Printf(`"%s",`, e[i].Error())
	}
	fmt.Printf(`"%s"`, e[len(e)-1].Error())
	print(`]}`)
}

func main() {
	defer panicGuard()

	importer, err := importer.New("input/test.xlsx", scheme1)
	if err != nil {
		panic(err)
	}

	validationErrs := importer.ValidateHeader()
	if len(validationErrs) == 0 {
		importer.Parse()
		query := importer.CypherString()
		println(query)
	} else {
		printErrorsJson(validationErrs)
	}
}
