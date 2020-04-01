package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
func WriteFile4byte(filename string, data *[]byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriterSize(f, 4096)
	if _, err := writer.Write(*data); err == nil {
		if err := writer.Flush(); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// 获取工作目录
func GetWorkDir() string {
	dir, _ := os.Getwd()
	// fplib.Debug("Work Dir：", dir)
	return dir
}

// 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//获取指定目录下的所有文件和目录
func ListDir(dirPth string) (files []string, files1 []string, err error) {
	//fmt.Println(dirPth)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}
	PthSep := string(os.PathSeparator)
	//    suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {

		if fi.IsDir() { // 忽略目录
			files1 = append(files1, dirPth+PthSep+fi.Name())
			ListDir(dirPth + PthSep + fi.Name())
		} else {
			//fmt.Println("s")
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, files1, nil
}
