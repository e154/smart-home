// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package urlpath

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func ExampleNew() {
	path := New("users/:id/files/*")
	fmt.Println(path.Segments[0].Const)
	fmt.Println(path.Segments[1].Param)
	fmt.Println(path.Segments[2].Const)
	fmt.Println(path.Trailing)
	// Output:
	//
	// users
	// id
	// files
	// true
}

func ExampleMatch() {
	path := New("users/:id/files/*")
	match, ok := path.Match("users/123/files/foo/bar/baz.txt")
	fmt.Println(ok)
	fmt.Println(match.Params)
	fmt.Println(match.Trailing)
	// Output:
	//
	// true
	// map[id:123]
	// foo/bar/baz.txt
}

func ExamplePath_Build() {
	path := New("users/:id/files/*")
	built, ok := path.Build(Match{
		Params:   map[string]string{"id": "123"},
		Trailing: "foo/bar/baz.txt",
	})
	fmt.Println(ok)
	fmt.Println(built)
	// Output:
	//
	// true
	// users/123/files/foo/bar/baz.txt
}

func TestNew(t *testing.T) {
	testCases := []struct {
		in  string
		out Path
	}{
		{
			"foo",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
			}},
		},

		{
			"/foo",
			Path{Segments: []Segment{
				Segment{Const: ""},
				Segment{Const: "foo"},
			}},
		},

		{
			":foo",
			Path{Segments: []Segment{
				Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"/:foo",
			Path{Segments: []Segment{
				Segment{Const: ""},
				Segment{IsParam: true, Param: "foo"},
			}},
		},

		{
			"foo/:bar",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:foo/bar/:bar",
			Path{Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "foo"},
				Segment{Const: "bar"},
				Segment{IsParam: true, Param: "bar"},
			}},
		},

		{
			"foo/:bar/:baz/*",
			Path{Trailing: true, Segments: []Segment{
				Segment{Const: "foo"},
				Segment{IsParam: true, Param: "bar"},
				Segment{IsParam: true, Param: "baz"},
			}},
		},

		{
			"/:/*",
			Path{Trailing: true, Segments: []Segment{
				Segment{Const: ""},
				Segment{IsParam: true, Param: ""},
			}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.in, func(t *testing.T) {
			out := New(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		Path string
		in   string
		out  Match
		ok   bool
	}{
		{
			"foo",
			"foo",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"foo",
			"bar",
			Match{},
			false,
		},

		{
			":foo",
			"bar",
			Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo",
			"/bar",
			Match{Params: map[string]string{"foo": "bar"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bar/baz",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz",
			Match{},
			false,
		},

		{
			"/:foo/:bar/:baz",
			"/foo/bar/baz",
			Match{Params: map[string]string{"foo": "foo", "bar": "bar", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"///",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz",
			"",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/a/b/c",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz",
			"/foo/bax/baz/",
			Match{},
			false,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/a/b/c",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: "a/b/c"},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz/",
			Match{Params: map[string]string{"foo": "foo", "baz": "baz"}, Trailing: ""},
			true,
		},

		{
			"/:foo/bar/:baz/*",
			"/foo/bar/baz",
			Match{},
			false,
		},

		{
			"/:foo/:bar/:baz/*",
			"////",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: ""},
			true,
		},

		{
			"/:foo/:bar/:baz/*",
			"/////",
			Match{Params: map[string]string{"foo": "", "bar": "", "baz": ""}, Trailing: "/"},
			true,
		},

		{
			"*",
			"",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"/*",
			"",
			Match{},
			false,
		},

		{
			"*",
			"/",
			Match{Params: map[string]string{}, Trailing: "/"},
			true,
		},

		{
			"/*",
			"/",
			Match{Params: map[string]string{}, Trailing: ""},
			true,
		},

		{
			"*",
			"/a/b/c",
			Match{Params: map[string]string{}, Trailing: "/a/b/c"},
			true,
		},

		{
			"*",
			"a/b/c",
			Match{Params: map[string]string{}, Trailing: "a/b/c"},
			true,
		},

		// Examples from documentation
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar",
			Match{Params: map[string]string{"shelf": "foo", "book": "bar"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/456",
			Match{Params: map[string]string{"shelf": "123", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/123/books/",
			Match{Params: map[string]string{"shelf": "123", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/456",
			Match{Params: map[string]string{"shelf": "", "book": "456"}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves//books/",
			Match{Params: map[string]string{"shelf": "", "book": ""}},
			true,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/bar/",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/shelves/foo/books/pages/baz",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"/SHELVES/foo/books/bar",
			Match{},
			false,
		},
		{
			"/shelves/:shelf/books/:book",
			"shelves/foo/books/bar",
			Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/",
			Match{Params: map[string]string{"user": "foo"}, Trailing: ""},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files/foo/bar/baz.txt",
			Match{Params: map[string]string{"user": "foo"}, Trailing: "foo/bar/baz.txt"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files////",
			Match{Params: map[string]string{"user": "foo"}, Trailing: "///"},
			true,
		},
		{
			"/users/:user/files/*",
			"/users/foo",
			Match{},
			false,
		},
		{
			"/users/:user/files/*",
			"/users/foo/files",
			Match{},
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tt.Path, tt.in), func(t *testing.T) {
			path := New(tt.Path)
			out, ok := path.Match(tt.in)

			if !reflect.DeepEqual(out, tt.out) {
				t.Errorf("out %#v, want %#v", out, tt.out)
			}

			if ok != tt.ok {
				t.Errorf("ok %#v, want %#v", ok, tt.ok)
			}

			// If no error was expected when matching the data, then we should be able
			// to round-trip back to the original data using Build.
			if tt.ok {
				if in, ok := path.Build(out); ok {
					if in != tt.in {
						t.Errorf("in %#v, want %#v", in, tt.in)
					}
				} else {
					t.Error("Build returned ok = false")
				}
			}
		})
	}
}

func BenchmarkMatch(b *testing.B) {
	b.Run("without trailing segments", func(b *testing.B) {
		b.Run("urlpath", func(b *testing.B) {
			path := New("/test/:foo/bar/:baz")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})

		b.Run("regex", func(b *testing.B) {
			regex := regexp.MustCompile("/test/(?P<foo>[^/]+)/bar/(?P<baz>[^/]+)")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})
	})

	b.Run("with trailing segments", func(b *testing.B) {
		b.Run("urlpath", func(b *testing.B) {
			path := New("/test/:foo/bar/:baz/*")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				path.Match(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})

		b.Run("regex", func(b *testing.B) {
			regex := regexp.MustCompile("/test/(?P<foo>[^/]+)/bar/(?P<baz>[^/]+)/.*")
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d", i, i))
				regex.FindStringSubmatch(fmt.Sprintf("/test/foo%d/bar/baz%d/extra", i, i))
			}
		})
	})
}
