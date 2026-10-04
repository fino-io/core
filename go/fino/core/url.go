package core

const UrlTypeName = "Url"
const UrlTypeFullName = "core.Url"

func NewUrl(url string) (*Url, error) {
	return ParseUrl(url)
}
