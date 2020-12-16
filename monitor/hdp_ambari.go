package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

//Ambaris Ambaris
type Ambaris struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	Path string `json:"path"`
	User
}

//User upmgj
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//MakeAmbaris  代理
func MakeAmbaris(ip, username, password, port, path string) (*Ambaris, error) {
	if username == "" || password == "" {
		return nil, errors.New("必须要密码")
	}
	return &Ambaris{
		IP:   ip,
		Port: port,
		Path: path,
	}, nil
}

//Proxy  代理
func (a *Ambaris) Proxy() http.Handler {
	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd, _ := forward.New()
	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server
		req.URL = testutils.ParseURI(fmt.Sprintf("http://%s:%s%s", a.IP, a.Port, a.Path))
		fmt.Println(req.URL)
		fwd.ServeHTTP(w, req)
	})
	return redirect
}

//HTTPProxy  代理
func (a *Ambaris) HTTPProxy() http.HandlerFunc {

	targetURL := &url.URL{
		Scheme: "http",
		//User:   url.UserPassword(a.Username, a.Password),
		Host: a.IP + ":" + a.Port,
	}
	fmt.Println(targetURL)
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		// 设置Content-Type
		// 检查是否POST请求
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}
		user := &User{}
		json.NewDecoder(r.Body).Decode(user)
		if user.Username != "sail" || user.Password != "1234" {
			w.WriteHeader(406)
			return
		}
		o := new(http.Request)
		*o = *r
		o.Host = targetURL.Host
		o.URL.Scheme = targetURL.Scheme
		o.URL.Host = targetURL.Host
		o.URL.Path = singleJoiningSlash(targetURL.Path, o.URL.Path)
		if q := o.URL.RawQuery; q != "" {
			o.URL.RawPath = o.URL.Path + "?" + q
		} else {
			o.URL.RawPath = o.URL.Path
		}
		o.URL.RawQuery = targetURL.RawQuery
		o.Proto = "HTTP/1.1"
		o.ProtoMajor = 1
		o.ProtoMinor = 1
		o.Close = false
		transport := http.DefaultTransport
		res, err := transport.RoundTrip(o)
		if err != nil {
			log.Printf("http: proxy error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		hdr := w.Header()
		for k, vv := range res.Header {
			for _, v := range vv {
				hdr.Add(k, v)
			}
		}
		// for _, c := range res.SetCookie {
		// 	w.Header().Add("Set-Cookie", c.Raw)
		// }
		w.WriteHeader(res.StatusCode)
		if res.Body != nil {
			io.Copy(w, res.Body)
		}
	}
}
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
