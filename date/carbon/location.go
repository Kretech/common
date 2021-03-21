package carbon

import "time"

var (
	UTC         = time.UTC
	Shanghai, _ = time.LoadLocation("Asia/Shanghai")
)
