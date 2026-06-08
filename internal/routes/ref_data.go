package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kofalt/go-memoize"
	"github.com/map-services/street-manager-relay/internal"
	"github.com/map-services/street-manager-relay/models"
)

func HandleRefData(repo *internal.DbRepository, cache *memoize.Memoizer) gin.HandlerFunc {
	return func(c *gin.Context) {

		refData, err, _ := memoize.Call(cache, "refdata", func() (*models.RefData, error) {
			return repo.RefData()
		})
		if err != nil {
			abortWithError(c, http.StatusInternalServerError, "Failed to fetch reference data", err)
			return
		}

		c.JSON(200, gin.H{
			"refdata":     refData,
			"attribution": internal.ATTRIBUTION,
		})
	}
}
