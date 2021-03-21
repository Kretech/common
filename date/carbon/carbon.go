package carbon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Kretech/xgo/date"
	"github.com/pkg/errors"
)

// Carbon
type Carbon struct {
	t time.Time
}

// 假定多数情况下，一个进程给一个location就够了
var currentLoc = time.Local

// 设置全局时区
func In(loc *time.Location) {
	currentLoc = loc
	time.Local = loc
}

func Now() *Carbon {
	return &Carbon{time.Now()}
}

func Today() *Carbon {
	return &Carbon{time.Time{}.Truncate(24 * Hour)}
}

func Yestoday() *Carbon {
	return Today().Sub(Day)
}

func Tomorrow() *Carbon {
	return Today().Add(Day)
}

func UnixOf(sec int64, nsec int64) *Carbon {
	t := time.Unix(sec, nsec)
	return &Carbon{t}
}

func TimeOf(t time.Time) *Carbon {
	return &Carbon{t}
}

// 解析时间文本
func strOf(str string) *Carbon {
	return parse(str)
}

// todo
// alias to StrOf
func parse(str string) *Carbon {
	t := time.Now().In(currentLoc)
	return &Carbon{t}
}

// wrap of time.Parse but with the carbon layout
// 使用 Go time.Time 解析时间文本
func ParseUTC(layout, value string) *Carbon {
	return ParseIn(layout, value, UTC)
}

func ParseLocal(layout, value string) *Carbon {
	return ParseIn(layout, value, time.Local)
}

func ParseIn(layout, value string, loc *time.Location) *Carbon {
	t, _ := time.ParseInLocation(date.ToGoFormat(layout), value, loc)
	return TimeOf(t)
}

func carbonOf(c *Carbon) *Carbon {
	return TimeOf(c.t)
}

func (c *Carbon) Clone() *Carbon {
	return carbonOf(c)
}

func (c *Carbon) Time() time.Time {
	return c.t
}

// 设置时区
func (c Carbon) In(loc *time.Location) *Carbon {
	c.t = c.t.In(loc)
	return &c
}

func (c Carbon) Add(d time.Duration) *Carbon {
	c.t = c.t.Add(d)
	return &c
}

// 减去另一个时间：计算两个时间的差
func (c *Carbon) Diff(t time.Time) time.Duration {
	return c.Time().Sub(t)
}

// 减去一段时间
func (c Carbon) Sub(d time.Duration) *Carbon {
	c.t = c.t.Add(-d)
	return &c
}

// Format 格式化代码
// c.Format("Y-m-d H:i:s")
// @see http://php.net/manual/zh/function.date.php#refsect1-function.date-parameters
func (c *Carbon) Format(layout string) string {
	layout = date.ToGoFormat(layout)
	return c.t.Format(layout)
}

func (this *Carbon) Unix() int64 {
	return this.t.Unix()
}

func (this *Carbon) UnixNano() int64 {
	return this.t.UnixNano()
}

func (this *Carbon) Truncate(d time.Duration) *Carbon {
	return TimeOf(this.t.Truncate(d))
}

func (c *Carbon) String() string {
	return c.DatetimeString()
}

func (c *Carbon) DateString() string {
	return c.Format("Y-m-d")
}

func (c *Carbon) DatetimeString() string {
	return c.Format("Y-m-d H:i:s")
}

func (c *Carbon) UnmarshalJSON(b []byte) (err error) {
	unix, err := strconv.Atoi(string(b))
	if err != nil {
		err = errors.WithMessage(err, `Carbon UnmarshalJSON`)
		return
	}

	*c = *UnixOf(int64(unix), 0)

	return
}

func (c *Carbon) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(c.Unix())), nil
}
