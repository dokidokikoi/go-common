package tools

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"golang.org/x/net/proxy"
)

func Socks5Proxy(p string, username, password string) (proxy.Dialer, error) {
	dialer, err := proxy.SOCKS5("tcp", p, &proxy.Auth{
		User:     username,
		Password: password,
	}, proxy.Direct)
	if err != nil {
		return nil, err
	}

	// return func(ctx context.Context, network, addr string) (net.Conn, error) { return dialer.Dial(network, addr) }, nil
	return dialer, err
}

func AbsUrl(base, uri string) string {
	baseURL, _ := url.Parse(base)
	absURL, _ := baseURL.Parse(uri)
	return absURL.String()
}

func GenQueryParams(p any) string {
	if p == nil {
		return ""
	}
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		p = v.Interface()
	}
	t := reflect.TypeOf(p)
	builder := strings.Builder{}
	for i := t.NumField() - 1; i >= 0; i-- {
		vf := v.Field(i)
		if vf.IsZero() {
			continue
		}
		if vf.Kind() == reflect.Ptr {
			vf = vf.Elem()
		}
		f := t.Field(i)
		if f.Type.Kind() == reflect.Slice {
			for i := 0; i < vf.Len(); i++ {
				builder.WriteString(fmt.Sprintf("%s=%v&", f.Tag.Get("query")+"[]", vf.Index(i).Interface()))
			}
			continue
		}
		builder.WriteString(fmt.Sprintf("%s=%v&", f.Tag.Get("query"), vf.Interface()))
	}
	return builder.String()
}
