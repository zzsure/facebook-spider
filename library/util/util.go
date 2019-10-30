package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// 将录入的url的格式转化成文章链接的地址
func GetOfficialAccountURL(url string) (string, error) {
	n, err := GetOfficialAccountName(url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://facebook.com/pg/%s/posts/", n), nil
}

// 获取公众号名称
func GetOfficialAccountName(url string) (string, error) {
	arr := strings.Split(url, "/")
	if len(arr) != 5 {
		return "", errors.New("url format is not correct. [http://facebook.com/xxx/]")
	}
	if len(arr[3]) == 0 {
		return "", errors.New("url not have official account name")
	}
	return arr[3], nil
}

// 存储string到文件中
func SaveStringToFile(dir, name, data string) error {
	if dir == "" || name == "" || data == "" {
		return errors.New("dir|name|data should not empty string")
	}

	var f *os.File
	var err error

	p := path.Join(dir, name)
	if CheckFileIsExist(p) == false {
		err = os.MkdirAll(dir, 0711)
		f, err = os.Create(p)
	} else {
		f, err = os.OpenFile(p, os.O_WRONLY|os.O_TRUNC, 0666)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(data)

	return err
}

func ReadStringFromFile(path string) (string, error) {
	if path == "" {
		return "", errors.New("path should not empty string")
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 获取当前的日期
func GetCurrentDate() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc).Format("20060102")
	return t
}

// 检查文件是否存在
func CheckFileIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
