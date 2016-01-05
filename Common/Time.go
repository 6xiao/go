package Common

import (
	"time"
)

// format a time.Time to string as 2006-01-02 15:04:05.999
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999")
}

// format time.Now() use FormatTime
func FormatNow() string {
	return FormatTime(time.Now())
}

// parse a string as "2006-01-02 15:04:05.999" to time.Time
func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05.999", s, time.Local)
}

// format a time.Time to number as 20060102150405999
func NumberTime(t time.Time) uint64 {
	y, m, d := t.Date()
	h, M, s := t.Clock()
	ms := t.Nanosecond() / 1000000
	return uint64(ms+s*1000+M*100000+h*10000000+d*1000000000) +
		uint64(m)*100000000000 + uint64(y)*10000000000000
}

// format time.Now() use NumberTime
func NumberNow() uint64 {
	return NumberTime(time.Now())
}

// parse a uint64 as 20060102150405999 to time.Time
func ParseNumber(t uint64) (time.Time, error) {
	ns := int((t % 1000) * 1000000)
	t /= 1000
	s := int(t % 100)
	t /= 100
	M := int(t % 100)
	t /= 100
	h := int(t % 100)
	t /= 100
	d := int(t % 100)
	t /= 100
	m := time.Month(t % 100)
	y := int(t / 100)
	
	return time.Date(y, m, d, h, M, s, ns, time.Local), nil
}
