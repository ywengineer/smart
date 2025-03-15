package utility

import (
	"regexp"
	"strings"
)

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
var mailPPattern = regexp.MustCompile(`^([a-z0-9_.-]+)@([\da-z.-]+)\.([a-z.]{2,6})$`)

func ValidMail(m string) bool {
	if len(m) == 0 {
		return false
	}
	return mailPPattern.MatchString(m)
}

func ValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	i := strings.LastIndex(ipAddress, ":")
	ipAddress = ipAddress[:i] //remove port

	return ip4Reg.MatchString(ipAddress)
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxInt[T int | int32 | int64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func MinInt[T int | int32 | int64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func StringWithNoSpace(str string) bool {
	return len(strings.TrimSpace(str)) > 0
}
