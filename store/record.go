package store

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/domgoodwin/bookscan/items"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/sirupsen/logrus"
)

const recordCsvFilePath = "./records.csv"

type RecordStorer struct {
	records    map[string]*items.Record
	databaseID string
}

func (s *RecordStorer) Setup() {
	s.records = make(map[string]*items.Record)
	s.LoadRecordsFromCSV()
	s.LoadRecordsFromNotion(context.Background(), "")
}

func (s *RecordStorer) StoreItem(r *items.Record) bool {
	_, found := s.records[r.Barcode]
	if !found {
		s.records[r.Barcode] = r
	}
	return found
}

func (s *RecordStorer) CheckIfItemInCache(barcode string) (*items.Record, bool) {
	record, found := s.records[barcode]
	return record, found
}

func (s *RecordStorer) ClearCache() int {
	logrus.Info("clearing cache")
	oldLength := len(s.records)
	s.records = make(map[string]*items.Record)
	return oldLength
}

func (s *RecordStorer) Length() int {
	return len(s.records)
}

func (s *RecordStorer) DatabaseID() string {
	return s.databaseID
}

func (s *RecordStorer) LoadRecordsFromCSV() {
	if os.Getenv("CSV_CACHE") == "false" {
		logrus.Info("CSV cache disabled")
		return
	}
	count := 0
	file, err := os.Open(csvFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
		logrus.Debug(scanner.Text())
		record := s.CSVRecordToRecord(scanner.Text())
		s.StoreItem(record)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Loaded %v records from CSV", count)
}

func (s *RecordStorer) LoadRecordsFromNotion(ctx context.Context, databaseID string) error {
	if os.Getenv("NOTION_CACHE") == "false" {
		logrus.Info("Notion cache disabled")
		return nil
	}
	records, databaseID, err := notion.GetRecordPagesFromDatabase(ctx, databaseID)
	if err != nil {
		return err
	}
	for _, r := range records {
		s.StoreItem(r)
	}
	logrus.Infof("Loaded %v records from Notion", len(records))
	s.databaseID = databaseID
	return nil
}

func (s *RecordStorer) CSVRecordToRecord(line string) *items.Record {
	parts := strings.Split(line, ",")
	if len(parts) != 6 {
		logrus.Warn("Line doesn't have 6 parts", map[string]string{"line": line})
		return nil
	}

	year, err := strconv.Atoi(removeQuotes(parts[3]))
	if err != nil {
		logrus.Error(err)
	}
	return &items.Record{
		Title:    removeQuotes(parts[0]),
		Artists:  strings.Split(removeQuotes(parts[1]), ";"),
		Barcode:  removeQuotes(parts[2]),
		Year:     year,
		Link:     removeQuotes(parts[4]),
		CoverURL: removeQuotes(parts[5]),
	}
}
