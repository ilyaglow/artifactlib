// Copyright (c) 2019 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

package goartifacts

import (
	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/filesystem/systemfs"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/forensicanalysis/fslib/filesystem/testfs"
)

func getInFS() fslib.FS {
	infs := &testfs.FS{}
	content := []byte("test")
	dirs := []string{"/dir/", "/dir/a/", "/dir/b/", "/dir/a/a/", "/dir/a/b/", "/dir/b/a/", "/dir/b/b/"}
	for _, dir := range dirs {
		infs.CreateDir(dir)
	}
	files := []string{"/foo.bin", "/dir/bar.bin", "/dir/baz.bin", "/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}
	for _, file := range files {
		infs.CreateFile(file, content)
	}
	return infs
}

func TestExpand(t *testing.T) {
	type args struct {
		infs                fslib.FS
		artifactDefinitions []ArtifactDefinition
	}
	tests := []struct {
		name    string
		args    args
		want    []ArtifactDefinition
		wantErr bool
	}{
		{
			"Expand", args{
			getInFS(),
			[]ArtifactDefinition{
				{Sources: []Source{{Type: "FILE", Attributes: Attributes{Paths: []string{"/*/bar.bin"}}}}},
			},
		},
			[]ArtifactDefinition{
				{Sources: []Source{{Type: "FILE", Attributes: Attributes{Paths: []string{"/dir/bar.bin"}}}}},
			}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.args.infs, tt.args.artifactDefinitions, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expand() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_expandPath(t *testing.T) {
	validPath, err := osfs.ToForensicPath("../test/artifacts/valid")
	if err != nil {
		t.Fatal(err)
	}
	invalidPath, err := osfs.ToForensicPath("../test/artifacts/invalid")
	if err != nil {
		t.Fatal(err)
	}
	winfs, err := systemfs.New()
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		fs fslib.FS
		in string
	}
	tests := []struct {
		name        string
		args        args
		want        []string
		windowsOnly bool
	}{
		{"Expand path 1", args{getInFS(), "/*/bar.bin"}, []string{"/dir/bar.bin"}, false},
		{"Expand path 2", args{getInFS(), "/dir/*.bin"}, []string{"/dir/bar.bin", "/dir/baz.bin"}, false},
		{"Expand path 3", args{getInFS(), "/dir/*/*/foo.bin"}, []string{"/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}, false},
		{"Expand path 4", args{getInFS(), "/**"}, []string{"/dir", "/dir/a", "/dir/a/a", "/dir/a/b", "/dir/b", "/dir/b/a", "/dir/b/b", "/dir/bar.bin", "/dir/baz.bin", "/foo.bin"}, false},
		{"Expand path 5", args{getInFS(), "/dir/**1"}, []string{"/dir/a", "/dir/b", "/dir/bar.bin", "/dir/baz.bin"}, false},
		{"Expand path 7", args{getInFS(), "/dir/**10"}, []string{"/dir/a", "/dir/a/a", "/dir/a/a/foo.bin", "/dir/a/b", "/dir/a/b/foo.bin", "/dir/b", "/dir/b/a", "/dir/b/a/foo.bin", "/dir/b/b", "/dir/b/b/foo.bin", "/dir/bar.bin", "/dir/baz.bin"}, false},
		{"Expand OSpath", args{osfs.New(), "../test/artifacts/*lid"}, []string{validPath, invalidPath}, false},
		{"Expand win path", args{osfs.New(), "C:/Windows"}, []string{"/C/Windows"}, true},
		{"Expand special file path", args{winfs, "C:/$MFT"}, []string{"/C/$MFT"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.windowsOnly || runtime.GOOS == "windows" {
				got, err := expandPath(tt.args.fs, tt.args.in, tt.args.fs.Name() == "OsFs" || tt.args.fs.Name() == "System FS")
				if err != nil {
					t.Fatal(err)
				}
				sort.Strings(tt.want)
				sort.Strings(got)
				if !reflect.DeepEqual(got, tt.want) {
					t.Error("are you admin?")
					t.Errorf("expandPath(%s) = %v, want %v", tt.args.in, got, tt.want)
				}
			}
		})
	}
}

func isSubset(superset []string, subset []string) bool {
	for _, subsetElem := range subset {
		if !contains(superset, subsetElem) {
			return false
		}
	}
	return true
}

func contains(set []string, elem string) bool {
	for _, a := range set {
		if a == elem {
			return true
		}
	}
	return false
}