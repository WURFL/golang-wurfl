package wurfl_test

import (
	"testing"

	wurfl "github.com/WURFL/golang-wurfl"
)

func BenchmarkDeviceAll(b *testing.B) {

	caps := []string{
		"mobile_browser_version",
		"pointing_method",
		"is_tablet",
	}

	wengine, err := wurfl.Create("/usr/share/wurfl/wurfl.zip", nil, caps, -1, wurfl.WurflCacheProviderLru, "100000")
	if err != nil {
		b.Fatal(err)
	}
	defer wengine.Destroy()

	for i := 0; i < b.N; i++ {
		wengine.GetAllDeviceIds()
	}
}
