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
			fmt.Printf(`{"fatal": "%s"}`, fatal.Error())
		} else {
			fmt.Printf(`{"fatal": "unexpected error: %+v"`, err)
		}
	}
}

func main() {

	defer panicGuard()

	importer, err := importer.New("input/test.xlsx", scheme1)
	if err != nil {
		panic(err)
	}

	errs := importer.ValidateHeader()
	if len(errs) != 0 {
		println("ERRORS:")
		for _, err := range errs {
			fmt.Printf("%+v\n", err)
		}
	}

	importer.Parse()
	query := importer.CypherString()
	println(query)
}
