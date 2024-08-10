package pages

import "strconv"

func getImageAlt(key int) string {
	return strconv.Itoa(key)
}

func getImageUrl(url string) string {
	return "/" + url
}
