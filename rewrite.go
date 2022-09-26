package rewrite

import (
	"net/http"
	"net/url"
	"regexp"
)

const headerField = "X-Rewrite-Original-URI"

type Rule struct {
	Pattern string
	To      string
	*regexp.Regexp
}

func NewRule(pattern, to string) (*Rule, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &Rule{
		pattern,
		to,
		reg,
	}, nil
}

func (r *Rule) Rewrite(req *http.Request) bool {
	oriPath := req.URL.String()
	if !r.MatchString(oriPath) {
		return false
	}

	to := r.ReplaceAllString(oriPath, r.To)
	uo, e := url.Parse(to)
	if e != nil {
		return false
	}

	req.URL.Path = uo.Path
	req.URL.RawPath = uo.RawPath
	return true
}

func NewHandler(rules map[string]string) RewriteHandler {
	var h RewriteHandler

	for key, val := range rules {
		r, e := NewRule(key, val)
		if e != nil {
			panic(e)
		}

		h.rules = append(h.rules, r)
	}

	return h
}

type RewriteHandler struct {
	rules []*Rule
}

func (h *RewriteHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, r := range h.rules {
		ok := r.Rewrite(req)
		if ok {
			break
		}
	}
}

func hit(pattern, str string) bool {
	r, e := regexp.MatchString(pattern, str)
	if e != nil {
		return false
	}

	return r
}
