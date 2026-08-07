package core

import (
	"reflect"
	"testing"
)

func TestNewUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *Url
		wantErr bool
	}{
		{name: "empty", url: "", want: &Url{Authority: &Url_Authority{}, Query: NewUrlQuery()}, wantErr: false},
		{name: "http-1", url: "http://g.cn", want: &Url{Scheme: "http", Authority: &Url_Authority{Host: "g.cn"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "https-1", url: "https://g.cn", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "mailto-1", url: "mailto://g.cn", want: &Url{Scheme: "mailto", Authority: &Url_Authority{Host: "g.cn"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "file-1", url: "file://g.cn", want: &Url{Scheme: "file", Authority: &Url_Authority{Host: "g.cn"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "https-port", url: "https://g.cn:8081", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn", Port: "8081"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "https-userinfo-no-password", url: "https://root@g.cn:8081", want: &Url{Scheme: "https", Authority: &Url_Authority{UserInfo: "root", Host: "g.cn", Port: "8081"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "https-userinfo-has-password", url: "https://root:123456@g.cn:8081", want: &Url{Scheme: "https", Authority: &Url_Authority{UserInfo: "root:123456", Host: "g.cn", Port: "8081"}, Query: NewUrlQuery()}, wantErr: false},
		{name: "https-path-1", url: "https://g.cn/resources", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery()}, wantErr: false},
		{name: "https-query-1", url: "https://g.cn/resources?key1=value1", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1")}, wantErr: false},
		{name: "https-query-2", url: "https://g.cn/resources?key1=value1&key2=value2", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1", "key2", "value2")}, wantErr: false},
		{name: "https-fragment-1", url: "https://g.cn/resources#top", want: &Url{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery(), Fragment: "top"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUrl(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrl_ToString(t *testing.T) {
	type fields struct {
		Scheme    string
		Authority *Url_Authority
		Path      string
		Query     *Url_Query
		Fragment  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "empty", fields: fields{}, want: ""},
		{name: "no-scheme", fields: fields{Authority: &Url_Authority{Host: "g.cn"}}, want: "g.cn"},
		{name: "http-scheme", fields: fields{Scheme: "http", Authority: &Url_Authority{Host: "g.cn"}}, want: "http://g.cn"},
		{name: "https-scheme", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}}, want: "https://g.cn"},
		{name: "mailto-scheme", fields: fields{Scheme: "mailto", Authority: &Url_Authority{Host: "g.cn"}}, want: "mailto://g.cn"},
		{name: "https-userinfo-no-password", fields: fields{Scheme: "https", Authority: &Url_Authority{UserInfo: "root", Host: "g.cn"}}, want: "https://root@g.cn"},
		{name: "https-userinfo-has-password", fields: fields{Scheme: "https", Authority: &Url_Authority{UserInfo: "root:123456", Host: "g.cn"}}, want: "https://root:123456@g.cn"},
		{name: "https-path-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery()}, want: "https://g.cn/resources"},
		{name: "https-query-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1")}, want: "https://g.cn/resources?key1=value1"},
		{name: "https-query-2", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1", "key2", "value2")}, want: "https://g.cn/resources?key1=value1&key2=value2"},
		{name: "https-fragment-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery(), Fragment: "top"}, want: "https://g.cn/resources#top"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Url{
				Scheme:    tt.fields.Scheme,
				Authority: tt.fields.Authority,
				Path:      tt.fields.Path,
				Query:     tt.fields.Query,
				Fragment:  tt.fields.Fragment,
			}
			if got := x.ToString(); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrl_FormatWithoutSchema(t *testing.T) {
	type fields struct {
		Scheme    string
		Authority *Url_Authority
		Path      string
		Query     *Url_Query
		Fragment  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "empty", fields: fields{}, want: ""},
		{name: "no-scheme", fields: fields{Authority: &Url_Authority{Host: "g.cn"}}, want: "g.cn"},
		{name: "http-scheme", fields: fields{Scheme: "http", Authority: &Url_Authority{Host: "g.cn"}}, want: "g.cn"},
		{name: "https-scheme", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}}, want: "g.cn"},
		{name: "mailto-scheme", fields: fields{Scheme: "mailto", Authority: &Url_Authority{Host: "g.cn"}}, want: "g.cn"},
		{name: "https-userinfo-no-password", fields: fields{Scheme: "https", Authority: &Url_Authority{UserInfo: "root", Host: "g.cn"}}, want: "root@g.cn"},
		{name: "https-userinfo-has-password", fields: fields{Scheme: "https", Authority: &Url_Authority{UserInfo: "root:123456", Host: "g.cn"}}, want: "root:123456@g.cn"},
		{name: "https-path-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery()}, want: "g.cn/resources"},
		{name: "https-query-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1")}, want: "g.cn/resources?key1=value1"},
		{name: "https-query-2", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery("key1", "value1", "key2", "value2")}, want: "g.cn/resources?key1=value1&key2=value2"},
		{name: "https-fragment-1", fields: fields{Scheme: "https", Authority: &Url_Authority{Host: "g.cn"}, Path: "/resources", Query: NewUrlQuery(), Fragment: "top"}, want: "g.cn/resources#top"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Url{
				Scheme:    tt.fields.Scheme,
				Authority: tt.fields.Authority,
				Path:      tt.fields.Path,
				Query:     tt.fields.Query,
				Fragment:  tt.fields.Fragment,
			}
			if got := x.FormatWithoutSchema(); got != tt.want {
				t.Errorf("FormatWithoutSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
