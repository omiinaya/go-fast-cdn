package database

import (
	"errors"

	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"gorm.io/gorm"
)

// ConfigRepo provides CRUD for config key/values
func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}

type ConfigRepo struct {
	db *gorm.DB
}

func (r *ConfigRepo) Get(key string) (string, error) {
	var config models.Config
	if err := r.db.First(&config, "key = ?", key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If the config is not found, create it with a default value
			defaultValue := getDefaultConfigValue(key)
			if err := r.Set(key, defaultValue); err != nil {
				return "", err
			}
			// Try to get it again after creating it
			if err := r.db.First(&config, "key = ?", key).Error; err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return config.Value, nil
}

// getDefaultConfigValue returns the default value for a given config key
func getDefaultConfigValue(key string) string {
	switch key {
	case "registration_enabled":
		return "true"
	default:
		return ""
	}
}

func (r *ConfigRepo) Set(key, value string) error {
	// Use a transaction to handle concurrent creation attempts
	return r.db.Transaction(func(tx *gorm.DB) error {
		var config models.Config
		err := tx.First(&config, "key = ?", key).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new config
			config.Key = key
			config.Value = value
			if err := tx.Create(&config).Error; err != nil {
				// Check if it's a unique constraint violation (record created by another goroutine)
				if isDuplicateKeyError(err) {
					// If another goroutine created it, try to get it again
					return tx.First(&config, "key = ?", key).Error
				}
				return err
			}
		} else if err != nil {
			return err
		} else {
			// Update existing config
			config.Value = value
			if err := tx.Save(&config).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// isDuplicateKeyError checks if the error is a duplicate key error
func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "UNIQUE constraint failed: configs.key" ||
		err.Error() == "duplicate key value violates unique constraint" ||
		err.Error() == "UNIQUE constraint violated")
}

// InitializeDefaultConfigs creates default configuration values if they don't exist
func InitializeDefaultConfigs() error {
	configRepo := NewConfigRepo(DB)

	// Check if "registration_enabled" config exists
	_, err := configRepo.Get("registration_enabled")
	if err != nil {
		// If it doesn't exist, create it with a default value of "true"
		err := configRepo.Set("registration_enabled", "true")
		if err != nil {
			// Log the error and return it to stop the application
			// In a production app, you might want to use a proper logger
			println("Failed to initialize default config 'registration_enabled':", err.Error())
			return err
		}
	}
	return nil
}

// EnsureDefaultConfigExists checks if a specific config exists and creates it with a default value if it doesn't
func EnsureDefaultConfigExists(key, defaultValue string) error {
	configRepo := NewConfigRepo(DB)

	// Check if the config exists
	_, err := configRepo.Get(key)
	if err != nil {
		// If it doesn't exist, create it with the default value
		err := configRepo.Set(key, defaultValue)
		if err != nil {
			println("Failed to initialize default config '"+key+"':", err.Error())
			return err
		}
	}
	return nil
}
