package httpapi

import "strings"

func normalizePath(path string) string {
	if strings.HasPrefix(path, "/students/") && path != "/students/" {
		return "/students/{id}"
	}
	return path
}
