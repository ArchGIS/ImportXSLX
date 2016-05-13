package parsers

import (
	"bytes"
	"enum"
	"errs"
	"fmt"
	"importer"
	"regexp"
	"strings"
	"xl"
)

const MAX_ROWS = 128

func escape(s string) string {
	return strings.Replace(s, `"`, `\"`, -1)
}

var (
	// "1000г."  "1000 г." "1000 г" "1000г"
	yearRx = regexp.MustCompile(`(\d+)\s*г[,. ]`)
	// "с. 10" "с 10" "с10" "с.10" "c.10-11"
	pagesRx = regexp.MustCompile(`[^\p{L}]с\.?\s*(\d+(-\d+)?)`)
	// "#10" "№10"
	numRx = regexp.MustCompile(`[#№](\d+)`)
	// "Имя И.И." "Имя" "Имя-Имя И.И."
	nameRx = regexp.MustCompile(`[\p{L}-]+\s*\p{L}\.\p{L}\.|\p{L}+`)
	// "1.0;1.3" "1.0,1.3" "1.0 1.3"
	coordsRx = regexp.MustCompile(`(\d+\.\d+)\s*[;,]?\s*(\d+\.\d+)`)
)

var scheme1 = importer.ParseScheme{
	"scheme1", []string{
		"Номер",
		"Название",
		"Тип памятника",
		"Культурная принадлежность",
		"Эпоха",
		"Описание",
		"Библиографические ссылки",
		"Координаты",
		"Страницы",
	},
}

var mandatoryCells = []string{
	"Номер",
	"Эпоха",
	"Название",
	"Тип памятника",
	"Культурная принадлежность",
}

var validators = map[string]func(string) bool{
	"Тип памятника": enum.MonumentTypeExists,
	"Эпоха":         enum.EpochExists,
}

func NewParser1() *Parser1 {
	return &Parser1{
		scheme:   &scheme1,
		epochs:   make(map[string]string),
		cultures: make(map[string]string),
	}
}

func fetchCoords(rawCoords string) (string, string, string) {
	coords := coordsRx.FindStringSubmatch(rawCoords)
	if 3 != len(coords) {
		return "", "", "Неправильный формат координат"
	}

	return coords[1], coords[2], ""
}

func (my *Parser1) Parse(table *xl.Table) error {
	// Отбираем уникальные эпохи и культуры
	const marker = "+"
	for _, row := range table.Rows {
		if enum.EpochExists(row.Cells["Эпоха"]) {
			my.epochs[row.Cells["Эпоха"]] = marker
		}
		my.cultures[row.Cells["Культурная принадлежность"]] = marker
	}
	if len(my.epochs) == 0 {
		return fmt.Errorf("Не найдено ни одной эпохи")
	}
	if len(my.cultures) == 0 {
		return fmt.Errorf("Не заполнена ни одна культурная принадлежность")
	}

	// Присваиваем им уникальные в рамках запроса идентификаторы
	index := 0
	for epochName := range my.epochs {
		my.epochs[epochName] = fmt.Sprintf("e%d", index)
		index += 1
	}
	index = 0
	for cultureName := range my.cultures {
		my.cultures[cultureName] = fmt.Sprintf("c%d", index)
		index += 1
	}

	return nil
}

func validateRow(row xl.Row, e *errs.RowError) {
	cells := row.Cells

	for _, cellName := range mandatoryCells {
		if "" == cells[cellName] {
			e.PushError(fmt.Sprintf(`Не задано поле "%s"`, cellName))
		}
	}
	for cellName, validator := range validators {
		if !validator(cells[cellName]) {
			e.PushError(fmt.Sprintf(`Неправильно задано поле "%s"`, cellName))
		}
	}
}

func (my *Parser1) CypherString(mapId string, table *xl.Table) (string, []error) {
	if "" == mapId {
		panic(errs.NewFatal("Не задан mapId"))
	}

	var buf bytes.Buffer
	e := []error{}

	buf.WriteString(fmt.Sprintf("MATCH (map:Literature {id:%s}\n", mapId))
	buf.WriteString(fmt.Sprintf("MATCH (research:Research)-[:Contains]->(map)\n"))

	// Собираем эпохи. Они "общие" для всего запроса.
	for epoch, id := range my.epochs {
		const epochPat = `MATCH (%s:Epoch {name:"%s"})` + "\n"
		buf.WriteString(fmt.Sprintf(epochPat, id, epoch))
	}
	for culture, id := range my.cultures {
		const culturePat = `MATCH (%s:Culture {name:"%s"})` + "\n"
		buf.WriteString(fmt.Sprintf(culturePat, id, culture))
	}

	if len(table.Rows) > MAX_ROWS {
		err := fmt.Errorf("Превышен лимит строк: %d/%d", len(table.Rows), MAX_ROWS)
		return "", []error{err}
	}

	// Пишем строки данных для памятников
	for rowIndex, row := range table.Rows {
		rowErrs := errs.NewRow(row, rowIndex)
		e = append(e, rowErrs)

		cells := row.Cells
		key := fmt.Sprintf("monument%d", rowIndex)

		// Сначала предварительная валидация.
		validateRow(row, rowErrs)
		if len(rowErrs.Texts) > 0 {
			continue
		}

		pages := cells["Страницы"]
		n := cells["Номер"]
		epoch := cells["Эпоха"]
		culture := cells["Культурная принадлежность"]
		name := cells["Название"]
		ty := cells["Тип памятника"]

		buf.WriteString(fmt.Sprintf("// Data for %s\n", key))

		// Создание самого памятника.
		if "" == cells["Координаты"] {
			buf.WriteString(
				fmt.Sprintf("CREATE (%s {typeId:%d})\n", key, enum.MonumentTypeId(ty)),
			)
		} else {
			x, y, err := fetchCoords(cells["Координаты"])
			if "" == err {
				buf.WriteString(fmt.Sprintf(
					"CREATE (%s {typeId:%d,x:%s,y:%s})\n",
					key,
					enum.MonumentTypeId(ty),
					x, y,
				))
			} else {
				rowErrs.PushError(err)
			}
		}

		// Ребро от карты к памятнику
		if "" == pages {
			const pat = `CREATE (map)-[:References {n:"%s"}]->(%s)` + "\n"
			buf.WriteString(fmt.Sprintf(pat, n, key))
		} else {
			const pat = `CREATE (map)-[:References {n:"%s", pages:"%s"}]->(%s)` + "\n"
			buf.WriteString(fmt.Sprintf(pat, n, pages, key))
		}

		// Эпоха памятника
		buf.WriteString(
			fmt.Sprintf("CREATE (%s)-[:Has]->(%s)\n", my.epochs[epoch], key),
		)

		// Knowledge
		knowledgeKey := key + "_k"
		description := cells["Описание"]
		if "" != description {
			const pat = `CREATE (%s:Knowledge {name:"%s",description:"%s"})-[:Describes]->(%s)` + "\n"
			buf.WriteString(fmt.Sprintf(pat, knowledgeKey, name, escape(description), key))
		} else {
			const pat = `CREATE (%s:Knowledge {name:"%s"})-[:Describes]->(%s)` + "\n"
			buf.WriteString(fmt.Sprintf(pat, knowledgeKey, name, key))
		}
		const cultureEdgePat = `CREATE (%s)-[:Mentions]->(%s)` + "\n"
		buf.WriteString(fmt.Sprintf(cultureEdgePat, knowledgeKey, my.cultures[culture]))
		buf.WriteString(fmt.Sprintf(
			"CREATE (research)-[:Contains]->(%s)\n", knowledgeKey,
		))

		// Библиографические ссылки. Они должны быть разделены через ";"
		litRefs := strings.Split(cells["Библиографические ссылки"], ";")
		if 0 != len(litRefs) {
			for litRefIndex, ref := range litRefs {
				name := nameRx.FindString(ref)
				if "" == name {
					rowErrs.PushError("Библиографическая ссылка: не найдено имя автора")
					continue // Имя обязательно
				}
				year := yearRx.FindStringSubmatch(ref)
				if 0 == len(year) {
					rowErrs.PushError("Библиографическая ссылка: не найден год исследования")
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

	return buf.String(), e
}

func (my *Parser1) Scheme() *importer.ParseScheme {
	return my.scheme
}
