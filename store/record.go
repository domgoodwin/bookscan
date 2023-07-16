package store

import (
	"github.com/domgoodwin/bookscan/items"
	"github.com/sirupsen/logrus"
)

const recordCsvFilePath = "./records.csv"

type RecordStorer struct {
	// map of database ids to map of book isbns to books
	records    map[string]map[string]*items.Record
	databaseID string
}

func (s *RecordStorer) Setup(in map[string]map[string]*items.Record) {
	s.records = in
}

func (s *RecordStorer) StoreItem(databaseID string, r *items.Record) bool {
	_, databaseFound := s.records[databaseID]
	if !databaseFound {
		s.records[databaseID] = make(map[string]*items.Record)
	}
	_, recordFound := s.records[databaseID][r.Barcode]
	if !recordFound {
		s.records[databaseID][r.Barcode] = r
	}
	return recordFound
}

func (s *RecordStorer) CheckIfItemInCache(databaseID, barcode string) (*items.Record, bool) {
	_, found := s.records[databaseID]
	if !found {
		s.records[databaseID] = make(map[string]*items.Record)
	}
	record, found := s.records[databaseID][barcode]
	return record, found
}

func (s *RecordStorer) ClearCache(databaseID string) int {
	logrus.Info("clearing cache")
	oldLength := len(s.records[databaseID])
	s.records = make(map[string]map[string]*items.Record)
	return oldLength
}

func (s *RecordStorer) Length() int {
	return len(s.records)
}

func (s *RecordStorer) DatabaseID() string {
	return s.databaseID
}
