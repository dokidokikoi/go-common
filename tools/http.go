package tools

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
	"resty.dev/v3"
)

var pool sync.Pool

func init() {
	pool = sync.Pool{
		New: func() interface{} {
			return resty.New().SetRetryCount(3).SetRetryWaitTime(time.Second)
		},
	}
}

func GetHttpClient() *resty.Client {
	return pool.Get().(*resty.Client)
}

type Option func(*resty.Request) error

func Req(method, url string, body any, options ...Option) (*resty.Response, error) {
	clinet := GetHttpClient()
	defer pool.Put(clinet)

	req := clinet.R().AddRetryConditions(func(r *resty.Response, err error) bool {
		return true
	}).SetBody(body)
	for _, o := range options {
		o(req)
	}
	return req.Execute(method, url)
}

func ReqWithProxy(method, url string, body any, proxy string, options ...Option) (*resty.Response, error) {
	proxy = strings.TrimSpace(proxy)
	if proxy == "" {
		return Req(method, url, body, options...)
	}
	client := GetHttpClient()
	defer func() {
		pool.Put(client.SetProxy(""))
	}()
	client.SetProxy(proxy)
	req := client.R().AddRetryConditions(func(r *resty.Response, err error) bool {
		return true
	}).SetBody(body)
	for _, o := range options {
		o(req)
	}
	return req.Execute(method, url)
}

func SetHeadersWithOption(headers map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetHeaders(headers)
		return nil
	}
}

func SetQueryParamsWithOption(params map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetQueryParams(params)
		return nil
	}
}

func SetCookiesWithOption(cookies ...*http.Cookie) Option {
	return func(r *resty.Request) error {
		r.SetCookies(cookies)
		return nil
	}
}

func SetFromWithOption(data map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetFormData(data)
		return nil
	}
}

func SetMultipartWithOption(fields ...*resty.MultipartField) Option {
	return func(r *resty.Request) error {
		r.SetMultipartFields(fields...)
		return nil
	}
}

func SetMultiFileWithOption(params map[string]string, files map[string][]string) Option {
	return func(r *resty.Request) error {
		buf := bytes.NewBuffer([]byte{})
		writer := multipart.NewWriter(buf)
		for k, v := range params {
			w, err := writer.CreateFormField(k)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte(v))
			if err != nil {
				return err
			}
		}
		for k, v := range files {
			for _, vv := range v {
				_, filename := filepath.Split(vv)
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition",
					fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
						escapeQuotes(k), escapeQuotes(filename)))
				h.Set("Content-Type", "application/octet-stream")
				w, err := writer.CreatePart(h)
				if err != nil {
					return err
				}
				f, err := os.Open(vv)
				if err != nil {
					return err
				}
				_, err = io.Copy(w, f)
				f.Close()
				if err != nil {
					return err
				}
			}
		}
		r.SetBody(buf)
		r.SetHeader("Content-Type", writer.FormDataContentType())
		return nil
	}
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func SetMultipartFormWithOption(params map[string]string) Option {
	return func(r *resty.Request) error {
		r.SetMultipartFormData(params)
		return nil
	}
}

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
		name := f.Tag.Get("query")
		if name == "" || name == "-" {
			continue
		}
		if f.Type.Kind() == reflect.Slice {
			for i := 0; i < vf.Len(); i++ {
				builder.WriteString(fmt.Sprintf("%s=%v&", name+"[]", vf.Index(i).Interface()))
			}
			continue
		}
		builder.WriteString(fmt.Sprintf("%s=%v&", name, vf.Interface()))
	}
	return builder.String()
}
