package utils

import "strings"

func ConvertTitleToSlug(title string) string {
	return strings.Join(strings.Split(title, " "), "-")
}

func CovertSlugToTitle(title string) string {
	return strings.Join(strings.Split(title, "-"), " ")
}
