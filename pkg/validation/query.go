package validation

import (
	"encoding/json"

	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
)

func ValidateQuery(queries map[string]string) error {
	if len(queries) == 0 {
		return nil
	}

	errMap := []map[string]interface{}{}

	for key, value := range queries {
		if value == "" {
			errMap = append(errMap, map[string]interface{}{
				"field":   key,
				"message": "query " + key + " is cannot be empty",
			})

		}
	}

	if len(errMap) == 0 {
		return nil
	}

	jsonErr, err := json.Marshal(errMap)
	if err != nil {
		return err
	}

	return &constant.ErrValidation{Message: string(jsonErr)}
}
