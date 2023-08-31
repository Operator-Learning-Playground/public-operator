package helpers

import (
	"io"
	"os"
)

func MustLoadFile(path string) []byte {
	b, err := LoadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return true

}

// MustFileExists 判断文件是否必须存在
func MustFileExists(file string) {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		panic("file name:" + file + "not found")
	}

}

// LoadFile 加载指定目录的文件, 全部取出内容
func LoadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return b, err
}
