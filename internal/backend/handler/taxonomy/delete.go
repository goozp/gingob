package taxonomy

import (
	"strconv"

	Response "github.com/puti-projects/puti/internal/backend/handler"
	"github.com/puti-projects/puti/internal/backend/service"
	"github.com/puti-projects/puti/internal/pkg/errno"

	"github.com/gin-gonic/gin"
)

// Delete delete the taxonomy directly without soft delete
func Delete(c *gin.Context) {
	ID, _ := strconv.Atoi(c.Param("id"))
	taxonomyType := c.Query("taxonomy") // TODO

	termID := uint64(ID)

	// check
	if err := checkIfCanDelete(termID, taxonomyType); err != nil {
		Response.SendResponse(c, err, nil)
		return
	}

	if err := service.DeleteTaxonomy(termID, taxonomyType); err != nil {
		Response.SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	Response.SendResponse(c, nil, nil)
}

func checkIfCanDelete(termID uint64, taxonomyType string) error {
	if taxonomyType != "category" && taxonomyType != "tag" {
		return errno.New(errno.ErrValidation, nil).Add("error taxonomy.")
	}

	if ifHasChild := service.IfTaxonomyHasChild(termID, taxonomyType); ifHasChild == true {
		return errno.New(errno.ErrValidation, nil).Add("taxonomy has children and can not be deleted")
	}

	return nil
}
