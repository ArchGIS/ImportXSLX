package enum

var epoch = map[string]struct{}{
	"Палеолит":                    struct{}{},
	"Мезолит":                     struct{}{},
	"Неолит":                      struct{}{},
	"Энеолит":                     struct{}{},
	"Бронзовый век":               struct{}{},
	"Ранний железный век":         struct{}{},
	"Великое переселение народов": struct{}{},
	"Средневековье":               struct{}{},
	"Новое время":                 struct{}{},
}

var monumentTypes = map[string]struct{}{
	"Не определено":          struct{}{},
	"Городище":               struct{}{},
	"Селище":                 struct{}{},
	"Местонахождение":        struct{}{},
	"Погребение":             struct{}{},
	"Могильник":              struct{}{},
	"Архитектурный памятник": struct{}{},
	"Грунтовый могильник":    struct{}{},
	"Курганный могильник":    struct{}{},
	"Курган":                 struct{}{},
	"Клад":                   struct{}{},
	"Мастерская":             struct{}{},
	"Стоянка":                struct{}{},
}

var monumentTypeIds = map[string]int{
	"Не определено":          0,
	"Городище":               1,
	"Селище":                 2,
	"Местонахождение":        3,
	"Погребение":             4,
	"Могильник":              5,
	"Архитектурный памятник": 6,
	"Грунтовый могильник":    7,
	"Курганный могильник":    8,
	"Курган":                 9,
	"Клад":                   10,
	"Мастерская":             11,
	"Стоянка":                12,
}

func MonumentTypeId(name string) int {
	return monumentTypeIds[name]
}

func EpochExists(name string) bool {
	_, ok := epoch[name]
	return ok
}

func MonumentTypeExists(name string) bool {
	_, ok := monumentTypes[name]
	return ok
}
