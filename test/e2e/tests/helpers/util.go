package helpers

import "net/url"

func GetParsedURL(uri string) (*url.URL, error) {
	return url.Parse(uri)
}
