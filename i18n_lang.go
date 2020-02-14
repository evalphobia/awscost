package awscost

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var defaultPrinter = message.NewPrinter(language.English)

func init() {
	for _, lang := range langs {
		for _, l := range lang.List {
			_ = message.SetString(l.Lang, lang.Key, l.Value)
		}
	}
	setDefaultLanguage()
}

func setDefaultLanguage() {
	var tag language.Tag
	switch getEnvLang() {
	case "ja":
		tag = language.Japanese
	default:
		tag = language.English
	}
	defaultPrinter = message.NewPrinter(tag)
}

func Message(key string, args ...interface{}) string {
	return defaultPrinter.Sprintf(key, args...)
}

func CommaNumber(n int) string {
	return defaultPrinter.Sprintf("%d", n)
}

var langs = []translation{
	{Key: "[AWS Estimate Costs] %s", List: []translationData{
		{language.Japanese, "[AWS概算コスト] %s"},
	}},
}

type translation struct {
	Key  string
	List []translationData
}

type translationData struct {
	Lang  language.Tag
	Value string
}
