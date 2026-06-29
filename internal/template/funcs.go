package template

import (
	"strings"
	"text/template"
	"unicode"
)

// FuncMap returns the helpers exposed inside every goforge template.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"trim":       strings.TrimSpace,
		"replace":    func(old, new, s string) string { return strings.ReplaceAll(s, old, new) },
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"contains":   strings.Contains,
		"split":      strings.Split,
		"join":       func(sep string, parts []string) string { return strings.Join(parts, sep) },
		"quote":      func(s string) string { return `"` + s + `"` },
		"pascal":     pascalCase,
		"camel":      camelCase,
		"kebab":      kebabCase,
		"snake":      snakeCase,
		"default":    defaultIfEmpty,
		"trimPrefix": strings.TrimPrefix,
		"trimSuffix": strings.TrimSuffix,
	}
}

func defaultIfEmpty(fallback, v string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func splitWords(s string) []string {
	var words []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			words = append(words, cur.String())
			cur.Reset()
		}
	}
	for i, r := range s {
		switch {
		case r == '-' || r == '_' || r == ' ' || r == '.' || r == '/':
			flush()
		case unicode.IsUpper(r) && i > 0:
			prev := []rune(s)[i-1]
			if unicode.IsLower(prev) || unicode.IsDigit(prev) {
				flush()
			}
			cur.WriteRune(r)
		default:
			cur.WriteRune(r)
		}
	}
	flush()
	return words
}

func pascalCase(s string) string {
	var b strings.Builder
	for _, w := range splitWords(s) {
		if w == "" {
			continue
		}
		runes := []rune(strings.ToLower(w))
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}

func camelCase(s string) string {
	p := pascalCase(s)
	if p == "" {
		return p
	}
	runes := []rune(p)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func kebabCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "-")
}

func snakeCase(s string) string {
	parts := splitWords(s)
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "_")
}
