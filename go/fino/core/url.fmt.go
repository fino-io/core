package core

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

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

	if u.Opaque != "" {
		return errors.New("opaque URL is not supported by the Url contract")
	}

	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	x.Scheme = u.Scheme
	x.Authority = &Url_Authority{
		UserInfo: u.User.String(),
		Host:     u.Hostname(),
		Port:     u.Port(),
	}
	x.Path = u.Path
	x.Fragment = u.Fragment
	x.Query = NewUrlQueryFrom(query)

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
			host = net.JoinHostPort(host, x.Authority.Port)
		} else if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}

		if x.Authority.UserInfo != "" {
			userinfo, err := url.Parse("//" + x.Authority.UserInfo + "@localhost")
			if err == nil {
				user = userinfo.User
			} else {
				user = url.User(x.Authority.UserInfo)
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

	return u.String()
}

func (x *Url) ToString() string {
	return x.Format()
}

// FormatWithoutScheme returns a display form without the scheme or leading //.
func (x *Url) FormatWithoutScheme() string {
	if x == nil {
		return ""
	}

	u := &Url{
		Authority: x.Authority,
		Path:      x.Path,
		Query:     x.Query,
		Fragment:  x.Fragment,
	}

	return strings.TrimPrefix(u.Format(), "//")
}
