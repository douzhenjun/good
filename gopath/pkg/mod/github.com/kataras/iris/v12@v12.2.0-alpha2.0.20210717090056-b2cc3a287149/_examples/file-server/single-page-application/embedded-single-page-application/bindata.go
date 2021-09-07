// Code generated by go-bindata. (@generated) DO NOT EDIT.

// Package main generated by go-bindata.// sources:
// data/public/app.js
// data/public/app2/index.html
// data/public/css/main.css
// data/views/index.html
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// ModTime return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

type assetFile struct {
	*bytes.Reader
	name            string
	childInfos      []os.FileInfo
	childInfoOffset int
}

type assetOperator struct{}

// Open implement http.FileSystem interface
func (f *assetOperator) Open(name string) (http.File, error) {
	var err error
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	content, err := Asset(name)
	if err == nil {
		return &assetFile{name: name, Reader: bytes.NewReader(content)}, nil
	}
	children, err := AssetDir(name)
	if err == nil {
		childInfos := make([]os.FileInfo, 0, len(children))
		for _, child := range children {
			childPath := filepath.Join(name, child)
			info, errInfo := AssetInfo(filepath.Join(name, child))
			if errInfo == nil {
				childInfos = append(childInfos, info)
			} else {
				childInfos = append(childInfos, newDirFileInfo(childPath))
			}
		}
		return &assetFile{name: name, childInfos: childInfos}, nil
	} else {
		// If the error is not found, return an error that will
		// result in a 404 error. Otherwise the server returns
		// a 500 error for files not found.
		if strings.Contains(err.Error(), "not found") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
}

// Close no need do anything
func (f *assetFile) Close() error {
	return nil
}

// Readdir read dir's children file info
func (f *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if len(f.childInfos) == 0 {
		return nil, os.ErrNotExist
	}
	if count <= 0 {
		return f.childInfos, nil
	}
	if f.childInfoOffset+count > len(f.childInfos) {
		count = len(f.childInfos) - f.childInfoOffset
	}
	offset := f.childInfoOffset
	f.childInfoOffset += count
	return f.childInfos[offset : offset+count], nil
}

// Stat read file info from asset item
func (f *assetFile) Stat() (os.FileInfo, error) {
	if len(f.childInfos) != 0 {
		return newDirFileInfo(f.name), nil
	}
	return AssetInfo(f.name)
}

// newDirFileInfo return default dir file info
func newDirFileInfo(name string) os.FileInfo {
	return &bindataFileInfo{
		name:    name,
		size:    0,
		mode:    os.FileMode(2147484068), // equal os.FileMode(0644)|os.ModeDir
		modTime: time.Time{}}
}

// AssetFile return a http.FileSystem instance that data backend by asset
func AssetFile() http.FileSystem {
	return &assetOperator{}
}

var _dataPublicAppJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x2a\xcf\xcc\x4b\xc9\x2f\xd7\x4b\xcc\x49\x2d\x2a\xd1\x50\x4a\x2c\x28\xd0\xcb\x2a\x56\xc8\xc9\x4f\x4c\x49\x4d\x51\x48\x2b\xca\xcf\x55\x88\x51\xd2\x57\xd2\xb4\x06\x04\x00\x00\xff\xff\xa9\x06\xf7\xa3\x27\x00\x00\x00")

func dataPublicAppJsBytes() ([]byte, error) {
	return bindataRead(
		_dataPublicAppJs,
		"data/public/app.js",
	)
}

func dataPublicAppJs() (*asset, error) {
	bytes, err := dataPublicAppJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/public/app.js", size: 39, mode: os.FileMode(438), modTime: time.Unix(1599156854, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataPublicApp2IndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\xce\xb1\x0a\x02\x31\x0c\xc6\xf1\x3d\x90\x77\xf8\x46\x5d\x2c\x77\x73\x28\xb8\xdd\xa0\x20\xe8\x0b\x54\x1b\x6d\xa1\x67\x8b\x64\xd0\xb7\x97\x3b\xeb\x18\xf8\xf3\xcb\x27\xc9\xe6\xe2\x99\x98\x24\x69\x88\x9e\x09\x00\xc4\xb2\x15\xf5\xfb\xd6\x30\x8a\xfb\x1d\x4c\xe2\x7a\xc2\x24\xd7\x1a\x3f\xff\x38\x0d\x1e\x9b\xb3\x05\xcb\x37\x4c\x97\xe3\x01\xa7\xf0\xd0\x2d\x26\x2d\xa5\xe2\xfe\xaa\x33\x42\x6b\xa3\xcb\xcf\xa8\xef\xdd\xf2\x0f\xe2\xd2\xb0\x82\x9d\x59\xed\x65\xc8\x37\x00\x00\xff\xff\xf4\x87\x93\x1a\x8f\x00\x00\x00")

func dataPublicApp2IndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_dataPublicApp2IndexHtml,
		"data/public/app2/index.html",
	)
}

func dataPublicApp2IndexHtml() (*asset, error) {
	bytes, err := dataPublicApp2IndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/public/app2/index.html", size: 143, mode: os.FileMode(438), modTime: time.Unix(1600097514, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataPublicCssMainCss = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\xca\x4f\xa9\x54\xa8\xe6\xe5\x52\x50\x50\x50\x48\x4a\x4c\xce\x4e\x2f\xca\x2f\xcd\x4b\xd1\x4d\xce\xcf\xc9\x2f\xb2\x52\x48\xca\x49\x4c\xce\xb6\xe6\xe5\xaa\xe5\xe5\x02\x04\x00\x00\xff\xff\x03\x25\x9c\x89\x29\x00\x00\x00")

func dataPublicCssMainCssBytes() ([]byte, error) {
	return bindataRead(
		_dataPublicCssMainCss,
		"data/public/css/main.css",
	)
}

func dataPublicCssMainCss() (*asset, error) {
	bytes, err := dataPublicCssMainCssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/public/css/main.css", size: 41, mode: os.FileMode(438), modTime: time.Unix(1599156854, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dataViewsIndexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\xce\x3d\x0e\xc2\x30\x0c\x05\xe0\x3d\x52\xee\xf0\xd4\x09\x96\x58\xdd\x4d\x66\x46\x86\x5e\x20\xb4\x86\x04\xa5\x34\x6a\x22\x7e\x54\xf5\xee\xa8\x3f\x8c\xd6\xf3\xf3\x67\xf6\xa5\x8f\x56\x2b\xad\xd8\x8b\xeb\xac\x56\x00\xc0\x25\x94\x28\x76\x9a\x60\x2e\xee\x2e\xa6\x59\x46\xcc\x33\xd3\x16\x68\xc5\xb4\xaf\x6b\xc5\xd7\xa1\xfb\xfe\x8b\xbe\xb6\x38\x34\xd2\xa7\xe8\x8a\x1c\x71\x96\x18\x07\xdc\xc6\xa1\xc7\x2b\xc8\x3b\x53\x78\x76\xf2\x31\x0b\x0a\x26\x5f\xaf\x07\xf6\x6a\x6e\xc7\x90\x0a\xf2\xd8\x9e\x2a\x72\x29\x99\x47\xae\x2c\xc0\xb4\x05\x2b\xba\x53\xab\xbf\x3c\xfe\x0b\x00\x00\xff\xff\x4a\xf7\x07\xf6\xbf\x00\x00\x00")

func dataViewsIndexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_dataViewsIndexHtml,
		"data/views/index.html",
	)
}

func dataViewsIndexHtml() (*asset, error) {
	bytes, err := dataViewsIndexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "data/views/index.html", size: 191, mode: os.FileMode(438), modTime: time.Unix(1600097531, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"data/public/app.js":          dataPublicAppJs,
	"data/public/app2/index.html": dataPublicApp2IndexHtml,
	"data/public/css/main.css":    dataPublicCssMainCss,
	"data/views/index.html":       dataViewsIndexHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("nonexistent") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"data": {nil, map[string]*bintree{
		"public": {nil, map[string]*bintree{
			"app.js": {dataPublicAppJs, map[string]*bintree{}},
			"app2": {nil, map[string]*bintree{
				"index.html": {dataPublicApp2IndexHtml, map[string]*bintree{}},
			}},
			"css": {nil, map[string]*bintree{
				"main.css": {dataPublicCssMainCss, map[string]*bintree{}},
			}},
		}},
		"views": {nil, map[string]*bintree{
			"index.html": {dataViewsIndexHtml, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}