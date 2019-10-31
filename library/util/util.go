package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// json pretty print
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// get official account post url by url
func GetOfficialAccountPostURL(url string) (string, error) {
	n, err := GetOfficialAccountName(url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(consts.POST_URL_FORMAT, n), nil
}

// get official account name
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

// save string to file
// TODO: filter repeat strings
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

// read string from file
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

// get current date string
func GetCurrentDate() string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc).Format("20060102")
	return t
}

// check file is exist
func CheckFileIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// read url and language information from csv file
func ReadUrlsFromCsv(path string) ([]*models.FileData, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fds := []*models.FileData{}
	if err := gocsv.UnmarshalFile(file, &fds); err != nil {
		return nil, err
	}
	return fds, nil
}

// get post file name
func GetPostFileName() string {
	// TODO：日期改成爬取的
	return fmt.Sprintf("%s%s", consts.POST_FILE_PREFIX, GetCurrentDate())
}

func GetCommentsFileName() string {
	// TODO: 日期改成爬取的
	return fmt.Sprintf("%s%s", consts.COMMENT_FILE_PREFIX, GetCurrentDate())
}

// get post dir
func GetPostsDir(dir, url string) (string, error) {
	//dir := strings.TrimRight(conf.Config.Spider.ArticleBaseDir, "/")
	name, err := GetOfficialAccountName(url)
	if err != nil {
		return "", err
	}
	return path.Join(dir, name), nil
	//return fmt.Sprintf("%s%s/%s%s", dir, url, util.GetCurrentDate(), "posts")
}
