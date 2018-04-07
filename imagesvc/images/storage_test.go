package images

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"testing/iotest"

	"github.com/powerman/check"
	"github.com/powerman/must"
)

const (
	testFileName = "file"
)

var (
	testFile = []byte{0, 1, 2, 3, 4}
)

func init() {
	must.AbortIf = must.PanicIf
}

func TestSaveImage(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)
	st := NewStorage(tmpDir)
	err := st.Save(testFileName, bytes.NewBuffer(testFile))
	t.Nil(err)

	bs := must.ReadFile(tmpDir + "/" + testFileName)
	t.Zero(bytes.Compare(bs, testFile))
}

func TestSaveImageDuplicate(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	err := st.Save(testFileName, bytes.NewBuffer(testFile))
	t.Nil(err)

	err = st.Save(testFileName, bytes.NewBuffer(append(testFile, testFile...)))
	t.EQ(err, errFileExist)

	bs := must.ReadFile(tmpDir + "/" + testFileName)
	t.Zero(bytes.Compare(bs, testFile))
}

func TestSaveImageReaderErr(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	err := st.Save(testFileName, iotest.TimeoutReader(bytes.NewBuffer(testFile)))
	t.NotNil(err)

	_, err = ioutil.ReadFile(tmpDir + "/" + testFileName)
	t.True(os.IsNotExist(err))
}

func TestSaveImageEmptyFilename(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	err := st.Save("", bytes.NewBuffer(testFile))
	t.EQ(err, errEmptyFilename)
}

func TestSaveImageTryEscapeFolderFilename(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	fileName := "escape"
	escapeFileName := "../" + fileName
	err := st.Save(escapeFileName, bytes.NewBuffer(testFile))
	t.Nil(err)

	//File created inside our directory
	_, err = ioutil.ReadFile(tmpDir + "/" + fileName)
	t.Nil(err)
}

func TestSaveImageDotFilename(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	err := st.Save(".", bytes.NewBuffer(testFile))
	t.EQ(err, errEmptyFilename)

	_, err = ioutil.ReadFile(tmpDir + "/" + testFileName)
	t.True(os.IsNotExist(err))
}

func TestDeleteImage(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	must.NoErr(st.Save(testFileName, bytes.NewBuffer(testFile)))

	err := st.Delete(testFileName)
	t.Nil(err)

	_, err = ioutil.ReadFile(tmpDir + "/" + testFileName)
	t.True(os.IsNotExist(err))
}

func TestDeleteImageNoFileNoErr(tt *testing.T) {
	t := check.T{tt}
	tmpDir := must.TempDir("", "")
	defer os.RemoveAll(tmpDir)

	st := NewStorage(tmpDir)
	err := st.Delete(testFileName)
	t.Nil(err)

	_, err = ioutil.ReadFile(tmpDir + "/" + testFileName)
	t.True(os.IsNotExist(err))
}
