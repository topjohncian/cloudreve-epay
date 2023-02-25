package appentry

import (
	"bufio"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func Eject(templateFS embed.FS) {
	walk := func(relPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Errorf("Failed to read info of %q: %s, skipping...", relPath, err)
		}

		if !d.IsDir() {
			// 写入文件
			out, err := CreateNestedFile(filepath.Join(lo.Must(os.Getwd()), "custom", relPath))
			if err != nil {
				return errors.Errorf("Failed to create file %q: %s, skipping...", relPath, err)
			}

			defer out.Close()

			logrus.Infof("导出 %s...", relPath)
			obj, _ := templateFS.Open(relPath)
			if _, err := io.Copy(out, bufio.NewReader(obj)); err != nil {
				return errors.Errorf("Cannot write file %q: %s, skipping...", relPath, err)
			}
		}
		return nil
	}

	err := fs.WalkDir(templateFS, ".", walk)
	if err != nil {
		logrus.WithError(err).Panic("导出模板文件失败")
		return
	}
	logrus.Info("成功导出模板文件")
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreateNestedFile 给定path创建文件，如果目录不存在就递归创建
func CreateNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if !Exists(basePath) {
		err := os.MkdirAll(basePath, 0700)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(path)
}

// IsEmpty 返回给定目录是否为空目录
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
