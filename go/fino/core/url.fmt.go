package core

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var urlSchemaPattern = regexp.MustCompile(`^([a-z0-9+\-.]+://)|mailto:|news:`)

func ParseUrl(rawUrl string) (*Url, error) {
	u := &Url{}
	if err := u.Parse(rawUrl); err != nil {
		return nil, err
	}
	return u, nil
}

func (x *Url) Parse(rawUrl string) error {
	if x == nil {
		return errors.New("url is nil")
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	x.Scheme = u.Scheme
	x.Authority = &Url_Authority{
		UserInfo: u.User.String(),
		Host:     u.Hostname(),
		Port:     u.Port(),
	}

	if x.Authority.Host == "" && (u.Scheme == "" || u.Opaque == "") && urlSchemaPattern.MatchString(rawUrl) {
		u, err := url.Parse("https://" + rawUrl)
		if err != nil {
			return err
		}
		x.Scheme = ""
		x.Authority.Host = u.Hostname()
		if x.Authority.Host == "" {
			return errors.New("failed to parse url")
		}
	}

	x.Path = u.Path
	x.Fragment = u.Fragment

	x.Query = &Url_Query{
		Vals: make(map[string]*StringValues),
	}

	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	for k, v := range query {
		x.Query.Vals[k] = &StringValues{
			Vals: v,
		}
	}

	return nil
}

func (x *Url) Format() string {
	if x == nil {
		return ""
	}

	var user *url.Userinfo // default nil
	host := ""
	if x.Authority != nil {
		host = x.Authority.Host
		if x.Authority.Port != "" {
			host += ":" + x.Authority.Port
		}

		if x.Authority.UserInfo != "" {
			userinfo := x.Authority.UserInfo
			segments := strings.Split(userinfo, ":")
			if len(segments) == 1 {
				user = url.User(segments[0])
			} else {
				user = url.UserPassword(segments[0], segments[1])
			}
		}
	}

	u := url.URL{
		Scheme:   x.Scheme,
		User:     user,
		Host:     host,
		Path:     x.Path,
		Fragment: x.Fragment,
	}

	if x.Query != nil {
		query := url.Values{}
		for k, v := range x.Query.Vals {
			if v != nil {
				query[k] = v.Vals
			}
		}
		u.RawQuery = query.Encode()
	}

	if u.Scheme == "" {
		return strings.TrimPrefix(u.String(), "//")
	}

	return u.String()
}

func (x *Url) ToString() string {
	return x.Format()
}

func (x *Url) FormatWithoutSchema() string {
	if x == nil {
		return ""
	}

	u := &Url{
		Path:     x.Path,
		Query:    x.Query,
		Fragment: x.Fragment,
	}

	if x.Authority != nil {
		u.Authority = &Url_Authority{
			UserInfo: x.GetAuthority().GetUserInfo(),
			Host:     x.GetAuthority().GetHost(),
			Port:     x.GetAuthority().GetPort(),
		}
	}

	return strings.TrimPrefix(u.Format(), "//")
}
