package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"test-project/internal/util"
)

func TestSortCollateStringGerman(t *testing.T) {
	slice := []string{"a", "ä", "e", "ö", "u", "ü", "o"}
	util.SortCollateStringSlice(slice, language.German)

	expected := []string{"a", "ä", "e", "o", "ö", "u", "ü"}
	assert.Equal(t, expected, slice)
}

func TestSortCollateStringEnglish(t *testing.T) {
	slice := []string{"a", "ä", "e", "ö", "u", "ü", "o"}
	util.SortCollateStringSlice(slice, language.English)

	expected := []string{"a", "ä", "e", "o", "ö", "u", "ü"}
	assert.Equal(t, expected, slice)
}

func TestSortCollateStringGermanAndOptions(t *testing.T) {
	slice := []string{"a", "ä", "e", "ö", "u", "ü", "o"}
	util.SortCollateStringSlice(slice, language.German, collate.IgnoreCase, collate.IgnoreWidth, collate.IgnoreDiacritics)

	expected := []string{"a", "ä", "e", "ö", "o", "u", "ü"}
	assert.Equal(t, expected, slice)
}
