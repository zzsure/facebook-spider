package test

import (
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"testing"
)

const EXPECTED_CURR_DATE = "20191031"
const EXPECTED_YEST_DATE = "20191030"
const EXPECTED_OCT_DATE = "20191009"

func TestParseTime(t *testing.T) {
	var timeTests = []struct {
		in       string
		expected string
	}{
		{"1 sec", EXPECTED_CURR_DATE},
		{"1 min", EXPECTED_CURR_DATE},
		{"1 min", EXPECTED_CURR_DATE},
		{"2 mins", EXPECTED_CURR_DATE},
		{"1 hr", EXPECTED_CURR_DATE},
		{"2 hrs", EXPECTED_CURR_DATE},
		{"Yesterday at 12:28 PM", EXPECTED_YEST_DATE},
		{"October 9 at 6:36 AM", EXPECTED_OCT_DATE},
	}
	for _, tt := range timeTests {
		actual := util.GetDateByCellTime(tt.in)
		if actual != tt.expected {
			t.Errorf("cellTime:%s, expected:%s, actual:%s", tt.in, tt.expected, actual)
		}
	}
}
