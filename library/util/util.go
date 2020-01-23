package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/facebook-spider/consts"
	"gitlab.azbit.cn/web/facebook-spider/models"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var logger = logging.MustGetLogger("library/util")

// json pretty print
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// get facebook mobile website
func GetMobilePostURL(url string) string {
	return strings.Replace(url, consts.FACEBOOK_DOMAIN, consts.BASIC_FACEBOOK_DOMAIN, 1)
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
	loc, _ := time.LoadLocation(consts.TIME_ZONE)
	t := time.Now().In(loc).Format("20060102")
	return t
}

// get current year
func GetCurrentYear() string {
	loc, _ := time.LoadLocation(consts.TIME_ZONE)
	return fmt.Sprintf("%v", time.Now().In(loc).Year())
}

// get current hour
func GetCurrentHour() int {
	loc, _ := time.LoadLocation(consts.TIME_ZONE)
	return time.Now().In(loc).Hour()
}

// get interval date
func GetOffsetDateStr(offset int) string {
	loc, _ := time.LoadLocation(consts.TIME_ZONE)
	nTime := time.Now().In(loc)
	yesTime := nTime.AddDate(0, 0, offset)
	return yesTime.Format("20060102")
}

// check file is exist
func CheckFileIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// request url
func RequestUrl(url string) ([]byte, error) {
	// Request the HTML page.
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
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
func GetPostFileName(d string) string {
	return fmt.Sprintf("%s%s", consts.POST_FILE_PREFIX, d)
}

// get comment file name
func GetCommentFileName(d string) string {
	return fmt.Sprintf("%s%s", consts.COMMENT_FILE_PREFIX, d)
}

// get article dir
func GetArticleDir(dir, lang, url string) (string, error) {
	//dir := strings.TrimRight(conf.Config.Spider.ArticleBaseDir, "/")
	name, err := GetOfficialAccountName(url)
	if err != nil {
		return "", err
	}
	return path.Join(dir, lang, name), nil
	//return fmt.Sprintf("%s%s/%s%s", dir, url, util.GetCurrentDate(), "posts")
}

// get date by facebook cell time
func GetDateByCellTime(cellTime string) string {
	// 1 sec - 59 secs, 1 min - 59 mins, 1 hr - 23 hrs
	// Sunday at 4:40 AM -> 20191112
	date := GetCurrentDate()
	cellTime = strings.Replace(cellTime, "小时", "hrs", -1)
	cellTime = strings.Replace(cellTime, "分钟", "mins", -1)
	loc, _ := time.LoadLocation(consts.TIME_ZONE)

	if strings.Contains(cellTime, "Today") {
		// Today at 3:22 AM
		date = GetCurrentDate()
	} else if strings.Contains(cellTime, "Yesterday") {
		// Yesterday at 12:28 PM
		date = GetOffsetDateStr(-1)
	} else if strings.Contains(cellTime, ",") {
		if strings.Contains(cellTime, "at") {
			// December 8, 2017 at 6:59 PM
			t, err := time.ParseInLocation("January 2, 2006 at 3:4 PM", cellTime, loc)
			if err == nil {
				date = fmt.Sprintf("%v", t.In(loc).Format("20060102"))
			} else {
				logger.Error("parse time err:", err.Error())
			}
		} else {
			// May 16, 2018
			t, err := time.ParseInLocation("Jan 2, 2006", cellTime, loc)
			if err == nil {
				date = fmt.Sprintf("%v", t.In(loc).Format("20060102"))
			} else {
				logger.Error("parse time err:", err.Error())
			}
		}
	} else if strings.Contains(cellTime, "at") {
		// October 29 at 11:33 PM
		layout := "2006 January 2 at 3:4 PM"
		if cellTime[3:4] == " " {
			// Oct 29 at 11:33 PM
			layout = "2006 Jan 2 at 3:4 PM"
		}
		if IsWeekFormat(cellTime) {
			// Monday at 9:25 AM
			arr := strings.Split(cellTime, " ")
			if len(arr) > 0 {
				offset := GetOffsetByToday(arr[0])
				date = GetOffsetDateStr(-1 * offset)
			}
		} else {
			tmp := GetCurrentYear() + " " + cellTime
			t, err := time.ParseInLocation(layout, tmp, loc)
			if err == nil {
				date = fmt.Sprintf("%v", t.In(loc).Format("20060102"))
			} else {
				fmt.Println("parse time err:", err.Error())
			}
		}
	} else if strings.Contains(cellTime, " hr") || strings.Contains(cellTime, " sec") || strings.Contains(cellTime, "min") {
		// 1 sec - 59 secs, 1 min - 59 mins, 1 hr - 23 hrs
		ch := GetCurrentHour()
		arr := strings.Split(cellTime, " hrs")
		if len(arr) >= 1 {
			ph, err := strconv.Atoi(arr[0])
			if err == nil && ph > ch {
				date = GetOffsetDateStr(-1)
			}
		}
	} else if len(cellTime) == 5 || len(cellTime) == 6 {
		// Aug 7
		tmp := cellTime + ", " + GetCurrentYear()
		t, err := time.ParseInLocation("Jan 2, 2006", tmp, loc)
		if err == nil {
			date = fmt.Sprintf("%v", t.In(loc).Format("20060102"))
		} else {
			logger.Error("parse time err:", err.Error())
		}
	}

	return date
}

// Is string contain number
func IsContainNumber(s string) bool {
	pattern := "\\d+"
	ret, _ := regexp.MatchString(pattern, s)
	return ret
}

func IsAllNumber(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return true
}

func GetCurrDay() string {
	loc, _ := time.LoadLocation(consts.TIME_ZONE)
	nTime := time.Now().In(loc)
	return nTime.Weekday().String()
}

func IsWeekFormat(s string) bool {
	a := strings.Split(s, " ")
	if len(a) <= 0 {
		return false
	}
	var longDayNames = []string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
	for _, day := range longDayNames {
		if a[0] == day {
			return true
		}
	}
	return false
}

// s2-s1
func GetOffsetByToday(s1 string) int {
	s2 := GetCurrDay()
	var longDayNames = []string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
	s1Index := 0
	s2Index := 0
	for idx, day := range longDayNames {
		if day == s1 {
			s1Index = idx
		}
		if day == s2 {
			s2Index = idx
		}
	}
	idx := s2Index - s1Index
	if idx < 0 {
		idx += 7
	}
	return idx
}
