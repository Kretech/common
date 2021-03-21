package carbon

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Kretech/xgo/test"
)

func TestUnixOf(t *testing.T) {
	In(Shanghai)
	test.AssertEqual(t, UnixOf(0, 0).Format("Y-m-d H:i:s"), "1970-01-01 08:00:00")

	In(time.UTC)
	test.AssertEqual(t, UnixOf(0, 0).In(Shanghai).Format("Y-m-d H:i:s"), "1970-01-01 08:00:00")
	test.AssertEqual(t, UnixOf(0, 0).Format("Y-m-d H:i:s"), "1970-01-01 00:00:00")

	In(time.UTC)
	test.AssertEqual(t, UnixOf(0, 0).Time().String(), "1970-01-01 00:00:00 +0000 UTC")
	test.AssertEqual(t, UnixOf(0, 0).Time().Unix(), 0)
	test.AssertEqual(t, UnixOf(0, 0).Format("Y-m-d H:i:s"), "1970-01-01 00:00:00")
	test.AssertEqual(t, UnixOf(0, 0).In(Shanghai).Time(), "1970-01-01 08:00:00 +0800 CST")
	test.AssertEqual(t, UnixOf(0, 0).In(Shanghai).Time().Unix(), 0)
	test.AssertEqual(t, UnixOf(0, 0).In(Shanghai).Format("Y-m-d H:i:s"), "1970-01-01 08:00:00")
}

func TestParse(t *testing.T) {
	parse("2012-1-1 01:02:03")
	parse("+1 day")
}

func TestCarbon_Format(t *testing.T) {
	In(Shanghai)
	test.AssertEqual(t, UnixOf(1524379525, 0).Format(`Y-m-d`), `2018-04-22`)
	test.AssertEqual(t, UnixOf(1524379525, 0).Format(`Y-n-j`), `2018-4-22`)
	test.AssertEqual(t, UnixOf(1523031400, 0).Format(`Y-n-j`), `2018-4-7`)
	test.AssertEqual(t, UnixOf(1524379525, 0).Format(`Y-m-d H:i:s`), `2018-04-22 14:45:25`)
}

func TestCarbon_In(t *testing.T) {
	t1 := Now()
	test.AssertEqual(t, t1.Time().Location(), time.Local)

	t2 := t1.In(Shanghai)
	// t1未修改
	test.AssertEqual(t, t1.Time().Location(), time.Local)
	// t2修改成功
	test.AssertEqual(t, t2.Time().Location(), Shanghai)
}

func TestCarbon_Sub(t *testing.T) {
	t1 := ParseUTC("Y-m-d H:i:s", "2018-01-02 09:00:00")
	t2 := t1.Sub(time.Hour)
	test.AssertEqual(t, t2.Format("Y-m-d H:i:s"), "2018-01-02 08:00:00")
}

func TestCarbon_MarshalJSO2N(t *testing.T) {

	t0 := ParseUTC("Y-m-d H:i:s", "1970-01-01 00:00:00")
	js, err := json.Marshal(t0)
	test.AssertNil(t, err)

	test.AssertEqual(t, string(js), `0`)

	t1 := ParseUTC("Y-m-d H:i:s", "2018-01-02 01:00:00")
	js, err = json.Marshal(t1)
	test.AssertNil(t, err)

	test.AssertEqual(t, string(js), `1514854800`)
}

func TestCarbon_MarshalJSON(t *testing.T) {
	type args struct {
		*time.Location
		str     string
		unixStr string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{``, args{UTC, `1970-01-01 00:00:00`, `0`}},
		{``, args{Shanghai, `1970-01-01 08:00:00`, `0`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ParseIn(`Y-m-d H:i:s`, tt.args.str, tt.args.Location)
			unixBytes, err := c.MarshalJSON()

			test.AssertNil(t, err)
			test.AssertEqual(t, string(unixBytes), tt.args.unixStr)

			unix, _ := strconv.Atoi(tt.args.unixStr)
			unix += 3600
			unixStr := fmt.Sprint(unix)
			var c2 Carbon
			err = json.Unmarshal([]byte(unixStr), &c2)
			test.AssertNil(t, err)
			test.AssertEqual(t, c.Add(time.Hour).String(), c2.String())
		})
	}
}
