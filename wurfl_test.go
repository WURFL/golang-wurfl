// WURFL InFuze Golang module
//

package wurfl_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	wurfl "github.com/WURFL/golang-wurfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func fixtureCreateEngine(t *testing.T) *wurfl.Wurfl {
	var wengine *wurfl.Wurfl
	var err error

	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		wengine, err = wurfl.Create("/usr/local/share/wurfl/wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
	} else {
		// all other systems (TODO windows)
		wengine, err = wurfl.Create("/usr/share/wurfl/wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
	}
	// cant use assert 'cause this method is also used in benchmarks (no *testing.T)
	// require.NoErrorf(t, err, "Create returned an error: %s", err)
	if err != nil {
		t.Errorf("Create returned an error: %s", err)
	}

	return wengine
}

func fixtureCreateEngineCachesize(t *testing.T, cachesize string) *wurfl.Wurfl {
	var wengine *wurfl.Wurfl
	var err error
	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		if cachesize == "" {
			wengine, err = wurfl.Create("/usr/local/share/wurfl/wurfl.zip",
				nil, nil, -1, wurfl.WurflCacheProviderNone, cachesize)
		} else {
			wengine, err = wurfl.Create("/usr/local/share/wurfl/wurfl.zip",
				nil, nil, -1, wurfl.WurflCacheProviderLru, cachesize)
		}
	} else {
		// all other systems (TODO windows)
		if cachesize == "" {
			wengine, err = wurfl.Create("/usr/share/wurfl/wurfl.zip",
				nil, nil, -1, wurfl.WurflCacheProviderNone, cachesize)
		} else {
			wengine, err = wurfl.Create("/usr/share/wurfl/wurfl.zip",
				nil, nil, -1, wurfl.WurflCacheProviderLru, cachesize)
		}
	}
	// cant use assert 'cause this method is also used in benchmarks (no *testing.T)
	// require.NoErrorf(t, err, "Create returned an error: %s", err)
	if err != nil {
		t.Errorf("Create returned an error: %s", err)
	}

	return wengine
}

func TestWurfl_WurflGetAPIVersion(t *testing.T) {
	ver := wurfl.APIVersion()
	assert.NotEmpty(t, ver)
	fmt.Printf("WURFL API Version: %s\n", ver)
}

func TestWurfl_Create(t *testing.T) {
	ua := "ArtDeviant/3.0.2 CFNetwork/711.3.18 Darwin/14.0.0"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	fmt.Printf("WURFL API Version: %s\n", wengine.GetAPIVersion())

	device, err := wengine.LookupUserAgent(ua)
	assert.Equal(t, nil, err)

	deviceid, err := device.GetDeviceID()
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, deviceid)
	assert.Equal(t, "apple_iphone_ver8_3_subuacfnetwork", deviceid)
}

func TestWurfl_Lookup(t *testing.T) {
	var wengine *wurfl.Wurfl
	var err error

	// Capability filtering is discouraged and will be deprecated. Here only for test coverage purposes
	capfilter := []string{
		"mobile_browser_version",
		"pointing_method",
		"is_tablet",
	}

	caps := capfilter[0:3]

	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		wengine, err = wurfl.Create("/usr/local/share/wurfl/wurfl.zip", nil, caps, -1, wurfl.WurflCacheProviderLru, "100000")
	} else {
		// all other systems (TODO windows)
		wengine, err = wurfl.Create("/usr/share/wurfl/wurfl.zip", nil, caps, -1, wurfl.WurflCacheProviderLru, "100000")
	}

	require.NoErrorf(t, err, "Create returned an error: %s", err)
	defer wengine.Destroy()

	wengine.SetLogPath("api.log")

	ua := "ArtDeviant/3.0.2 CFNetwork/711.3.18 Darwin/14.0.0"

	device, err := wengine.LookupUserAgent(ua)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)

	deviceid, err := device.GetDeviceID()
	assert.NoError(t, err)

	newdevice, err := wengine.LookupDeviceID(deviceid)
	assert.NoError(t, err)

	deviceid2, err2 := newdevice.GetDeviceID()
	assert.NoError(t, err2)

	if deviceid != deviceid2 {
		t.Errorf("Error, devices do not match %s, %s", deviceid, deviceid2)
	}

	_, uaerr := device.GetUserAgent()
	assert.NoError(t, uaerr)

	oua, uaerr := device.GetOriginalUserAgent()
	assert.NoError(t, uaerr)

	if oua != ua {
		t.Errorf("Error, ua matched >%s< and device original ua >%s< do not match", ua, oua)
	}

	if device.IsRoot() {
		fmt.Printf("Device is root\n")
	}

	if device.GetCapability("mobile_browser_version") != "8.0" {
		t.Errorf("device.GetCapability(\"mobile_browser_version\") does not return 8.0 : %s\n", device.GetCapability("mobile_browser_version"))
	}

	vcaps := device.GetVirtualCapabilities(wengine.GetAllVCaps())

	if vcaps["advertised_device_os"] != "iOS" {
		t.Errorf("device.GetVirtualCapabilities() \"advertised_device_os\" != \"iOS\" : %s\n", vcaps["advertised_device_os"])
	}

	allcaps := device.GetCapabilities(wengine.GetAllCaps())

	if allcaps["device_os"] != "iOS" {
		t.Errorf("device.GetCapabilities() \"device_os\" != \"iOS\" : %s\n", allcaps["device_os"])
	}

	device.GetMatchType()
	device.GetRootID()

	device.Destroy()
	newdevice.Destroy()
	fi, apilogoserr := os.Stat("api.log")
	assert.Nil(t, apilogoserr)
	assert.NotNil(t, fi)
	os.Remove("api.log")
}

func TestWurfl_GetAllVCaps(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	s := wengine.GetAllVCaps()
	if len(s) < 28 {
		t.Errorf("Vcaps should be 28 or more, they are %d", len(s))
	}

}

// Test_GetCapability : various cases on GetCapability / GetVirtualCapability
func Test_GetCapability(t *testing.T) {
	assert := assert.New(t)
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	assert.Nil(err)
	if assert.NotNil(device) {
		assert.Equal("Android", device.GetCapability("device_os"))
		// non existent
		assert.Equal("", device.GetCapability("duvice_as"))
		// non existent, new method with error WURFL_ERROR_CAPABILITY_NOT_FOUND
		val, err := device.GetStaticCap("duvice_as")
		assert.Equal("", val)
		assert.NotEqual(nil, err)
		// non existing vcap
		assert.Equal("", device.GetVirtualCapability("notexist"))
		// non existing vcap, with error WURFL_ERROR_VIRTUAL_CAPABILITY_NOT_FOUND
		val, err = device.GetVirtualCap("notexist")
		assert.Equal("", val)
		assert.NotEqual(nil, err)
		// existing vcap
		assert.Equal("428", device.GetVirtualCapability("pixel_density"))
		// falback on vcap from cap
		assert.Equal("Smartphone", device.GetCapability("form_factor"))
	}
	device.Destroy()
}

func Test_LookupRequest(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// fmt.Println(wengine.ImportantHeaderNames)

	// lookup both UAs. When using both headers in LookupRequest, the XUCBrowserDeviceUA has precedence.
	UserAgentDevice, err := wengine.LookupUserAgent(UserAgent)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)

	UserAgentDeviceID, _ := UserAgentDevice.GetDeviceID()

	XUCBrowserDeviceDevice, err := wengine.LookupUserAgent(XUCBrowserDeviceUA)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)

	XUCBrowserDeviceDeviceID, _ := XUCBrowserDeviceDevice.GetDeviceID()

	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-UCBrowser-Device-UA", XUCBrowserDeviceUA)

	reqDevice, err := wengine.LookupRequest(req)
	assert.NoErrorf(t, err, "LookupRequest returned an error: %s", err)

	reqDeviceID, _ := reqDevice.GetDeviceID()

	// now verify that device retrieved with headers is the the same as XUCBrowserDeviceDevice
	if UserAgentDeviceID == reqDeviceID {
		t.Errorf("Devices are the same, should be different : %s, %s\n", UserAgentDeviceID, reqDeviceID)
	}

	if XUCBrowserDeviceDeviceID != reqDeviceID {
		t.Errorf("Devices are different, should be the same: %s, %s\n", XUCBrowserDeviceDeviceID, reqDeviceID)
	}

	UserAgentDevice.Destroy()
	XUCBrowserDeviceDevice.Destroy()
	reqDevice.Destroy()
}

// Test_LookupRequestExperimental : test new sec-ch headers
func Test_LookupRequestExperimental(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	// set experimental headers
	wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, 1)

	// User-Agent
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	UserAgentDevice, err := wengine.LookupUserAgent(UserAgent)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)

	UserAgentDeviceID, _ := UserAgentDevice.GetDeviceID()

	fmt.Println(wengine.ImportantHeaderNames)

	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Sec-CH-UA", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
	req.Header.Add("Sec-CH-UA-Full-Version", "90.0.4430.91")
	req.Header.Add("Sec-CH-UA-Platform", "Android")
	req.Header.Add("Sec-CH-UA-Platform-Version", "11")
	req.Header.Add("Sec-CH-UA-Model", "SM-M315F")

	reqDevice, err := wengine.LookupRequest(req)
	assert.NoErrorf(t, err, "LookupRequest returned an error: %s", err)

	reqDeviceID, _ := reqDevice.GetDeviceID()

	// now verify that device retrieved with headers is the the same as XUCBrowserDeviceDevice
	if UserAgentDeviceID != reqDeviceID {
		t.Errorf("Devices should be the same : %s, %s\n", UserAgentDeviceID, reqDeviceID)
	}

	UserAgentDevice.Destroy()
	reqDevice.Destroy()
}

func Test_LookupWithImportantHeaderMap(t *testing.T) {
	var err error
	var deviceLA *wurfl.Device
	var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// do a LookupUserAgent() with UA

	// fmt.Println(wengine.ImportantHeaderNames)

	deviceLA, err = wengine.LookupUserAgent(UserAgent)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)
	LaDeviceID, _ := deviceLA.GetDeviceID()

	// create IHMap and lookup using headers

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["X-UCBrowser-Device-UA"] = XUCBrowserDeviceUA

	deviceIHM, err = wengine.LookupWithImportantHeaderMap(IHMap)
	assert.NoErrorf(t, err, "LookupWithImportantHeaderMap returned an error: %s", err)

	IHMDeviceID, _ := deviceIHM.GetDeviceID()

	if LaDeviceID == IHMDeviceID {
		t.Errorf("Devices are the same, should be different : %s, %s\n", LaDeviceID, IHMDeviceID)
	}

	deviceLA.Destroy()
	deviceIHM.Destroy()
}

func lookUpUserAgent(wengine *wurfl.Wurfl, ua string, capabilities []string) map[string]string {
	device, _ := wengine.LookupUserAgent(ua)
	defer device.Destroy()
	return device.GetCapabilities(capabilities)
}

func TestJira_INFUZE1053(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()
	capabilitiesVector := []string{"max_image_width", "dual_orientation", "density_class", "is_ios"}
	capabilitiesMap := lookUpUserAgent(wengine, "Callpod Keeper for Android 1.0 (10.0.0/234) Dalvik/2.1.0 (Linux; U; Android 5.0.1; SAMSUNG-SGH-I337 Build/LRX22C)", capabilitiesVector)
	fmt.Println(capabilitiesMap)
}

func Test_LookupDeviceIDWithImportantHeaderMap(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// create IHMap and lookup using headers

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["X-UCBrowser-Device-UA"] = XUCBrowserDeviceUA

	deviceIHM, err := wengine.LookupWithImportantHeaderMap(IHMap)
	assert.NoErrorf(t, err, "LookupWithImportantHeaderMap returned an error: %s", err)

	IHMDeviceID, _ := deviceIHM.GetDeviceID()
	AdvBrow1, _ := deviceIHM.GetVirtualCap("advertised_browser")

	// now lookup by deviceID and no header and check that an advertised vcap behaves correctly
	deviceIDIHM, err := wengine.LookupDeviceID(IHMDeviceID)
	assert.NoErrorf(t, err, "LookupDeviceID returned an error: %s", err)

	AdvBrow2, _ := deviceIDIHM.GetVirtualCap("advertised_browser")
	if AdvBrow1 == AdvBrow2 {
		t.Errorf("advertised_browser are the same, should be different : %s, %s\n", AdvBrow1, AdvBrow2)
	}

	// now lookup by deviceID and header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceIDWithImportantHeaderMap(IHMDeviceID, IHMap)
	assert.NoErrorf(t, err, "LookupDeviceIDWithImportantHeaderMap returned an error: %s", err)

	AdvBrow3, _ := deviceIDIHM.GetVirtualCap("advertised_browser")

	if AdvBrow1 != AdvBrow3 {
		t.Errorf("advertised_browser are different, should be the same: %s, %s\n", AdvBrow1, AdvBrow3)
	}

	deviceIHM.Destroy()
}

func Test_LookupWithImportantHeaderMapCaseInsensitive(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	UserAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

	deviceLA, err := wengine.LookupUserAgent(UserAgent)
	assert.NoErrorf(t, err, "LookupUserAgent returned an error: %s", err)

	LaDeviceID, _ := deviceLA.GetDeviceID()

	// create IHMap and lookup using headers. Headers names will have some case differences
	// to check that LookuupWithImportantHeaderMap() is case insensitive

	// Where strings could contain double quotes (") delimit values using backticks (`) to avoid escaping
	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["accept-encoding"] = "gzip, deflate, br, zstd"
	IHMap["Sec-ch-UA-Platform"] = "Android"
	IHMap["Sec-CH-ua"] = `"Chromium";v="122", "Not(A:Brand";v="24", "Veera";v="122"`
	IHMap["SEC-CH-UA-MOBILE"] = "?1"
	IHMap["seC-cH-uA-fulL-versioN-lisT"] = `"Chromium";v="122.0.0.0", "Not(A:Brand";v="24.0.0.0", "Veera";v="122.0.0.0"`

	deviceIHM, err := wengine.LookupWithImportantHeaderMap(IHMap)
	assert.NoErrorf(t, err, "LookupWithImportantHeaderMap returned an error: %s", err)

	IHMDeviceID, _ := deviceIHM.GetDeviceID()

	if LaDeviceID == IHMDeviceID {
		t.Errorf("Devices are the same, should be different : %s, %s\n", LaDeviceID, IHMDeviceID)
	}

	deviceLA.Destroy()
	deviceIHM.Destroy()
}

func Test_LookupDeviceIDWithRequest(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-UCBrowser-Device-UA", XUCBrowserDeviceUA)

	deviceIHM, err := wengine.LookupRequest(req)
	assert.NoErrorf(t, err, "LookupRequest returned an error: %s", err)

	IHMDeviceID, _ := deviceIHM.GetDeviceID()
	AdvBrow1, _ := deviceIHM.GetVirtualCap("advertised_browser")

	// now lookup by deviceID and no header and check that an advertised vcap behaves correctly
	deviceIDIHM, err := wengine.LookupDeviceID(IHMDeviceID)
	assert.NoErrorf(t, err, "LookupDeviceID returned an error: %s", err)

	AdvBrow2, _ := deviceIDIHM.GetVirtualCap("advertised_browser")

	if AdvBrow1 == AdvBrow2 {
		t.Errorf("advertised_browser are the same, should be different : %s, %s\n", AdvBrow1, AdvBrow2)
	}

	// now lookup by deviceID and header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceIDWithRequest(IHMDeviceID, req)
	assert.NoErrorf(t, err, "LookupDeviceIDWithRequest returned an error: %s", err)

	AdvBrow3, _ := deviceIDIHM.GetVirtualCap("advertised_browser")

	if AdvBrow1 != AdvBrow3 {
		t.Errorf("advertised_browser are different, should be the same: %s, %s\n", AdvBrow1, AdvBrow3)
	}

	deviceIHM.Destroy()
}

func TestWurfl_UpdaterRunonce(t *testing.T) {
	var wurflZip string

	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		wurflZip = "/usr/local/share/wurfl/wurfl.zip"
	} else {
		// all other systems (TODO windows)
		wurflZip = "/usr/share/wurfl/wurfl.zip"
	}

	// copy wurfl.zip to /tmp
	cpCmd := exec.Command("cp", "-rf", wurflZip, "/tmp/wurfl.zip")
	_ = cpCmd.Run()

	info, _ := os.Stat("/tmp/wurfl.zip")
	mdt1 := info.ModTime()

	wengine, err := wurfl.Create("/tmp/wurfl.zip", nil, nil, wurfl.WurflEngineTargetHighAccuray, wurfl.WurflCacheProviderDoubleLru, "100000")
	require.NoErrorf(t, err, "Create returned an error: %s", err)
	defer wengine.Destroy()

	// set env var SM_UPDATER_DATA_URL to your updater URL (from scientiamobile Vault)

	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// set logger file
	_ = wengine.SetUpdaterLogPath("/tmp/wurfl-updater-log.txt")

	// set updater path
	uerr := wengine.SetUpdaterDataURL(URL)
	assert.NoErrorf(t, uerr, "SetUpdaterDataURL returned an error: %s", err)

	// set updater user-agent
	expUa := fmt.Sprintf("golang_wurfl_test/%s", wurfl.Version)

	t.Logf("Specific golang binding User-Agent string: %s", expUa)

	uerr = wengine.SetUpdaterUserAgent(expUa)
	assert.NoErrorf(t, uerr, "SetUpdaterUserAgent returned an error: %s", err)

	ua := wengine.GetUpdaterUserAgent()
	assert.NotEmptyf(t, ua, "SetUpdaterUserAgent returned an error: %s", err)

	// set timeout to defaults
	uerr = wengine.SetUpdaterDataURLTimeout(-1, -1)

	uerr = wengine.UpdaterRunonce()
	assert.NoErrorf(t, err, "UpdaterRunonce returned an error: %s", err)

	// check if the modification time of wurfl.zip has changed
	info, _ = os.Stat("/tmp/wurfl.zip")
	mdt2 := info.ModTime()

	if mdt1.Equal(mdt2) {
		t.Errorf("/tmp/wurfl.zip not downloaded\n")
	}

}
func TestWurfl_UpdaterThread(t *testing.T) {
	var wengine *wurfl.Wurfl
	var err error
	var wurflZip string

	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		wurflZip = "/usr/local/share/wurfl/wurfl.zip"
	} else {
		// all other systems (TODO windows)
		wurflZip = "/usr/share/wurfl/wurfl.zip"
	}

	// copy wurfl.zip to /tmp
	cpCmd := exec.Command("cp", "-rf", wurflZip, "/tmp/wurfl.zip")
	_ = cpCmd.Run()

	info, _ := os.Stat("/tmp/wurfl.zip")
	mdt1 := info.ModTime()

	wengine, err = wurfl.Create("/tmp/wurfl.zip", nil, nil, wurfl.WurflEngineTargetHighAccuray, wurfl.WurflCacheProviderDoubleLru, "100000")
	require.NoErrorf(t, err, "Create returned an error: %s", err)
	defer wengine.Destroy()

	_ = wengine.SetUpdaterLogPath("/tmp/wurfl-updater-log.txt")

	lastLoadTime := wengine.GetLastLoadTime()

	// set env var SM_UPDATER_DATA_URL to your updater URL (from scientiamobile Vault)

	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// set updater path
	uerr := wengine.SetUpdaterDataURL(URL)
	assert.NoErrorf(t, uerr, "SetUpdaterDataURL returned an error: %s", err)

	uerr = wengine.UpdaterStart()
	assert.NoErrorf(t, uerr, "UpdaterStart returned an error: %s", err)

	// wait for updater thread to finish first update
	time.Sleep(30 * time.Second)

	// check if the modification time of wurfl.zip has changed
	info, _ = os.Stat("/tmp/wurfl.zip")
	mdt2 := info.ModTime()

	if mdt1.Equal(mdt2) {
		t.Errorf("/tmp/wurfl.zip not downloaded\n")
	}

	newLastLoadTime := wengine.GetLastLoadTime()

	if newLastLoadTime == lastLoadTime {
		t.Errorf("Engine not reloaded\n")
	}

}

func TestWurfl_Getters(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	assert.NotEmpty(t, wengine.GetLastLoadTime(), "GetLastLoadTime should return non-empty wurfl info")
	assert.NotEmpty(t, wengine.GetInfo(), "GetInfo should return non-empty wurfl info")
	assert.False(t, wengine.HasVirtualCapability("pippo"), "HasVirtualCapability should return false for non-existent capability")
	assert.True(t, wengine.HasVirtualCapability("is_ios"), "HasVirtualCapability should return true for existent capability")
	assert.True(t, wengine.HasCapability("device_os"), "HasCapability should return true for existent capability")
	assert.False(t, wengine.HasCapability("pippo"), "HasCapability should return false for non-existent capability")

}

func TestWurfl_GetAllDeviceIds(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	ids := wengine.GetAllDeviceIds()
	assert.NotEqual(t, 0, len(ids))
}

// Test_SetAttr :
func Test_SetAttr(t *testing.T) {

	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	t.Run("setattr", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, 10)
		assert.NoErrorf(t, err, "SetAttr returned an error: %s", err)
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)
		assert.Equal(t, 10, attrValue)
	})

	t.Run("setattr negative value", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, -10)
		assert.NoErrorf(t, err, "SetAttr returned an error: %s", err)
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)
		assert.Equal(t, -10, attrValue)
	})

	t.Run("setattr invalid attr", func(t *testing.T) {
		err = wengine.SetAttr(44, 10)
		assert.NotNil(t, err)
	})
}

func Test_SetAttr_FallbackCache(t *testing.T) {
	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	t.Run("setAttr fallback cache - default", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		assert.NoErrorf(t, err, "SetAttr returned an error: %s", err)

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDefault {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		}
	})

	t.Run("setAttr fallback cache - disabled", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
		assert.NoErrorf(t, err, "SetAttr returned an error: %s", err)

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDisabled {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
		}
	})

	t.Run("setAttr fallback cache - limited", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		assert.NoErrorf(t, err, "SetAttr returned an error: %s", err)

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheLimited {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		}

		// check that a new set overwrites the old one
		wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDefault {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		}

	})
}

// Test_GetAttr :
func Test_GetAttr(t *testing.T) {

	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	t.Run("getattr", func(t *testing.T) {
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)
		assert.NoErrorf(t, err, "GetAttr returned an error: %s", err)
		assert.Equal(t, 1, attrValue)
	})

	t.Run("getattr invalid attr", func(t *testing.T) {
		attrValue, err = wengine.GetAttr(44)
		assert.NotNil(t, err)
	})
}

func TestWurfl_IsUserAgentFrozen(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want bool
	}{
		{
			name: "not frozen UA",
			ua:   "ArtDeviant/3.0.2 CFNetwork/711.3.18 Darwin/14.0.0",
			want: false,
		},
		{
			name: "frozen UA",
			ua:   "Mozilla/5.0 (Linux; Android 6.0; AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.0.0 Mobile Safari/537.36",
			want: false,
		},
	}

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wengine.IsUserAgentFrozen(tt.ua)
			if got != tt.want {
				t.Errorf("Wurfl.IsUserAgentFrozen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWurfl_GetHeaderQuality(t *testing.T) {
	tests := []struct {
		name    string
		h       http.Header
		want    wurfl.HeaderQuality
		wantErr bool
	}{
		{
			name: "full-1",
			h: http.Header{
				"User-Agent":                 {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"},
				"Sec-Ch-Ua":                  {`" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`},
				"Sec-Ch-Ua-Full-Version":     {"96.0.4664.93"},
				"Sec-Ch-Ua-Platform":         {"Linux"},
				"Sec-Ch-Ua-Platform-Version": {"5.4.0"},
			},
			want: wurfl.HeaderQualityFull,
		},
		{
			name: "full-2",
			h: http.Header{
				"User-Agent":                 {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.57"},
				"Sec-Ch-Ua":                  {`" Not A;Brand";v="99", "Chromium";v="96", "Microsoft Edge";v="96"`},
				"Sec-Ch-Ua-Full-Version":     {"96.0.1054.57"},
				"Sec-Ch-Ua-Platform":         {"Windows"},
				"Sec-Ch-Ua-Platform-Version": {"10.0.0"},
			},
			want: wurfl.HeaderQualityFull,
		},
		{
			name: "full-3",
			h: http.Header{
				"User-Agent":                 {"Mozilla/5.0 (Linux; Android 12; Pixel 4 XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36"},
				"Sec-Ch-Ua":                  {`" Not A;Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
				"Sec-Ch-Ua-Full-Version":     {"97.0.4692.70"},
				"Sec-Ch-Ua-Platform":         {"Android"},
				"Sec-Ch-Ua-Platform-Version": {"12.0.0"},
				"Sec-Ch-Ua-Model":            {"Pixel 4 XL"},
			},
			want: wurfl.HeaderQualityFull,
		},
		{
			name: "full-4",
			h: http.Header{
				"User-Agent":                 {"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36"},
				"Sec-Ch-Ua":                  {`" Not A;Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
				"Sec-Ch-Ua-Full-Version":     {"97.0.4692.70"},
				"Sec-Ch-Ua-Platform":         {"Android"},
				"Sec-Ch-Ua-Platform-Version": {"12.0.0"},
				"Sec-Ch-Ua-Model":            {"Pixel 4 XL"},
			},
			want: wurfl.HeaderQualityFull,
		},
		{
			name: "basic-non-frozen-1",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`},
				"Sec-Ch-Ua-Platform": {"Linux"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "basic-non-frozen-2",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.57"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Chromium";v="96", "Microsoft Edge";v="96"`},
				"Sec-Ch-Ua-Platform": {"Windows"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "basic-non-frozen-3",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (Linux; Android 12; Pixel 4 XL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
				"Sec-Ch-Ua-Platform": {"Android"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "basic-frozen-1",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.0.0 Safari/537.36"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Chromium";v="96", "Google Chrome";v="96"`},
				"Sec-Ch-Ua-Platform": {"Linux"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "basic-frozen-2",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.0.0 Safari/537.36 Edg/96.0.0.0"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Chromium";v="96", "Microsoft Edge";v="96"`},
				"Sec-Ch-Ua-Platform": {"Windows"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "basic-frozen-3",
			h: http.Header{
				"User-Agent":         {"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36"},
				"Sec-Ch-Ua":          {`" Not A;Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`},
				"Sec-Ch-Ua-Platform": {"Android"},
			},
			want: wurfl.HeaderQualityBasic,
		},
		{
			name: "none-1",
			h: http.Header{
				"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.0.0 Safari/537.36"},
			},
			want: wurfl.HeaderQualityNone,
		},
		{
			name: "none-2",
			h: http.Header{
				"User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.0.0 Safari/537.36 Edg/96.0.0.0"},
				"Sec-Ch-Ua":  {`" Not A;Brand";v="99", "Chromium";v="96", "Microsoft Edge";v="96"`},
			},
			want: wurfl.HeaderQualityNone,
		},
		{
			name: "none-3",
			h: http.Header{
				"User-Agent": {"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36"},
			},
			want: wurfl.HeaderQualityNone,
		},
	}

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header = tt.h
		t.Run(tt.name, func(t *testing.T) {
			got, err := wengine.GetHeaderQuality(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wurfl.GetHeaderQualitys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Wurfl.GetHeaderQuality() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_GetCapabilityAsInt(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	defer device.Destroy()

	t.Run("Test GetCapabilityAsInt OK", func(t *testing.T) {
		cap, err := device.GetCapabilityAsInt("resolution_height")
		assert.NoErrorf(t, err, "GetCapabilityAsInt error %v", err)
		assert.GreaterOrEqual(t, cap, 0)
	})

	t.Run("Test GetCapabilityAsInt calling a non int capability, must return a not nil error", func(t *testing.T) {
		capname := "brand_name"
		// from 1.12.7.1 libwurfl returns error when asked for non numeric capabilities (ie: brand_name)
		_, err := device.GetCapabilityAsInt(capname)
		assert.Errorf(t, err, "GetCapabilityAsInt should return an error for non int capability %s", capname)
	})

	t.Run("Test GetCapabilityAsInt calling a capability using an empty string, must return a not nil error", func(t *testing.T) {
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetCapabilityAsInt("")
		assert.NotNil(t, err)
	})
}

func TestDevice_GetVirtualCapabilityAsInt(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)
	defer device.Destroy()

	t.Run("Test GetVirtualCapabilityAsInt OK", func(t *testing.T) {
		cap, err := device.GetVirtualCapabilityAsInt("pixel_density")
		assert.NoErrorf(t, err, "GetCapabilityAsInt error %v", err)
		assert.GreaterOrEqual(t, cap, 0)
	})

	t.Run("Test GetVirtualCapabilityAsInt calling a non int virtual capability, must return a not nil error", func(t *testing.T) {
		capname := "form_factor"
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetVirtualCapabilityAsInt(capname)
		assert.NotNil(t, err)
	})

	t.Run("Test GetVirtualCapabilityAsInt calling a virtual capability using an empty string, must return a not nil error", func(t *testing.T) {
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetVirtualCapabilityAsInt("")
		assert.NotNil(t, err)
	})
}

func TestDevice_GetRootID(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	/**
	* This device have no device roots in its fall back tree,
	* since no devices above it (itself included) are real devices (actual device roots),
	* in this case "" is expected.
	**/
	device, err := wengine.LookupDeviceID("generic_android_ver11_0_subff102_tablet")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	rootID := device.GetRootID()
	assert.Equal(t, "", rootID)

	// generic has empty root id
	device, err = wengine.LookupDeviceID("generic")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	rootID = device.GetRootID()
	assert.Equal(t, "", rootID)

	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1_subu3k10")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	rootID = device.GetRootID()
	assert.Equal(t, "natec_smart_tv_dongle_hd221_ver1", rootID)

	//is an actual device root , root is itself
	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	rootID = device.GetRootID()
	assert.Equal(t, "natec_smart_tv_dongle_hd221_ver1", rootID)
}

func TestDevice_GetParentID(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupDeviceID("generic_android_ver11_0_subff102_tablet")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	parentID := device.GetParentID()
	assert.Equal(t, "generic_android_ver11_0_subff101_tablet", parentID)

	// generic has empty parent id
	device, err = wengine.LookupDeviceID("generic")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	parentID = device.GetParentID()
	assert.Equal(t, "", parentID)

	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1_subu3k10")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	parentID = device.GetParentID()
	assert.Equal(t, "natec_smart_tv_dongle_hd221_ver1", parentID)

	//is an actual device root , but has a parent
	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1")
	require.NoErrorf(t, err, "LookupDeviceID error %v", err)

	parentID = device.GetParentID()
	assert.Equal(t, "generic_android_ver4_2", parentID)
}

// The Wurfl handler can be mocked implementing the WurflHandler interface
// Below is a simple mock that does not require external dependencies.
// If additional features are required third party mocking libraries can be used, like:
// - https://github.com/golang/mock
// - https://github.com/vektra/mockery
type MockWurfl struct {
}

func (m *MockWurfl) GetAPIVersion() string {
	return "version 1.10.0.0"
}

func (m *MockWurfl) Download(URL string, folder string) error {
	return nil
}

func (m *MockWurfl) SetLogPath(LogFile string) error {
	return nil
}

func (m *MockWurfl) LookupDeviceIDWithImportantHeaderMap(DeviceID string, IHMap map[string]string) (wurfl.DeviceHandler, error) {
	return nil, nil
}

func (m *MockWurfl) LookupWithImportantHeaderMap(IHMap map[string]string) (wurfl.DeviceHandler, error) {
	return nil, nil
}

func (m *MockWurfl) LookupDeviceIDWithRequest(DeviceID string, r *http.Request) (wurfl.DeviceHandler, error) {
	return nil, nil
}

func (m *MockWurfl) LookupRequest(r *http.Request) (wurfl.DeviceHandler, error) {
	return nil, nil
}

func (m *MockWurfl) LookupUserAgent(ua string) (wurfl.DeviceHandler, error) {
	return nil, nil
}
func (m *MockWurfl) GetAllDeviceIds() []string {
	return []string{"generic"}
}
func (m *MockWurfl) LookupDeviceID(DeviceID string) (wurfl.DeviceHandler, error) {
	return nil, nil
}
func (m *MockWurfl) IsUserAgentFrozen(ua string) bool {
	if ua == "frozen" {
		return true
	}
	return false
}
func (m *MockWurfl) GetHeaderQuality(r *http.Request) (wurfl.HeaderQuality, error) {
	switch r.Header.Get("x-wurfl-mock-getheaderquality") {
	case wurfl.HeaderQualityFull.String():
		return wurfl.HeaderQualityFull, nil
	case wurfl.HeaderQualityBasic.String():
		return wurfl.HeaderQualityBasic, nil
	case wurfl.HeaderQualityNone.String():
		return wurfl.HeaderQualityNone, nil
	}
	return wurfl.HeaderQualityNone, fmt.Errorf("GetHeaderQuality error")
}
func (m *MockWurfl) Destroy() {
}
func (m *MockWurfl) GetAllVCaps() []string {
	return []string{"is_ios", "is_app"}
}
func (m *MockWurfl) GetAllCaps() []string {
	return []string{"brand_name", "model_name"}
}
func (m *MockWurfl) GetInfo() string {
	return "the Wurfl info..."
}
func (m *MockWurfl) GetLastLoadTime() string {
	return "the Wurfl last load time"
}
func (m *MockWurfl) GetLastUpdated() string {
	return "the Wurfl last updated time"
}
func (m *MockWurfl) HasCapability(cap string) bool {
	for _, c := range m.GetAllCaps() {
		if cap == c {
			return true
		}
	}
	return false
}
func (m *MockWurfl) HasVirtualCapability(vcap string) bool {
	for _, vc := range m.GetAllVCaps() {
		if vcap == vc {
			return true
		}
	}
	return false
}

func (m *MockWurfl) SetAttr(attr int, value int) error {
	return nil
}

func (m *MockWurfl) GetAttr(attr int) (int, error) {
	return 0, nil
}

func TestWurfl_Mock(t *testing.T) {

	// ExampleService simulates a service that uses the Wurfl Handler
	type ExampleService struct {
		wurfl.WurflHandler
	}

	m := &MockWurfl{}

	s := &ExampleService{WurflHandler: m}

	t.Run("Test GetAllCaps mock", func(t *testing.T) {
		caps := s.GetAllCaps()
		assert.Equal(t, 2, len(caps))
	})

	t.Run("Test HasCapability mock", func(t *testing.T) {
		hasCap := s.HasCapability("brand_name")
		assert.Equal(t, true, hasCap)

		hasCap = s.HasCapability("not_exists")
		assert.Equal(t, false, hasCap)
	})

	t.Run("Test HasVirtualCapability mock", func(t *testing.T) {
		hasVCap := s.HasVirtualCapability("is_ios")
		assert.Equal(t, true, hasVCap)

		hasVCap = s.HasVirtualCapability("not_exists")
		assert.Equal(t, false, hasVCap)
	})
}
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.12.3.4", "1.12.3.4", 0},
		{"1.12.3.4", "1.12.3.5", -1},
		{"1.12.3.5", "1.12.3.4", 1},
		{"1.12.3.4", "2.0.0.0", -1},
		{"2.0.0.0", "1.2.3.4", 1},
		{"1.12.3.4", "1.13.0.0", -1},
		{"1.0.0.0", "1.0.0.1", -1},
		{"1.0.0.1", "1.0.0.0", 1},
		{"1.0.0.0", "1.0.0.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			result := wurfl.CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%s, %s) = %d; want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

// TestDownloadJira1236 check the issue detected from a customer, using capability filter in
// engine creation.
// It first attempts to create a WURFL engine with an evaluation version of the WURFL data,
// but this fails due to the capability filter that includes capabilities not present in the
// evaluation version. It then downloads a fresh copy of the WURFL data, and successfully
// creates a WURFL engine with the same capability filter.
func TestDownloadJira1236(t *testing.T) {
	// First create an engine with eval version ov wurfl.zip, using a capfilter
	// The capfilter will contain capabilities not inclueded in eval version, so
	// the engine creation will fail
	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// Copy evaluation wurfl.zip installed with libwurfl here
	// this file could be under /usr/share/wurfl or /usr/local/share/wurfl depending on the system
	srcPath := "/usr/share/wurfl/wurfl.zip"
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		srcPath = "/usr/local/share/wurfl/wurfl.zip"
	}
	destPath := "wurfl.zip"

	err := copyFile(srcPath, destPath)
	if err != nil {
		t.Fatalf("Failed to copy wurfl.zip: %v", err)
	}

	// Set a filter containing capabilities not included in eval version
	// Capability filtering is discouraged and will be deprecated. Here only for testing purposes

	capfilter := []string{
		"brand_name",
		"model_name",
		"resolution_height",
		"resolution_width",
		"is_bot",
		"release_date",
	}

	// Create the engine with eval version, the capfilter makes the engine creation fail
	_, err = wurfl.Create("wurfl.zip", nil, capfilter, -1, wurfl.WurflCacheProviderLru, "100000")

	if err == nil || err.Error() != "specified capability is missing" {
		t.Errorf("Create ok, it should fail\n")
	}

	// Now download first a fresh wurfl.zip, and then create the engine
	err = wurfl.Download(URL, ".")
	require.NoError(t, err)

	_, err = wurfl.Create("wurfl.zip", nil, capfilter, -1, wurfl.WurflCacheProviderLru, "100000")
	require.NoError(t, err)

	// Remove local wurfl.zip
	err = os.Remove("wurfl.zip")
	assert.NoError(t, err)
}

// TestDownload tests the Download function of the wurfl package. It creates a temporary
// directory, downloads the WURFL data file from the environment variable SM_UPDATER_DATA_URL,
// and checks for various error scenarios, including invalid URLs, invalid folders, and
// an empty URL.
func TestDownload(t *testing.T) {
	// get env var SM_UPDATER_DATA_URL URL value (from scientiamobile Vault)
	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// create a temporary directory for the test, removed after the test
	tempDir, err := os.MkdirTemp("", "wurfl_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// tests is a slice of test cases for the Download function.
	// Each test case has a name, a URL to download from, a folder to download to,
	// and an expected error value.
	// The tests cover various scenarios, including a valid download, an invalid URL,
	// an invalid folder, and an empty URL.
	tests := []struct {
		name        string
		URL         string
		folder      string
		expectedErr bool
	}{
		{
			name:        "Valid download",
			URL:         URL,
			folder:      tempDir,
			expectedErr: false,
		},
		{
			name:        "Invalid URL",
			URL:         "htt://invalid-URL.com/wurfl.zip",
			folder:      tempDir,
			expectedErr: true,
		},
		{
			name:        "Invalid folder",
			URL:         URL,
			folder:      "/nonexistent/folder",
			expectedErr: true,
		},
		{
			name:        "Empty URL",
			URL:         "",
			folder:      tempDir,
			expectedErr: true,
		},
		// Empty folder case: running it on TC will work, because unit tests are
		// run on dockers. It will fail on personal laptops, unless executed with
		// root privileges, because wurfl.zip will be downloaded under "/" dir
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wurfl.Download(tt.URL, tt.folder)
			if tt.expectedErr && err == nil {
				t.Errorf("Expected an error, but got nil for %s test", tt.name)
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v for test %s", err, tt.name)
			}
		})
	}
}

// TestDownloadAndLoad tests the Download function by checking that the file was,
// really downloaded, and if it can be used to create a WURFL engine.
func TestDownloadAndLoad(t *testing.T) {
	// get env var SM_UPDATER_DATA_URL URL value (from scientiamobile Vault)
	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// Creates a temporary directory for the WURFL integration test and defers its removal.
	// The temporary directory is used to store the downloaded WURFL data file during the test.
	tempDir, err := os.MkdirTemp("", "wurfl_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Download() downloads the WURFL data file from the specified URL and saves it to the specified directory.
	// If the download is successful, the function returns nil. Otherwise, it returns an error.
	err = wurfl.Download(URL, tempDir)
	require.NoError(t, err)

	// Check if the file was downloaded
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)

	if len(files) == 0 {
		t.Errorf("No files downloaded")
	}

	for _, file := range files {
		if file.Name() == "wurfl.zip" {
			wurflfile := tempDir + "/" + file.Name()
			// Try to create and load an engine with the freshly downloaded file
			_, err = wurfl.Create(wurflfile, nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
			require.NoError(t, err)
			return
		}
	}

	t.Errorf("wurfl.zip not found in downloaded files")
}
func TestGetStaticCap(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	tests := []struct {
		name     string
		cap      string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid capability",
			cap:      "mobile_browser_version",
			expected: "12.0",
			wantErr:  false,
		},
		{
			name:     "Another valid capability",
			cap:      "pointing_method",
			expected: "touchscreen",
			wantErr:  false,
		},
		{
			name:     "Boolean capability",
			cap:      "is_tablet",
			expected: "false",
			wantErr:  false,
		},
		{
			name:     "Non-existent capability",
			cap:      "non_existent_cap",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty capability string",
			cap:      "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := device.GetStaticCap(tt.cap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Device.GetStaticCap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Device.GetStaticCap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetStaticCapConcurrency(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := device.GetStaticCap("mobile_browser_version")
			if err != nil {
				t.Errorf("Concurrent GetStaticCap() failed: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestGetVirtualCap(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	tests := []struct {
		name     string
		vcap     string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid virtual capability",
			vcap:     "is_smartphone",
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "Non-existent virtual capability",
			vcap:     "non_existent_cap",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty virtual capability",
			vcap:     "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := device.GetVirtualCap(tt.vcap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Device.GetVirtualCap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("Device.GetVirtualCap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetVirtualCapConcurrency(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := device.GetVirtualCap("is_smartphone")
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestGetStaticCaps(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	t.Run("ValidCapabilities", func(t *testing.T) {
		caps := []string{"mobile_browser", "is_tablet", "pointing_method"}
		result, err := device.GetStaticCaps(caps)
		assert.NoError(t, err)

		if len(result) != len(caps) {
			t.Errorf("Expected %d capabilities, got %d", len(caps), len(result))
		}
		for _, cap := range caps {
			if _, ok := result[cap]; !ok {
				t.Errorf("Expected capability %s not found in result", cap)
			}
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		result, err := device.GetStaticCaps([]string{})
		assert.NoError(t, err)
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d capabilities", len(result))
		}
	})

	t.Run("NonExistentCapability", func(t *testing.T) {
		caps := []string{"non_existent_cap"}
		result, err := device.GetStaticCaps(caps)
		assert.NotNil(t, err)    // error returned
		assert.NotNil(t, result) // empty map returned
	})

	t.Run("MixedCapabilities", func(t *testing.T) {
		caps := []string{"mobile_browser", "non_existent_cap", "is_tablet"}
		result, err := device.GetStaticCaps(caps)
		assert.NotNil(t, err)
		assert.NotNilf(t, result, "Expected valid map for mixed capabilities, got nil")

		// We expect wrong capability to be skipped
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})

	t.Run("DuplicateCaps", func(t *testing.T) {
		caps := []string{"mobile_browser", "mobile_browser", "is_tablet"}
		result, err := device.GetStaticCaps(caps)
		assert.NoError(t, err)
		// We expect one of the duplicate capabilities to be removed
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})
}

func TestGetVirtualCaps(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)
	require.NotNil(t, wengine)
	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	require.NoError(t, err)

	t.Run("ValidVCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "advertised_device_os", "pixel_density"}
		result, err := device.GetVirtualCaps(caps)
		assert.NoError(t, err)
		if len(result) != len(caps) {
			t.Errorf("Expected %d capabilities, got %d", len(caps), len(result))
		}
		for _, cap := range caps {
			if _, ok := result[cap]; !ok {
				t.Errorf("Expected capability %s not found in result", cap)
			}
		}
	})

	t.Run("EmptyMap", func(t *testing.T) {
		result, err := device.GetVirtualCaps([]string{})
		assert.NoError(t, err)
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d capabilities", len(result))
		}
	})

	t.Run("NonExistentVCap", func(t *testing.T) {
		caps := []string{"non_existent_vcap"}
		result, err := device.GetVirtualCaps(caps)
		assert.NotNil(t, err)    // error returned
		assert.NotNil(t, result) // empty map returned
	})

	t.Run("MixedVCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "non_existent_cap", "pixel_density"}
		result, err := device.GetVirtualCaps(caps)
		assert.NotNil(t, err)
		assert.NotNilf(t, result, "Expected valid map for mixed capabilities, got nil")
		// We expect wrong virtual capability to be skipped
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})

	t.Run("DuplicateCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "advertised_device_os", "is_smartphone"}
		result, err := device.GetVirtualCaps(caps)
		assert.NoError(t, err)
		// We expect one of the duplicate capabilities to be removed
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})
}

func TestGetLastUpdated(t *testing.T) {
	// get env var SM_UPDATER_DATA_URL URL value (from scientiamobile Vault)
	URL := os.Getenv("SM_UPDATER_DATA_URL")
	if URL == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// Download() downloads the WURFL data file from the specified URL and saves it to the specified directory.
	// If the download is successful, the function returns nil. Otherwise, it returns an error.
	err := wurfl.Download(URL, ".")
	assert.NoError(t, err)

	wengine, err := wurfl.Create("./wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
	require.NoError(t, err)
	defer wengine.Destroy()

	lastUpdated := wengine.GetLastUpdated()
	if lastUpdated == "" {
		t.Errorf("Got empty string as wurfl last updated time")
	}

	t.Logf("Last updated: %s", lastUpdated)

	// Remove the local wurfl.zip file
	err = os.Remove("./wurfl.zip")
	assert.NoError(t, err)
}
