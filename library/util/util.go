package util

import (
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"gitlab.azbit.cn/web/facebook-spider/library/util/net"
)

func GenServerUUID() string {
	ip, mac := net.NewLAN().NetInfo()
	return fmt.Sprintf("%s-%s", ip, mac)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func CoverGBKToUTF8(src string) string {
	return mahonia.NewDecoder("gbk").ConvertString(src)
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}