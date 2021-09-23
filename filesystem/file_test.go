package filesystem

import (
	"bytes"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file"
	"github.com/pingcap/failpoint"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const testDir = "testDir"

func TestEnsureFileExistErr(t *testing.T) {
	var err error
	fp := "test/111.txt"
	testContent := []byte("TestEnsureFileExist")

	{
		err = failpoint.Enable("github.com/meshplus/hyperbench/filesystem/openfile-err", "return")
		assert.NoError(t, err)
		_, err = EnsureFileExist(fp, testContent)
		assert.Error(t, err)
		err = failpoint.Disable("github.com/meshplus/hyperbench/filesystem/openfile-err")
		assert.NoError(t, err)
		_ = os.RemoveAll(fp)
	}

	{
		_ = failpoint.Enable("github.com/meshplus/hyperbench/filesystem/writefile-err", "return")
		_, err = EnsureFileExist(fp, testContent)
		assert.Error(t, err)
		_ = failpoint.Disable("github.com/meshplus/hyperbench/filesystem/writefile-err")
		_ = os.RemoveAll(fp)
	}

	{
		_ = failpoint.Enable("github.com/meshplus/hyperbench/filesystem/closefile-err", "return")
		_, err = EnsureFileExist(fp, testContent)
		assert.Error(t, err)
		_ = failpoint.Disable("github.com/meshplus/hyperbench/filesystem/closefile-err")
		_ = os.RemoveAll(fp)
	}
}

func TestEnsureFileExist(t *testing.T) {
	var err error
	var b []byte
	var fp string
	var testContent []byte
	path := []string{
		"test/test.txt",
		"test.txt",
	}

	content := []byte("TestEnsureFileExist")
	for _, p := range path {
		fp = filepath.Join(testDir, p)
		testContent = append(content, p...)
		_, err = EnsureFileExist(fp, testContent)
		assert.Equal(t, nil, err, "make sure test file exist")

		b, err = ioutil.ReadFile(fp)
		assert.Equal(t, nil, err, "open test file")
		assert.Equal(t, testContent, b, "check test file content")

		_, err = EnsureFileExist(fp, append(testContent, testContent...))
		assert.Equal(t, nil, err, "check exit file")
		b, err = ioutil.ReadFile(fp)
		assert.Equal(t, testContent, b, "check exit file")
		assert.Equal(t, nil, err, "check exit file")
	}

	if err = os.RemoveAll(testDir); err != nil {
		assert.Equal(t, nil, err, "remove test dir")
	}
}

func TestUnpackErr(t *testing.T) {
	var err error

	_ = failpoint.Enable("github.com/meshplus/hyperbench/filesystem/unpack-err", "return")
	err = Unpack("111")
	assert.NoError(t, err)
	_ = failpoint.Disable("github.com/meshplus/hyperbench/filesystem/unpack-err")
}

func TestUnpack(t *testing.T) {
	var err error
	var fp string
	var f http.File
	var buf *bytes.Buffer
	path := []string{
		"test/test.txt",
		"test.txt",
	}
	content := []byte("TestUnpack")

	for _, p := range path {
		fp = filepath.Join(testDir, p)
		err = FileSystem.AddBytes(fp, content)
		assert.Equal(t, nil, err)
	}

	if err = Unpack(testDir); err != nil {
		assert.Equal(t, nil, err, "unpack prefix")
	}

	for _, p := range path {
		fp = filepath.Join(testDir, p)
		t.Logf("check file %v", fp)

		f, err = os.Open(fp)
		assert.Equal(t, nil, err, "open from local file system")

		buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(f)
		assert.Equal(t, nil, err, "read file")

		assert.Equal(t, content, buf.Bytes(), "check file content")
	}

	if err = os.RemoveAll(testDir); err != nil {
		assert.Equal(t, nil, err, "remove test dir")
	}

	if len(path) > 0 {
		fp = filepath.Join(testDir, path[0])
		if err = Unpack(fp); err != nil {
			assert.Equal(t, nil, err, "unpack file path")
		}
		f, _ = os.Open(fp)
		buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(f)
		assert.Equal(t, nil, err, "check file content")
		if err = os.RemoveAll(testDir); err != nil {
			assert.Equal(t, nil, err, "remove test dir")
		}
	}

}

func TestWalkErr(t *testing.T) {
	var err error
	var f packr.File
	var name = "filesystem_test.txt"
	var content []byte

	content = []byte(content)

	f, err = file.NewFile(name, content)
	assert.NoError(t, err)

	_ = failpoint.Enable("github.com/meshplus/hyperbench/filesystem/walk-buf-err", "return")
	err = walk(name, f)
	assert.Error(t, err)
	_ = failpoint.Disable("github.com/meshplus/hyperbench/filesystem/walk-buf-err")
}

func TestWalk(t *testing.T) {
	var err error
	var name string
	var f packr.File
	//var fp string
	var content []byte

	content = []byte(content)
	name = "filesystem_test.txt"

	f, err = file.NewFile(name, content)
	assert.NoError(t, err)
	//fp = filepath.Join(testDir, name)

	err = walk(name, f)
	assert.NotEqual(t, err, nil, "nil file")

	//if err = ioutil.WriteFile(name, content, os.ModePerm); err != nil {
	//	return
	//}
	//
	//err = FileSystem.AddBytes(fp, content)

}

//func TestUnpack2(t *testing.T) {
//	var err error
//	var name string
//	var f packr.File
//	name = "test"
//
//	err = walk(name, )
//	err = Unpack("wrongtest")
//	assert.NoError(t, err)
//}
