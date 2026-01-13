package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	// Compile regex patterns once at package level
	specialCharsRegex    = regexp.MustCompile("[^a-z0-9-]+")
	multipleHyphensRegex = regexp.MustCompile("-+")
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = specialCharsRegex.ReplaceAllString(slug, "")

	// Remove multiple hyphens
	slug = multipleHyphensRegex.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateUniqueSlug generates a unique slug by checking database and appending suffix if needed
func GenerateUniqueSlug(db *gorm.DB, tableName, baseSlug, excludeID string) string {
	slug := baseSlug
	counter := 1

	for {
		var count int64
		query := db.Table(tableName).Where("slug = ?", slug)

		// Exclude current record if updating
		if excludeID != "" {
			query = query.Where("id != ?", excludeID)
		}

		query.Count(&count)

		if count == 0 {
			return slug
		}

		// Try with counter suffix
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++

		// Fallback to timestamp after too many attempts
		if counter > 100 {
			slug = fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())
			break
		}
	}

	return slug
}
