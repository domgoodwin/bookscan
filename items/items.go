package items

const (
	ItemTypeBook   = "BOOK"
	ItemTypeRecord = "RECORD"
)

type Item interface {
	Book | Record
	Type() string
}
