package date

import (
	"strings"
	"time"
)

var timeOffset time.Duration = 0

func TimeOffset(offset time.Duration) {
	timeOffset = offset
}

// ToGoFormat 把 Y-m-d 转成 Go 格式
// cache ?
func ToGoFormat(format string) string {
	rule := []string{
		`Y`, `2006`, //  年
		`y`, `06`,

		`m`, `01`, //  月
		`n`, `1`,

		`d`, `02`, //  日
		`j`, `2`,

		`H`, `15`, //  时
		`h`, `03`,
		`g`, `3`,
		`G`, `15`, // 应该木有前导零

		`i`, `04`, //  分
		`s`, `05`, //  秒

		`D`, `Mon`, //  周
		`N`, `1`,
	}

	specs := func(expr string) string {
		switch expr {
		case `S`:
			return `st/nd/rd/th`
		case `z`:
			return `The day of the year`
		case `t`:
			return `Number of days in the given month`
		case `L`:
			return `Whether it's a leap year`
		case `a`:
			return `am or pm`
		case `A`:
			return `AM or PM`
		default:
			return expr
		}
	}

	size := len(rule)
	for i := 0; i < size; i += 2 {
		format = strings.Replace(format, rule[i], rule[i+1], -1)
	}

	format = specs(format)

	return format
}

func LocalFormat(format string, timestamp ...int64) string {
	var seconds int64
	if len(timestamp) > 0 {
		seconds = timestamp[0]
	} else {
		seconds = time.Now().Local().Add(timeOffset).Unix()
	}

	format = ToGoFormat(format)

	return time.Unix(seconds, 0).Local().Format(format)
}

// @see http://php.net/manual/en/function.strtotime.php#example-2803
func StrToTime(expr string) int64 {
	t := time.Now().Local()
	//pattern := `(\d+|[(last)|(next)]) [(year)|(month)|(day)|(hour)|(minute)|(second)|(Monday)]`

	for i := 1; i < 10; i++ {

	}
	return t.Unix()
}
