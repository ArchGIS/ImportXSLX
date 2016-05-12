package importer

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/tealeg/xlsx"
)

var (
	// "1000г."  "1000 г." "1000 г" "1000г"
	yearRx = regexp.MustCompile(`(\d+)\s*г[,. ]`)
	// "с. 10" "с 10" "с10" "с.10" "c.10-11"
	pagesRx = regexp.MustCompile(`[^\p{L}]с\.?\s*(\d+(-\d+)?)`)
	// "#10" "№10"
	numRx = regexp.MustCompile(`[#№](\d+)`)
	// "Имя И.И." "Имя" "Имя-Имя И.И."
	nameRx = regexp.MustCompile(`[\p{L}-]+\s*\p{L}\.\p{L}\.|\p{L}+`)
)

func fetch(cell *xlsx.Cell) string {
	value, err := cell.String()
	if err != nil {
		// Не знаю, что здесь за ошибка может быть, лень
		// смотреть исходники функции String(). Возможно она всегда nil
		// ошибку возвращает. Это можно проверить как-нибудь потом.
		panic(err)
	}

	return value
}

func escape(s string) string {
	return strings.Replace(s, `"`, `\"`, -1)
}

func New(xslxFilePath string, scheme ParseScheme) (*Importer, error) {
	file, err := xlsx.OpenFile(xslxFilePath)
	if err != nil {
		return nil, err
	}

	sheet := file.Sheets[0] // Работаем только с первым листом
	header := sheet.Rows[0] // Первый ряд - заголок
	rows := sheet.Rows[1:]  // Остальные - данные

	return &Importer{
		scheme:  scheme,
		header:  header,
		rows:    rows,
		epochs:  make(map[string]string),
		indexes: make(map[string]int),
	}, nil
}

func (my *Importer) ValidateHeader() []error {
	errs := []error{}

	for _, cell := range my.header.Cells {
		cellName, err := cell.String()
		if err != nil {
			panic(err)
		}

		info := my.scheme.Find(cellName)
		if "" == info.Name {
			errs = append(
				errs, fmt.Errorf("Не найдена информация по полю %s", cellName),
			)
		}
	}

	return errs
}

func (my *Importer) Parse() {
	// Запись индексов.
	for index, cell := range my.scheme.Cells {
		my.indexes[cell.Name] = index
	}

	const marker = "+"
	// Отбираем уникальные эпохи
	for _, row := range my.rows {
		cellName, _ := row.Cells[my.indexes["Эпоха"]].String()
		my.epochs[cellName] = marker
	}

	// Присваиваем им уникальные в рамках запроса идентификаторы
	index := 0
	for epoch := range my.epochs {
		my.epochs[epoch] = fmt.Sprintf("e%d", index)
		index += 1
	}
}

func (my *Importer) CypherString() string {
	var buf bytes.Buffer

	// Собираем эпохи
	for epoch, id := range my.epochs {
		const epochPat = "MATCH (%s:Epoch {name:%s})\n"
		buf.WriteString(fmt.Sprintf(epochPat, id, epoch))
	}

	// Пишем строки данных для памятников
	for index, row := range my.rows {
		println(index)
		cells := row.Cells
		key := fmt.Sprintf("monument%d", index)

		pages := fetch(cells[my.indexes["Страницы"]])
		n := fetch(cells[my.indexes["Номер"]])

		buf.WriteString(fmt.Sprintf("// Data for %s\n", key))
		if "" == n {
			panic("n is required") // #FIXME: Не паниковать, а запоминать ошибку
		} else {
			if "" == pages {
				const pat = `CREATE (map)-[:References {n:"%s"}]->(%s)` + "\n"
				buf.WriteString(fmt.Sprintf(pat, n, key))
			} else {
				const pat = `CREATE (map)-[:References {n:"%s", pages:"%s"}]->(%s)` + "\n"
				buf.WriteString(fmt.Sprintf(pat, n, pages, key))
			}
		}

		buf.WriteString(
			fmt.Sprintf("CREATE (%s {})\n", key),
		)

		epoch := my.epochs[fetch(cells[my.indexes["Эпоха"]])]
		buf.WriteString(
			fmt.Sprintf("CREATE (%s)-[:Has]->(%s)\n", epoch, key),
		)

		name := fetch(cells[my.indexes["Название"]])
		if "" != name {
			description := fetch(cells[my.indexes["Описание"]])
			if "" != description {
				const pat = `CREATE (:Knowledge {name:"%s",description:"%s"})-[:Describes]->(%s)` + "\n"
				buf.WriteString(fmt.Sprintf(pat, name, escape(description), key))
			} else {
				const pat = `CREATE (:Knowledge {name:"%s"})-[:Describes]->(%s)` + "\n"
				buf.WriteString(fmt.Sprintf(pat, name, key))
			}
		}

		litRefs := strings.Split(fetch(cells[my.indexes["Библиографические ссылки"]]), ";")
		if 0 != len(litRefs) {
			for litRefIndex, ref := range litRefs {
				name := nameRx.FindString(ref)
				if "" == name {
					continue // Имя обязательно
				}
				year := yearRx.FindStringSubmatch(ref)
				if 0 == len(year) {
					continue // Год обязателен
				}

				authorKey := fmt.Sprintf("%s_a%d", key, litRefIndex)
				buf.WriteString(
					fmt.Sprintf(`MERGE (%s {name:"%s"})`+"\n", authorKey, name),
				)
				researchKey := fmt.Sprintf("%s_r%d", key, litRefIndex)
				buf.WriteString(
					fmt.Sprintf(
						"MERGE (%s)-[:Created]->(%s:Research {year:%s})\n",
						authorKey,
						researchKey,
						year[1],
					),
				)
				literatureKey := fmt.Sprintf("%s_l%d", key, litRefIndex)
				buf.WriteString(
					fmt.Sprintf(
						"CREATE (%s)-[:Has]->(%s:Literature {year:%s})\n",
						researchKey,
						literatureKey,
						year[1],
					),
				)

				pages := pagesRx.FindStringSubmatch(ref)
				n := numRx.FindStringSubmatch(ref)
				if len(n) > 0 {
					if len(pages) > 0 {
						const pat = `CREATE (%s)-[:References {pages:"%s", n:"%s"}]->(%s)` + "\n"
						buf.WriteString(
							fmt.Sprintf(pat, literatureKey, pages[1], n[1], key),
						)
					} else {
						const pat = `CREATE (%s)-[:References {n:"%s"}]->(%s)` + "\n"
						buf.WriteString(
							fmt.Sprintf(pat, literatureKey, n[1], key),
						)
					}
				}
			}
		}
	}

	return buf.String()
}
