package rewrite

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testFixture struct {
	from string
	to   string
}

type testCase struct {
	pattern  string
	to       string
	fixtures []testFixture
}

func TestRewrite0(t *testing.T) {
	tests := []testCase{
		{
			pattern: "(_region=)CN",
			to:      "${1}KR",
			fixtures: []testFixture{
				{from: "/a?_region=CN", to: "/a?_region=KR"},
			},
		},
	}

	for _, test := range tests {
		t.Logf("Test - pattern: %s, to: %s", test.pattern, test.to)

		for _, fixture := range test.fixtures {
			req, err := http.NewRequest("GET", fixture.from, nil)
			if err != nil {
				t.Fatalf("Fixture %v - create HTTP request error: %v", fixture, err)
			}

			h := NewHandler(map[string]string{
				test.pattern: test.to,
			})

			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)

			t.Logf("Rewrited: %s", req.URL.String())
		}
	}
}

func TestRewrite(t *testing.T) {
	tests := []testCase{
		{
			pattern: "/a",
			to:      "/b",
			fixtures: []testFixture{
				{from: "/a", to: "/b"},
			},
		},
		{
			pattern: "/a/(.*)",
			to:      "/bb",
			fixtures: []testFixture{
				{from: "/a", to: "/a"},
				{from: "/a/", to: "/bb"},
				{from: "/a/a", to: "/bb"},
				{from: "/a/b/c", to: "/bb"},
			},
		},
		{
			pattern: "/r/(.*)",
			to:      `/r/v1/$1`,
			fixtures: []testFixture{
				{from: "/a", to: "/a"},
				{from: "/r", to: "/r"},
				{from: "/r/a", to: "/r/v1/a"},
				{from: "/r/a/b", to: "/r/v1/a/b"},
			},
		},
		{
			pattern: "/r/(.*)/a/(.*)",
			to:      `/r/v1/$1/a/$2`,
			fixtures: []testFixture{
				{from: "/r/1/2", to: "/r/1/2"},
				{from: "/r/1/a/2", to: "/r/v1/1/a/2"},
				{from: "/r/1/a/2/3", to: "/r/v1/1/a/2/3"},
			},
		},
		{
			pattern: "/r/(.*)/a/(.*)",
			to:      `/r/v1/$2/a/$1`,
			fixtures: []testFixture{
				{from: "/r/1/a/2", to: "/r/v1/2/a/1"},
				{from: "/r/1/a/2/3", to: "/r/v1/2/3/a/1"},
			},
		},
		{
			pattern: "/from/:one/to/:two",
			to:      "/from/:two/to/:one",
			fixtures: []testFixture{
				{from: "/from/123/to/456", to: "/from/456/to/123"},
				{from: "/from/abc/to/def", to: "/from/def/to/abc"},
			},
		},
		{
			pattern: "/from/:one/to/:two",
			to:      "/:one/:two/:three/:two/:one",
			fixtures: []testFixture{
				{from: "/from/123/to/456", to: "/123/456/:three/456/123"},
				{from: "/from/abc/to/def", to: "/abc/def/:three/def/abc"},
			},
		},
		{
			pattern: "/from/(.*)",
			to:      "/to/$1",
			fixtures: []testFixture{
				{from: "/from/untitled-1%2F/upload", to: "/to/untitled-1%2F/upload"},
			},
		},
	}

	for _, test := range tests {
		t.Logf("Test - pattern: %s, to: %s", test.pattern, test.to)

		for _, fixture := range test.fixtures {
			req, err := http.NewRequest("GET", fixture.from, nil)
			if err != nil {
				t.Fatalf("Fixture %v - create HTTP request error: %v", fixture, err)
			}

			h := NewHandler(map[string]string{
				test.pattern: test.to,
			})

			t.Logf("From: %s", req.URL.EscapedPath())
			if req.URL.EscapedPath() != fixture.from {
				t.Errorf("Invalid test fixture: %s", fixture.from)
			}

			res := httptest.NewRecorder()
			h.ServeHTTP(res, req)

			t.Logf("Rewrited: %s", req.URL.EscapedPath())
			if req.URL.EscapedPath() != fixture.to {
				t.Errorf("Test failed \n pattern: %s, to: %s, \n fixture: %s to %s, \n result: %s",
					test.pattern, test.to, fixture.from, fixture.to, req.URL.EscapedPath())
			}

			if req.Header.Get(headerField) != "" {
				// matched
				if req.Header.Get(headerField) != fixture.from {
					t.Error("incorrect flag")
				}
			}
		}
	}
}
