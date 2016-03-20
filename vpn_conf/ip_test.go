package vpn_conf

import "testing"

func TestDigitWrap(t *testing.T) {
	startIp := "10.8.0.0"
	nextIp := startIp
	for i := 0; i < 256*2+1; i++ {
		nextIp = NextAddress4(nextIp)
	}
	shouldBe := "10.8.2.1"
	if nextIp != shouldBe {
		t.Errorf("Expexted %v got %v", shouldBe, nextIp)
	}
}

func TestIP6(t *testing.T) {
	startIp := "2001:db8:1:2:0:0:0:2"
	nextIp := startIp
	for i := 0; i < 65536; i++ {
		nextIp = NextAddress6(nextIp)
		// log.Println(nextIp)
	}
	shouldBe := "2001:db8:1:2:0:0:1:2"
	if nextIp != shouldBe {
		t.Errorf("Expexted %v got %v", shouldBe, nextIp)
	}
}
