package vpn_conf

import (
	"strconv"
	"strings"
)

func NextAddress4(ip string) string {
	digits := addrDigits(ip)
	digits = incrementIP4(digits)
	return ipString(digits)
}

func NextAddress6(ip string) string {
	digits := addrDigits6(ip)
	digits = incrementIP6(digits)
	return ipString6(digits)
}

func ipString(ip []int) string {
	strs := make([]string, len(ip))
	for i, a := range ip {
		strs[i] = strconv.Itoa(a)
	}
	return strings.Join(strs, ".")
}

func ipString6(ip []int) string {
	strs := make([]string, len(ip))
	for i, a := range ip {
		strs[i] = strconv.FormatInt(int64(a), 16)
	}
	return strings.Join(strs, ":")
}

func addrDigits(ip string) []int {
	parts := strings.Split(ip, ".")
	digits := make([]int, len(parts))
	for i, a := range parts {
		digits[i], _ = strconv.Atoi(a)
	}
	return digits
}

func addrDigits6(ip string) []int {
	parts := strings.Split(ip, ":")
	digits := make([]int, len(parts))
	for i, a := range parts {
		tmp, _ := strconv.ParseInt(a, 16, 32)
		digits[i] = int(tmp)
	}
	return digits
}

func incrementIP4(addr []int) []int {
	return incrementIP(addr, 256)
}

func incrementIP6(addr []int) []int {
	return incrementIP(addr, 65536)
}

func incrementIP(addr []int, base int) []int {
	addr[len(addr)-1] += 1
	for x := len(addr) - 1; x > 0; x-- {
		if addr[x] > base-1 {
			addr[x-1] += 1
			addr[x] = 0
		}
	}
	return addr
}
