package test

import (
	"gitlab.azbit.cn/web/facebook-spider/library/util"
	"testing"
)

const EXPECTED_CURR_DATE = "20191112"
const EXPECTED_YEST_DATE = "20191111"
const EXPECTED_OCT_DATE = "20191009"
const EXPECTED_EIGHT_DATE = "20171208"
const EXPECTED_SUNDY_DATE = "20191110"

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
		{"December 8, 2017 at 6:59 PM", EXPECTED_EIGHT_DATE},
		{"Today at 3:22 AM", EXPECTED_CURR_DATE},
		{"Sunday at 4:40 AM", EXPECTED_SUNDY_DATE},
	}
	for _, tt := range timeTests {
		actual := util.GetDateByCellTime(tt.in)
		if actual != tt.expected {
			t.Errorf("cellTime:%s, expected:%s, actual:%s", tt.in, tt.expected, actual)
		}
	}
}
