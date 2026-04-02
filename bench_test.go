package wurfl_test

import (
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"
	"unsafe"

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

// Benchmark_LookupWithImportantHeaderMap_Cache
func Benchmark_LookupWithImportantHeaderMap_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["Sec-CH-UA"] = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	IHMap["Sec-CH-UA-Full-Version"] = "90.0.4430.91"
	IHMap["Sec-CH-UA-Platform"] = "Android"
	IHMap["Sec-CH-UA-Platform-Version"] = "11"
	IHMap["Sec-CH-UA-Model"] = "SM-M315F"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupWithImportantHeaderMap(IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupWithImportantHeaderMap_AllSecChUa_Cache
func Benchmark_LookupWithImportantHeaderMap_AllSecChUa_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	IHMap := make(map[string]string)
	IHMap["Accept-Encoding"] = "gzip, br"
	IHMap["X-Requested-With"] = "com.instagram.android"
	IHMap["User-Agent"] = UserAgent
	IHMap["Sec-CH-UA"] = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	IHMap["Sec-CH-UA-Full-Version"] = "90.0.4430.91"
	IHMap["Sec-CH-UA-Platform"] = "Android"
	IHMap["Sec-CH-UA-Platform-Version"] = "11"
	IHMap["Sec-CH-UA-Model"] = "SM-M315F"
	IHMap["Sec-CH-UA-Mobile"] = "?1"
	IHMap["Sec-CH-UA-Arch"] = "x86"
	IHMap["Sec-Ch-Ua-Full-Version-List"] = "\"Chromium\";v=\"146.0.7680.157\", \"Not-A.Brand\";v=\"24.0.0.0\", \"Android WebView\";v=\"146.0.7680.157\""

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupWithImportantHeaderMap(IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupWithImportantHeaderMap_NoCache
func Benchmark_LookupWithImportantHeaderMap_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["Sec-CH-UA"] = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	IHMap["Sec-CH-UA-Full-Version"] = "90.0.4430.91"
	IHMap["Sec-CH-UA-Platform"] = "Android"
	IHMap["Sec-CH-UA-Platform-Version"] = "11"
	IHMap["Sec-CH-UA-Model"] = "SM-M315F"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupWithImportantHeaderMap(IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupRequest_Long_Cache
func Benchmark_LookupRequest_Long_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 14; XQ-CT54 Build/64.2.A.2.258; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.157 Mobile Safari/537.36 Instagram 422.0.0.44.64 Android (34/14; 420dpi; 1096x2560; Sony; XQ-CT54; XQ-CT54; qcom; en_GB; 916494010; IABMV/1)"
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, br")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("Cdn-Loop", "cloudflare; loops=1")
	req.Header.Add("Cf-Connecting-Ip", "217.216.97.90")
	req.Header.Add("Cf-Ew-Via", "15")
	req.Header.Add("Cf-Ipcountry", "DE")
	req.Header.Add("Cf-Ray", "9e3709414e65a043-FRA")
	req.Header.Add("Cf-Visitor", "{\"scheme\":\"https\"}")
	req.Header.Add("Cf-Worker", "html-load.cc")
	req.Header.Add("Forwarded", "for=217.216.97.90")
	req.Header.Add("Host", "prebid.wurflcloud.com")
	req.Header.Add("Referer", "https://vsco.co/")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"146\", \"Not-A.Brand\";v=\"24\", \"Android WebView\";v=\"146\"")
	req.Header.Add("Sec-Ch-Ua-Full-Version", "\"146.0.7680.157\"")
	req.Header.Add("Sec-Ch-Ua-Full-Version-List", "\"Chromium\";v=\"146.0.7680.157\", \"Not-A.Brand\";v=\"24.0.0.0\", \"Android WebView\";v=\"146.0.7680.157\"")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Add("Sec-Ch-Ua-Model", "\"XQ-CT54\"")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Android\"")
	req.Header.Add("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")
	req.Header.Add("Sec-Fetch-Dest", "script")
	req.Header.Add("Sec-Fetch-Mode", "no-cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Storage-Access", "active")
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-Amzn-Trace-Id", "Root=1-69c7d9dc-618fdb3e4d9456fc38f66cab")
	req.Header.Add("X-As-Script-Type", "ESSENTIAL")
	req.Header.Add("X-Device-Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("X-Device-Ip", "217.216.97.90")
	req.Header.Add("X-Device-Referer", "https://vsco.co/")
	req.Header.Add("X-Device-User-Agent", UserAgent)
	req.Header.Add("X-Forwarded-For", "217.216.97.90, 172.71.144.58")
	req.Header.Add("X-Requested-With", "com.instagram.android")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupRequest(req)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupRequest_Long_NoCache
func Benchmark_LookupRequest_Long_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 14; XQ-CT54 Build/64.2.A.2.258; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.157 Mobile Safari/537.36 Instagram 422.0.0.44.64 Android (34/14; 420dpi; 1096x2560; Sony; XQ-CT54; XQ-CT54; qcom; en_GB; 916494010; IABMV/1)"
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, br")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("Cdn-Loop", "cloudflare; loops=1")
	req.Header.Add("Cf-Connecting-Ip", "217.216.97.90")
	req.Header.Add("Cf-Ew-Via", "15")
	req.Header.Add("Cf-Ipcountry", "DE")
	req.Header.Add("Cf-Ray", "9e3709414e65a043-FRA")
	req.Header.Add("Cf-Visitor", "{\"scheme\":\"https\"}")
	req.Header.Add("Cf-Worker", "html-load.cc")
	req.Header.Add("Forwarded", "for=217.216.97.90")
	req.Header.Add("Host", "prebid.wurflcloud.com")
	req.Header.Add("Referer", "https://vsco.co/")
	req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"146\", \"Not-A.Brand\";v=\"24\", \"Android WebView\";v=\"146\"")
	req.Header.Add("Sec-Ch-Ua-Full-Version", "\"146.0.7680.157\"")
	req.Header.Add("Sec-Ch-Ua-Full-Version-List", "\"Chromium\";v=\"146.0.7680.157\", \"Not-A.Brand\";v=\"24.0.0.0\", \"Android WebView\";v=\"146.0.7680.157\"")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Add("Sec-Ch-Ua-Model", "\"XQ-CT54\"")
	req.Header.Add("Sec-Ch-Ua-Platform", "\"Android\"")
	req.Header.Add("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")
	req.Header.Add("Sec-Fetch-Dest", "script")
	req.Header.Add("Sec-Fetch-Mode", "no-cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-Storage-Access", "active")
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-Amzn-Trace-Id", "Root=1-69c7d9dc-618fdb3e4d9456fc38f66cab")
	req.Header.Add("X-As-Script-Type", "ESSENTIAL")
	req.Header.Add("X-Device-Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("X-Device-Ip", "217.216.97.90")
	req.Header.Add("X-Device-Referer", "https://vsco.co/")
	req.Header.Add("X-Device-User-Agent", UserAgent)
	req.Header.Add("X-Forwarded-For", "217.216.97.90, 172.71.144.58")
	req.Header.Add("X-Requested-With", "com.instagram.android")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupRequest(req)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupWithImportantHeaderMap_Long_Cache
func Benchmark_LookupWithImportantHeaderMap_Long_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 14; XQ-CT54 Build/64.2.A.2.258; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.157 Mobile Safari/537.36 Instagram 422.0.0.44.64 Android (34/14; 420dpi; 1096x2560; Sony; XQ-CT54; XQ-CT54; qcom; en_GB; 916494010; IABMV/1)"

	IHMap := make(map[string]string)
	IHMap["Accept"] = "*/*"
	IHMap["Accept-Encoding"] = "gzip, br"
	IHMap["Accept-Language"] = "en-GB,en-US;q=0.9,en;q=0.8"
	IHMap["Cdn-Loop"] = "cloudflare; loops=1"
	IHMap["Cf-Connecting-Ip"] = "217.216.97.90"
	IHMap["Cf-Ew-Via"] = "15"
	IHMap["Cf-Ipcountry"] = "DE"
	IHMap["Cf-Ray"] = "9e3709414e65a043-FRA"
	IHMap["Cf-Visitor"] = "{\"scheme\":\"https\"}"
	IHMap["Cf-Worker"] = "html-load.cc"
	IHMap["Forwarded"] = "for=217.216.97.90"
	IHMap["Host"] = "prebid.wurflcloud.com"
	IHMap["Referer"] = "https://vsco.co/"
	IHMap["Sec-Ch-Ua"] = "\"Chromium\";v=\"146\", \"Not-A.Brand\";v=\"24\", \"Android WebView\";v=\"146\""
	IHMap["Sec-Ch-Ua-Full-Version"] = "\"146.0.7680.157\""
	IHMap["Sec-Ch-Ua-Full-Version-List"] = "\"Chromium\";v=\"146.0.7680.157\", \"Not-A.Brand\";v=\"24.0.0.0\", \"Android WebView\";v=\"146.0.7680.157\""
	IHMap["Sec-Ch-Ua-Mobile"] = "?1"
	IHMap["Sec-Ch-Ua-Model"] = "\"XQ-CT54\""
	IHMap["Sec-Ch-Ua-Platform"] = "\"Android\""
	IHMap["Sec-Ch-Ua-Platform-Version"] = "\"14.0.0\""
	IHMap["Sec-Fetch-Dest"] = "script"
	IHMap["Sec-Fetch-Mode"] = "no-cors"
	IHMap["Sec-Fetch-Site"] = "cross-site"
	IHMap["Sec-Fetch-Storage-Access"] = "active"
	IHMap["User-Agent"] = UserAgent
	IHMap["X-Amzn-Trace-Id"] = "Root=1-69c7d9dc-618fdb3e4d9456fc38f66cab"
	IHMap["X-As-Script-Type"] = "ESSENTIAL"
	IHMap["X-Device-Accept-Language"] = "en-GB,en-US;q=0.9,en;q=0.8"
	IHMap["X-Device-Ip"] = "217.216.97.90"
	IHMap["X-Device-Referer"] = "https://vsco.co/"
	IHMap["X-Device-User-Agent"] = UserAgent
	IHMap["X-Forwarded-For"] = "217.216.97.90, 172.71.144.58"
	IHMap["X-Requested-With"] = "com.instagram.android"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupWithImportantHeaderMap(IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupWithImportantHeaderMap_Long_NoCache
func Benchmark_LookupWithImportantHeaderMap_Long_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 14; XQ-CT54 Build/64.2.A.2.258; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.157 Mobile Safari/537.36 Instagram 422.0.0.44.64 Android (34/14; 420dpi; 1096x2560; Sony; XQ-CT54; XQ-CT54; qcom; en_GB; 916494010; IABMV/1)"

	IHMap := make(map[string]string)
	IHMap["Accept"] = "*/*"
	IHMap["Accept-Encoding"] = "gzip, br"
	IHMap["Accept-Language"] = "en-GB,en-US;q=0.9,en;q=0.8"
	IHMap["Cdn-Loop"] = "cloudflare; loops=1"
	IHMap["Cf-Connecting-Ip"] = "217.216.97.90"
	IHMap["Cf-Ew-Via"] = "15"
	IHMap["Cf-Ipcountry"] = "DE"
	IHMap["Cf-Ray"] = "9e3709414e65a043-FRA"
	IHMap["Cf-Visitor"] = "{\"scheme\":\"https\"}"
	IHMap["Cf-Worker"] = "html-load.cc"
	IHMap["Forwarded"] = "for=217.216.97.90"
	IHMap["Host"] = "prebid.wurflcloud.com"
	IHMap["Referer"] = "https://vsco.co/"
	IHMap["Sec-Ch-Ua"] = "\"Chromium\";v=\"146\", \"Not-A.Brand\";v=\"24\", \"Android WebView\";v=\"146\""
	IHMap["Sec-Ch-Ua-Full-Version"] = "\"146.0.7680.157\""
	IHMap["Sec-Ch-Ua-Full-Version-List"] = "\"Chromium\";v=\"146.0.7680.157\", \"Not-A.Brand\";v=\"24.0.0.0\", \"Android WebView\";v=\"146.0.7680.157\""
	IHMap["Sec-Ch-Ua-Mobile"] = "?1"
	IHMap["Sec-Ch-Ua-Model"] = "\"XQ-CT54\""
	IHMap["Sec-Ch-Ua-Platform"] = "\"Android\""
	IHMap["Sec-Ch-Ua-Platform-Version"] = "\"14.0.0\""
	IHMap["Sec-Fetch-Dest"] = "script"
	IHMap["Sec-Fetch-Mode"] = "no-cors"
	IHMap["Sec-Fetch-Site"] = "cross-site"
	IHMap["Sec-Fetch-Storage-Access"] = "active"
	IHMap["User-Agent"] = UserAgent
	IHMap["X-Amzn-Trace-Id"] = "Root=1-69c7d9dc-618fdb3e4d9456fc38f66cab"
	IHMap["X-As-Script-Type"] = "ESSENTIAL"
	IHMap["X-Device-Accept-Language"] = "en-GB,en-US;q=0.9,en;q=0.8"
	IHMap["X-Device-Ip"] = "217.216.97.90"
	IHMap["X-Device-Referer"] = "https://vsco.co/"
	IHMap["X-Device-User-Agent"] = UserAgent
	IHMap["X-Forwarded-For"] = "217.216.97.90, 172.71.144.58"
	IHMap["X-Requested-With"] = "com.instagram.android"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupWithImportantHeaderMap(IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupDeviceIDWithRequest_Cache
func Benchmark_LookupDeviceIDWithRequest_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	DeviceID := "samsung_sm_m315f_ver1_suban110"
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

		device, _ := wengine.LookupDeviceIDWithRequest(DeviceID, req)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupDeviceIDWithRequest_NoCache
func Benchmark_LookupDeviceIDWithRequest_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	DeviceID := "samsung_sm_m315f_ver1_suban110"
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

		device, _ := wengine.LookupDeviceIDWithRequest(DeviceID, req)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupDeviceIDWithImportantHeaderMap_Cache
func Benchmark_LookupDeviceIDWithImportantHeaderMap_Cache(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	DeviceID := "samsung_sm_m315f_ver1_suban110"

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["Sec-CH-UA"] = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	IHMap["Sec-CH-UA-Full-Version"] = "90.0.4430.91"
	IHMap["Sec-CH-UA-Platform"] = "Android"
	IHMap["Sec-CH-UA-Platform-Version"] = "11"
	IHMap["Sec-CH-UA-Model"] = "SM-M315F"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupDeviceIDWithImportantHeaderMap(DeviceID, IHMap)
		device.Destroy()
	}

	b.StopTimer()
}

// Benchmark_LookupDeviceIDWithImportantHeaderMap_NoCache
func Benchmark_LookupDeviceIDWithImportantHeaderMap_NoCache(b *testing.B) {
	wengine := fixtureCreateEngineCachesize(nil, "")
	defer wengine.Destroy()
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	DeviceID := "samsung_sm_m315f_ver1_suban110"

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["Sec-CH-UA"] = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\""
	IHMap["Sec-CH-UA-Full-Version"] = "90.0.4430.91"
	IHMap["Sec-CH-UA-Platform"] = "Android"
	IHMap["Sec-CH-UA-Platform-Version"] = "11"
	IHMap["Sec-CH-UA-Model"] = "SM-M315F"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		device, _ := wengine.LookupDeviceIDWithImportantHeaderMap(DeviceID, IHMap)
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

// sink prevents the compiler from optimizing away benchmark results
var benchSink unsafe.Pointer

func Benchmark_TrieGet(b *testing.B) {
	headers := []string{
		"Accept-Encoding", "CAST-DEVICE-CAPABILITIES", "Device-Stock-UA",
		"Sec-CH-UA", "Sec-CH-UA-Arch", "Sec-CH-UA-Full-Version",
		"Sec-CH-UA-Full-Version-List", "Sec-CH-UA-Mobile", "Sec-CH-UA-Model",
		"Sec-CH-UA-Platform", "Sec-CH-UA-Platform-Version", "User-Agent",
		"X-OperaMini-Phone-UA", "X-Requested-With", "X-UCBrowser-Device-UA",
	}
	trieGet := wurfl.BenchmarkableTrieGet(headers)

	b.Run("ExactCase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchSink = trieGet("Sec-CH-UA-Full-Version-List")
		}
	})
	b.Run("MixedCase", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchSink = trieGet("sEc-cH-uA-fUlL-vErSiOn-LiSt")
		}
	})
	b.Run("NotFound", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchSink = trieGet("X-Not-A-Real-Header")
		}
	})
	b.Run("ShortKey", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchSink = trieGet("User-Agent")
		}
	})
}

// Benchmark_HeaderLookupStrategies compares four strategies for case-insensitive
// header name lookup using the same set of 15 WURFL important header names.
//
// Strategies:
//   - Trie: | 0x20 byte folding per byte, O(len(key)), 0 allocs
//   - Map: strings.ToLower + map access, 1 alloc (unavoidable for hashing)
//   - SequentialEqualFold: linear scan with strings.EqualFold, 0 allocs
//   - BinarySearchFold: sort.Search with | 0x20 comparator + EqualFold confirm, 0 allocs
func Benchmark_HeaderLookupStrategies(b *testing.B) {
	headers := []string{
		"Accept-Encoding", "CAST-DEVICE-CAPABILITIES", "Device-Stock-UA",
		"Sec-CH-UA", "Sec-CH-UA-Arch", "Sec-CH-UA-Full-Version",
		"Sec-CH-UA-Full-Version-List", "Sec-CH-UA-Mobile", "Sec-CH-UA-Model",
		"Sec-CH-UA-Platform", "Sec-CH-UA-Platform-Version", "User-Agent",
		"X-OperaMini-Phone-UA", "X-Requested-With", "X-UCBrowser-Device-UA",
	}

	trieGet := wurfl.BenchmarkableTrieGet(headers)
	mapGet := wurfl.BenchmarkableMapGet(headers)
	seqGet := wurfl.BenchmarkableSequentialEqualFoldGet(headers)
	bsfGet := wurfl.BenchmarkableBinarySearchFoldGet(headers)

	keys := []struct {
		name string
		key  string
	}{
		{"LongKey_MixedCase", "sEc-cH-uA-fUlL-vErSiOn-LiSt"},
		{"ShortKey_MixedCase", "uSeR-aGeNt"},
		{"NotFound", "X-Not-A-Real-Header"},
	}

	for _, k := range keys {
		b.Run("Trie/"+k.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchSink = trieGet(k.key)
			}
		})
		b.Run("Map/"+k.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchSink = mapGet(k.key)
			}
		})
		b.Run("SequentialEqualFold/"+k.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchSink = seqGet(k.key)
			}
		})
		b.Run("BinarySearchFold/"+k.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchSink = bsfGet(k.key)
			}
		})
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

func TestP99_LookupRequest_NoCache(t *testing.T) {
	wengine := fixtureCreateEngineCachesize(t, "")
	defer wengine.Destroy()
	ua := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", ua)
	req.Header.Add("Sec-CH-UA", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Add("Sec-CH-UA-Full-Version", "90.0.4430.91")
	req.Header.Add("Sec-CH-UA-Platform", "Android")
	req.Header.Add("Sec-CH-UA-Platform-Version", "11")
	req.Header.Add("Sec-CH-UA-Model", "SM-M315F")

	const iterations = 1000000
	durations := make([]time.Duration, iterations)

	// warmup
	for i := 0; i < 100; i++ {
		device, _ := wengine.LookupRequest(req)
		device.Destroy()
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		device, _ := wengine.LookupRequest(req)
		device.Destroy()
		durations[i] = time.Since(start)
	}

	slices.Sort(durations)
	p50 := durations[iterations*50/100]
	p95 := durations[iterations*95/100]
	p99 := durations[iterations*99/100]
	p100 := durations[iterations-1]
	t.Logf("Iterations: %d", iterations)
	t.Logf("P50:  %v", p50)
	t.Logf("P95:  %v", p95)
	t.Logf("P99:  %v", p99)
	t.Logf("P100: %v", p100)
}

func TestP99_LookupUserAgent_NoCache(t *testing.T) {
	wengine := fixtureCreateEngineCachesize(t, "")
	defer wengine.Destroy()
	ua := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	const iterations = 1000000
	durations := make([]time.Duration, iterations)

	// warmup
	for i := 0; i < 1000; i++ {
		device, _ := wengine.LookupUserAgent(ua)
		device.Destroy()
	}

	for i := 0; i < iterations; i++ {
		start := time.Now()
		device, _ := wengine.LookupUserAgent(ua)
		device.Destroy()
		durations[i] = time.Since(start)
	}

	slices.Sort(durations)
	p50 := durations[iterations*50/100]
	p95 := durations[iterations*95/100]
	p99 := durations[iterations*99/100]
	p100 := durations[iterations-1]
	t.Logf("Iterations: %d", iterations)
	t.Logf("P50:  %v", p50)
	t.Logf("P95:  %v", p95)
	t.Logf("P99:  %v", p99)
	t.Logf("P100: %v", p100)
}

// used to prevent compiler optimizations in the BenchmarkStringToLower
var sinkString string

func Benchmark_StringToLower(b *testing.B) {
	headerName := "Sec-CH-UA-Platform-Version"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkString = strings.ToLower(headerName)
	}
	b.StopTimer()
}

var sinkBool bool

func Benchmark_StringEqualFold(b *testing.B) {
	headerName := "Sec-CH-UA-Platform-Version"
	headerNameLower := "sec-ch-ua-platform-version"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkBool = strings.EqualFold(headerName, headerNameLower)
	}
	b.StopTimer()
}
