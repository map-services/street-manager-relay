package cmd

import (
	"log/slog"

	"github.com/cockroachdb/errors"
	"github.com/map-services/street-manager-relay/internal"
)

func RegenerateIndex(dbPath string) error {
	repo, err := internal.NewDbRepository(dbPath)
	if err != nil {
		return errors.Wrap(err, "failed to initialize db repository")
	}
	defer func() {
		if err := repo.Close(); err != nil {
			slog.Error("Error closing database", "error", err)
		}
	}()

	affected, total, err := repo.RegenerateIndex()
	if err != nil {
		return errors.Wrap(err, "error regenerating index")
	}

	if total > 0 {
		slog.Info("Affected records", "affected", affected, "total", total, "percentage", float64(affected)/float64(total)*100.0)
	} else {
		slog.Info("No records found to process.")
	}

	return nil
}
