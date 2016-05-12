package main

import (
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

func main() {
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
