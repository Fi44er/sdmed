package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",

	'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "yo",
	'Ж': "zh", 'З': "z", 'И': "i", 'Й': "y", 'К': "k", 'Л': "l", 'М': "m",
	'Н': "n", 'О': "o", 'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u",
	'Ф': "f", 'Х': "h", 'Ц': "ts", 'Ч': "ch", 'Ш': "sh", 'Щ': "sch",
	'Ъ': "", 'Ы': "y", 'Ь': "", 'Э': "e", 'Ю': "yu", 'Я': "ya",
}

func CreateSlugRU(s string) string {
	var result strings.Builder

	for _, r := range s {
		if unicode.Is(unicode.Cyrillic, r) {
			if eng, ok := translitMap[r]; ok {
				result.WriteString(eng)
			}
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	slug := result.String()

	reMultiDash := regexp.MustCompile(`-+`)
	slug = reMultiDash.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
