package lookup

import (
	"errors"

	"github.com/domgoodwin/bookscan/items"
	"github.com/sirupsen/logrus"
)

var recordDatasources = map[RecordDataSource]bool{
	discogsDataStore{}: true,
}

type RecordDataSource interface {
	Name() string
	LookupBarcode(barcode string) (*items.Record, error)
}

func LookupRecordBarcode(barcode string) (*items.Record, error) {
	logrus.Info("looking up barcode ", map[string]string{
		"barcode": barcode,
	})
	for d, enabled := range recordDatasources {
		if !enabled {
			logrus.Info("skipping data source: ", d.Name())
			continue
		}
		record, err := d.LookupBarcode(barcode)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		if record != nil {
			logrus.Infof("found record: %v", record)
			return record, nil
		}
	}
	logrus.Infof("didn't find record")
	return nil, errors.New("not found")
}
