package wurfl_test

import (
	"net/http"
	"testing"

	wurfl "github.com/WURFL/golang-wurfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Benchmark_GetStaticCap_FallbackCacheDefault_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	require.NotNil(b, wengine)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName, _ := device.GetStaticCap("model_name")
	if modelName == "" {
		b.Error("capability model_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDefault {
		b.Error("fallback cache mode should be DEFAULT")
	}

	defer device.Destroy()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheDefault_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	require.NotNil(b, wengine)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv, _ := device.GetStaticCap("is_smarttv")
	assert.NotEmpty(b, isSmarttv)

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")
	assert.NoError(b, err)

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	assert.NoError(b, err)

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDefault {
		b.Error("fallback cache mode should be DEFAULT")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheLimited_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName, _ := device.GetStaticCap("model_name")
	if modelName == "" {
		b.Error("capability model_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheLimited {
		b.Error("fallback cache mode should be LIMITED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheLimited_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv, _ := device.GetStaticCap("is_smarttv")
	if isSmarttv == "" {
		b.Error("capability is_smarttv must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheLimited {
		b.Error("fallback cache mode should be LIMITED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheDisabled_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName, _ := device.GetStaticCap("model_name")
	if modelName == "" {
		b.Error("capability model_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDisabled {
		b.Error("fallback cache mode should be DISABLED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheDisabled_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv, _ := device.GetStaticCap("is_smarttv")
	if isSmarttv == "" {
		b.Error("capability is_smarttv must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDisabled {
		b.Error("fallback cache mode should be DISABLED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheDisabled_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName, _ := device.GetStaticCap("marketing_name")
	if mktName == "" {
		b.Error("capability marketing_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("samsung_sm_g990b_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDisabled {
		b.Error("fallback cache mode should be DISABLED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("marketing_name")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheDefault_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName, _ := device.GetStaticCap("marketing_name")
	if mktName == "" {
		b.Error("capability marketing_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("samsung_sm_g990b_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDefault {
		b.Error("fallback cache mode should be DEFAULT")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("marketing_name")
	}

	b.StopTimer()
}

func Benchmark_GetStaticCap_FallbackCacheLimited_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName, _ := device.GetStaticCap("marketing_name")
	if mktName == "" {
		b.Error("capability marketing_name must have a value")
	}

	device.Destroy() // destroy and recreate to have the capability cache clean
	device, err = wengine.LookupDeviceID("samsung_sm_g990b_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	fallbackMode, err := wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
	if err != nil {
		b.Error("error from GetAttr should be nil")
	}

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheLimited {
		b.Error("fallback cache mode should be LIMITED")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("marketing_name")
	}

	b.StopTimer()
}

// Benchmark_GetCapability : test time in single const get_capability
func Benchmark_GetStaticCap(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetStaticCap("device_os")
	}

	b.StopTimer()
}

func Benchmark_GetCapabilityAsInt(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetCapabilityAsInt("resolution_width")
	}

	b.StopTimer()
}

func Benchmark_GetDeviceID(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		b.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		b.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetDeviceID()
	}

	b.StopTimer()
}
func Benchmark_LookupDeviceID(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupDeviceID("google_pixel_5_ver1")
		device.Destroy()
	}

	b.StopTimer()
}

func Benchmark_LookupUserAgent_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupUserAgent(UserAgent)
		device.Destroy()
	}

	b.StopTimer()
}

func Benchmark_LookupUserAgent_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupUserAgent(UserAgent)
		device.Destroy()
	}

	b.StopTimer()
}

// func Benchmark_LookupRequest_Cache
func Benchmark_LookupRequest_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Sec-CH-UA", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Add("Sec-CH-UA-Full-Version", "90.0.4430.91")
	req.Header.Add("Sec-CH-UA-Platform", "Android")
	req.Header.Add("Sec-CH-UA-Platform-Version", "11")
	req.Header.Add("Sec-CH-UA-Model", "SM-M315F")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupRequest(req)
		device.Destroy()
	}

	b.StopTimer()
}

// func Benchmark_LookupRequest_NoCache
func Benchmark_LookupRequest_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Sec-CH-UA", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Add("Sec-CH-UA-Full-Version", "90.0.4430.91")
	req.Header.Add("Sec-CH-UA-Platform", "Android")
	req.Header.Add("Sec-CH-UA-Platform-Version", "11")
	req.Header.Add("Sec-CH-UA-Model", "SM-M315F")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupRequest(req)
		device.Destroy()
	}

	b.StopTimer()
}

// GetAllCaps benchmarks
func Benchmark_GetAllCaps(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		wengine.GetAllCaps()
	}

	b.StopTimer()
}

// GetAllVCaps benchmarks
func Benchmark_GetAllVCaps(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		wengine.GetAllVCaps()
	}

	b.StopTimer()
}

// GetAllDeviceIds benchmark
func Benchmark_GetAllDeviceIds(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		wengine.GetAllDeviceIds()
	}

	b.StopTimer()
}

// benchmark CString()/Free() couple compared to accessing a map to get the desired CString
// to understand whether a caps names CString cache might be interesting or not
func Benchmark_CStringCFree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wurfl.GoStringToCStringAndFree("capability")
	}
}

func Benchmark_MapAccess(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wengine.GoStringToCStringUsingMap("capability")
		// if i%10000000 == 0 {
		// 	fmt.Println(CCap)
		// }
	}
	b.StopTimer()
}

func BenchmarkGetStaticCap(b *testing.B) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(nil)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := device.GetStaticCap("mobile_browser_version")
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
func BenchmarkGetVirtualCap(b *testing.B) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(nil)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := device.GetVirtualCap("advertised_device_os")
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
