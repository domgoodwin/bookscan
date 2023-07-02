package store

var BookStore *BookStorer
var RecordStore *RecordStorer

func SetupStore() {
	BookStore = &BookStorer{}
	BookStore.Setup()
	RecordStore = &RecordStorer{}
	RecordStore.Setup()
}
