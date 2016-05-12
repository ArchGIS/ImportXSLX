package xl

type Header []string

type Row struct {
	// Тут не подойдёт map[string]string потому что нам важен порядок
	// элементов в ряду. Map в Go имеет рандомный порядок обхода.
	Header *Header
	// Cells  []string
	Cells map[string]string
}

type Table struct {
	Header Header
	Rows   []Row
}
