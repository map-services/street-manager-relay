package cmd

import (
	"encoding/csv"
	"log/slog"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/map-services/street-manager-relay/internal/favicon"
	"github.com/map-services/street-manager-relay/internal/promoter"
	"github.com/map-services/street-manager-relay/models"
)

func UpdateFaviconsInCSV(csvFile string) error {

	orgs, err := promoter.GetPromoterOrgsList()
	if err != nil {
		return err
	}

	updated := make([]*models.PromoterOrg, 0, len(orgs))
	for idx, record := range orgs {

		slog.Info("Processing record", "index", idx, "url", record.Url)

		iconInfo, err := favicon.Extract(record.Url)
		if err != nil {
			slog.Warn("failed to extract favicon", "url", record.Url, "error", err)
		} else {
			record.Favicon = &iconInfo.Href
		}
		updated = append(updated, record)
	}

	f, err := os.OpenFile(csvFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", csvFile)
	}
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("error closing file", "error", err)
		}
	}()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	for _, record := range updated {
		row := record.ToCSV()
		if err := csvWriter.Write(row); err != nil {
			return errors.Wrapf(err, "failed to write row=%v", row)
		}
	}
	return nil
}
