package tools

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func SaveFileWithMd5Name(data io.Reader, dir string, ext string) (string, error) {
	tmpF, err := os.CreateTemp(dir, "tmp_*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpF.Name())
	defer tmpF.Close()

	r := io.TeeReader(data, tmpF)
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	newPath := filepath.Join(dir, fmt.Sprintf("%x%s", h.Sum(nil), ext))
	_, err = os.Stat(newPath)
	if err == nil {
		return newPath, nil
	}
	err = os.Rename(tmpF.Name(), newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func Move(src, dst string) error {
	_, err := os.Stat(dst)
	if err == nil {
		return os.ErrExist
	}
	if os.IsExist(err) {
		return err
	}
	err = os.Rename(src, dst)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(src, dst)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Open(source)")
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return errors.Wrap(err, "Create(destination)")
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return errors.Wrap(err, "Copy")
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	os.Remove(source)
	return nil
}

func Cp(src, dst string) error {
	_, err := os.Stat(dst)
	if err == nil {
		return os.ErrExist
	}
	if os.IsExist(err) {
		return err
	}
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	_, err = io.Copy(dstF, srcF)
	return err
}
