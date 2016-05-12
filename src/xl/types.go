package xl

type Header []string

type Row struct {
	Header *Header
	Cells  []string
}

type Table struct {
	Header Header
	rows   []Row
}
