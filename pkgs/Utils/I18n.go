package Utils

import (
	"encoding/json"
	"github.com/kevinu2/ngo2/pkgs/BigMath"
)

type I18n struct {
}

func (I18n) String (universal, i18n, lang string) string {
	i18nMap := make(map[string]string)
	err := json.Unmarshal([]byte(i18n), &i18nMap)
	if err != nil {
		return universal
	}
	if _, ok := i18nMap[lang]; ok {
		return i18nMap[lang]
	}
	return universal
}

func (I18n) ByPN (numberStr, universalPositive, i18nPositive, universalNegative, i18nNegative, lang string) string {
	if BigMath.Gt(numberStr, "0") {
		return I18nString(universalPositive, i18nPositive, lang)
	} else if BigMath.Lt(numberStr, "0") {
		return I18nString(universalNegative, i18nNegative, lang)
	} else {
		return ""
	}
}
