package httpext

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
)

func Abort(w http.ResponseWriter, err error, code int) {
	response := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}

func AbortWithoutContent(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func AbortPlain(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	String(w, err.Error())
}

func AbortJSON(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func JSONWithStatus(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func JSON(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func XML(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	_, _ = w.Write([]byte(xml.Header))
	_ = xml.NewEncoder(w).Encode(body)
}

func Data(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", http.DetectContentType(data))
	_, _ = w.Write(data)
}

func String(w http.ResponseWriter, str string) {
	w.Header().Add("Content-Type", "text/plain")
	_, _ = w.Write([]byte(str))
}

func SanitizeURL(uri string) string {
	cleanRegexp, _ := regexp.Compile("[\\d\\w\\-/]+")
	return cleanRegexp.FindString(uri)
}

func MakeAbsoluteURL(r *http.Request, target string, isSecured bool) string {
	if isSecured {
		return fmt.Sprintf("https://%s%s", r.Host, target)
	} else {
		return fmt.Sprintf("http://%s%s", r.Host, target)
	}
}

func GetHostname(r *http.Request) string {
	host := r.Host
	if strings.Contains(host, ":") {
		host, _, _ = net.SplitHostPort(host)
	}
	return host
}

func GenerateURL(
	r *http.Request,
	absolute bool,
	isSecured bool) string {
	ctx := r.Context()
	// this function only supports generating the current route at the moment
	// a next step could be to add the ability to insert an alias from the routing table
	// to generate the URL
	var paramRegExp *regexp.Regexp

	if rctx := chi.RouteContext(ctx); rctx != nil {
		url := rctx.RoutePattern()

		// remove last slash for main pages (with and without lang prefix)
		if url == "/" {
			url = ""
		}

		//to avoid runtime panic when reconstructing URL, leave it
		//as requested URL in case of incorrect url to provide 500 page on it
		if len(rctx.URLParams.Keys) != len(rctx.URLParams.Values) {
			return r.URL.RequestURI()
		}

		for k := len(rctx.URLParams.Keys) - 1; k >= 0; k-- {
			search := rctx.URLParams.Keys[k]

			if search == "*" {
				paramRegExp, _ = regexp.Compile("\\*")
			} else {
				paramRegExp, _ = regexp.Compile("{" + search + ":?.*}")
			}

			url = string(paramRegExp.ReplaceAll([]byte(url), []byte(rctx.URLParams.Values[k])))
		}
		if absolute {
			url = MakeAbsoluteURL(r, url, isSecured)
		}
		return url
	}
	return ""
}
