package validation

import (
	"encoding/json"

	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
)

type ValidateContentTypeConfig struct {
	AllowedTypes map[string]struct{}
}

type ValidateContentTypeOption func(*ValidateContentTypeConfig)

// WithValidateContentTypeAllowedTypes is a function to add allowed types to ValidateContentTypeConfig
func WithValidateContentTypeAllowedTypes(allowedTypes []string) ValidateContentTypeOption {
	return func(config *ValidateContentTypeConfig) {
		if config.AllowedTypes == nil {
			config.AllowedTypes = make(map[string]struct{})
		}

		for _, allowedType := range allowedTypes {
			config.AllowedTypes[allowedType] = struct{}{}
		}
	}
}

// WithValidateContentTypeImage is a function to add allowed image types to ValidateContentTypeConfig
func WithValidateContentTypeImage() ValidateContentTypeOption {
	return func(config *ValidateContentTypeConfig) {
		if config.AllowedTypes == nil {
			config.AllowedTypes = make(map[string]struct{})
		}

		config.AllowedTypes["image/png"] = struct{}{}
		config.AllowedTypes["image/jpg"] = struct{}{}
		config.AllowedTypes["image/jpeg"] = struct{}{}
		config.AllowedTypes["image/gif"] = struct{}{}
		config.AllowedTypes["application/octet-stream"] = struct{}{}
	}
}

var defaultValidateContentTypeConfig = ValidateContentTypeConfig{
	AllowedTypes: map[string]struct{}{},
}

func ValidateContentType(field, contentType string, opt ...ValidateContentTypeOption) error {
	config := defaultValidateContentTypeConfig

	for _, o := range opt {
		o(&config)
	}

	if len(config.AllowedTypes) == 0 {
		return nil
	}

	if _, ok := config.AllowedTypes[contentType]; ok {
		return nil
	}

	errMap := []map[string]interface{}{
		{
			"field":   field,
			"message": contentType + " type not allowed",
		},
	}

	jsonErr, err := json.Marshal(errMap)
	if err != nil {
		return err
	}

	return &constant.ErrValidation{Message: string(jsonErr)}
}

type ValidateFileSizeConfig struct {
	MaxSize int64
	MinSize int64
}

type ValidateFileSizeOption func(*ValidateFileSizeConfig)

// WithValidateFileSizeMaxSize is a function to add max size to ValidateFileSizeConfig
func WithValidateFileSizeMaxSize(maxSize int64) ValidateFileSizeOption {
	return func(config *ValidateFileSizeConfig) {
		config.MaxSize = maxSize
	}
}

// WithValidateFileSizeMinSize is a function to add min size to ValidateFileSizeConfig
func WithValidateFileSizeMinSize(minSize int64) ValidateFileSizeOption {
	return func(config *ValidateFileSizeConfig) {
		config.MinSize = minSize
	}
}

var defaultValidateFileSizeConfig = ValidateFileSizeConfig{
	// 2MB
	MaxSize: 2 * 1024 * 1024,

	MinSize: 0,
}

func ValidateFileSize(field string, fileSize int64, opt ...ValidateFileSizeOption) error {
	config := defaultValidateFileSizeConfig

	for _, o := range opt {
		o(&config)
	}

	if fileSize <= config.MaxSize {
		return nil
	}

	errMap := []map[string]interface{}{
		{
			"field":   field,
			"message": "file size too large",
		},
	}

	jsonErr, err := json.Marshal(errMap)
	if err != nil {
		return err
	}

	return &constant.ErrValidation{Message: string(jsonErr)}
}
