package monitor

import "testing"

func TestAmbaris(t *testing.T) {
	wants := map[string]Ambaris{
		"baidu": Ambaris{
			IP: "www.baidu.com",
		},
		"amnaris": Ambaris{
			IP:   "192.168.2.144",
			Port: "8080",
			Path: "/",
		},
	}
	for k, v := range wants {
		t.Log(k, v.Proxy())
	}
}
