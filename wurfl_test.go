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
	if err != nil {
		t.Errorf("Create error : %s\n", err.Error())
		return nil
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
	if err != nil {
		t.Errorf("Create error : %s\n", err.Error())
		return nil
	}

	return wengine
}

func TestWurfl_WurflGetAPIVersion(t *testing.T) {
	ver := wurfl.APIVersion()
	if ver == "" {
		t.Errorf("APIVersion error, ret nil\n")
	}
	fmt.Printf("WURFL API Version: %s\n", ver)
}

func TestWurfl_Create(t *testing.T) {
	ua := "ArtDeviant/3.0.2 CFNetwork/711.3.18 Darwin/14.0.0"
	var device *wurfl.Device
	var deviceid string
	var err error

	wengine := fixtureCreateEngine(t)

	fmt.Printf("WURFL API Version: %s\n", wengine.GetAPIVersion())

	device, err = wengine.LookupUserAgent(ua)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}

	deviceid, err = device.GetDeviceID()

	if err == nil {
		if deviceid != "apple_iphone_ver8_3_subuacfnetwork" {
			t.Errorf("Lookup mismatch ? >%s< instead of >apple_iphone_ver8_3_subuacfnetwork<", deviceid)
		}
	}
	// uncle firefox

	wengine.Destroy()
}

func TestWurfl_Lookup(t *testing.T) {
	var wengine *wurfl.Wurfl
	var err error
	var err2 error
	var device *wurfl.Device
	var newdevice *wurfl.Device
	var deviceid string
	var deviceid2 string
	var verbose bool = false

	capfilter := []string{
		"mobile_browser_version",
		"pointing_method",
		"is_tablet",
	}

	caps := capfilter[0:3]

	if verbose {
		fmt.Println("Loading engine ...")
	}

	_, oserr := os.Stat("/usr/local/share/wurfl/wurfl.zip")
	if oserr == nil {
		// macosx rootless
		wengine, err = wurfl.Create("/usr/local/share/wurfl/wurfl.zip", nil, caps, -1, wurfl.WurflCacheProviderLru, "100000")
	} else {
		// all other systems (TODO windows)
		wengine, err = wurfl.Create("/usr/share/wurfl/wurfl.zip", nil, caps, -1, wurfl.WurflCacheProviderLru, "100000")
	}

	if err != nil {
		// some error here
		t.Errorf("Create error : %s\n", err.Error())
	}

	wengine.SetLogPath("api.log")

	// fmt.Println("Engine loaded, version : ", wengine.GetAPIVersion(), "wurfl info ", wengine.GetInfo())

	ua := "ArtDeviant/3.0.2 CFNetwork/711.3.18 Darwin/14.0.0"
	// ua := "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:48.0) Gecko/20100101 Firefox/48.0"
	// for i:= 0; i < 1000000; i++ {
	if verbose {
		fmt.Println("Lookup of ", ua)
	}

	device, err = wengine.LookupUserAgent(ua)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}

	deviceid, err = device.GetDeviceID()
	if err != nil {
		t.Errorf("GetDeviceID error : %s\n", err.Error())
	}

	// fmt.Println("wurfl_id : ", deviceid)

	newdevice, err = wengine.LookupDeviceID(deviceid)
	if err != nil {
		t.Errorf("LookupDeviceID error : %s\n", err.Error())
	}

	deviceid2, err2 = newdevice.GetDeviceID()
	if err2 != nil {
		t.Errorf("GetDeviceID error : %s\n", err.Error())
	}

	if deviceid != deviceid2 {
		t.Errorf("Error, devices do not match %s, %s", deviceid, deviceid2)
	}

	_, uaerr := device.GetUserAgent()
	if uaerr != nil {
		t.Errorf("GetUserAgent error : %s\n", err.Error())
	}

	oua, uaerr := device.GetOriginalUserAgent()
	if uaerr != nil {
		t.Errorf("GetOriginalUserAgent error : %s\n", err.Error())
	}

	if oua != ua {
		t.Errorf("Error, ua matched >%s< and device original ua >%s< do not match", ua, oua)
		// fmt.Println("Error, ua matched (", ua, ") and matched device ua (", dua, ") do not match")
	}

	nua, uaerr := device.GetNormalizedUserAgent()
	if uaerr != nil {
		t.Errorf("GetNormalizedUserAgent error : %s\n", err.Error())
	}

	if len(nua) >= len(ua) {
		if verbose {
			fmt.Printf("Error, Normalized ua >%s< longer than original ua >%s<", nua, ua)
		}
	}

	if device.IsRoot() {
		fmt.Printf("Device is root\n")
	}

	if device.GetCapability("mobile_browser_version") != "8.0" {
		t.Errorf("device.GetCapability(\"mobile_browser_version\") does not return 8.0 : %s\n", device.GetCapability("mobile_browser_version"))
	}

	vcaps := make(map[string]string)
	vcaps = device.GetVirtualCapabilities(wengine.GetAllVCaps())

	if vcaps["advertised_device_os"] != "iOS" {
		t.Errorf("device.GetVirtualCapabilities() \"advertised_device_os\" != \"iOS\" : %s\n", vcaps["advertised_device_os"])
	}

	allcaps := make(map[string]string)
	allcaps = device.GetCapabilities(wengine.GetAllCaps())

	if allcaps["device_os"] != "iOS" {
		t.Errorf("device.GetCapabilities() \"device_os\" != \"iOS\" : %s\n", allcaps["device_os"])
	}

	device.GetMatchType()
	device.GetRootID()

	if verbose {
		fmt.Println(deviceid)
		fmt.Println(device.GetCapability("mobile_browser_version"))
		fmt.Println(device.GetVirtualCapability("is_android"))
		fmt.Println(device.GetCapabilities(capfilter))
		fmt.Println(device.GetMatchType())
		fmt.Println(device.GetRootID())
	}

	device.Destroy()
	newdevice.Destroy()
	wengine.Destroy()
	fi, apilogoserr := os.Stat("api.log")
	assert.Nil(t, apilogoserr)
	assert.NotNil(t, fi)
	// os.Remove("api.log")
}

func TestWurfl_GetAllVCaps(t *testing.T) {

	wengine := fixtureCreateEngine(t)
	s := wengine.GetAllVCaps()
	if len(s) != 28 {
		t.Errorf("Vcaps should be 28, they are %d", len(s))
	}

	wengine.Destroy()
}

// Test_GetCapability : various cases on GetCapability / GetVirtualCapability
func Test_GetCapability(t *testing.T) {
	assert := assert.New(t)
	wengine := fixtureCreateEngine(nil)
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
	wengine.Destroy()
}

func Test_LookupRequest(t *testing.T) {

	var err error
	// var deviceLA *wurfl.Device
	// var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// fmt.Println(wengine.ImportantHeaderNames)

	// lookup both UAs. When using both headers in LookupRequest, the XUCBrowserDeviceUA has precedence.
	UserAgentDevice, err := wengine.LookupUserAgent(UserAgent)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}
	UserAgentDeviceId, _ := UserAgentDevice.GetDeviceID()

	XUCBrowserDeviceDevice, err := wengine.LookupUserAgent(XUCBrowserDeviceUA)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}
	XUCBrowserDeviceDeviceId, _ := XUCBrowserDeviceDevice.GetDeviceID()

	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-UCBrowser-Device-UA", XUCBrowserDeviceUA)

	reqDevice, err := wengine.LookupRequest(req)
	if err != nil {
		t.Errorf("LookupRequest error : %s\n", err.Error())
	}
	reqDeviceId, _ := reqDevice.GetDeviceID()

	// now verify that device retrieved with headers is the the same as XUCBrowserDeviceDevice
	if UserAgentDeviceId == reqDeviceId {
		t.Errorf("Devices are the same, should be different : %s, %s\n", UserAgentDeviceId, reqDeviceId)
	}

	if XUCBrowserDeviceDeviceId != reqDeviceId {
		t.Errorf("Devices are different, should be the same: %s, %s\n", XUCBrowserDeviceDeviceId, reqDeviceId)
	}

	UserAgentDevice.Destroy()
	XUCBrowserDeviceDevice.Destroy()
	reqDevice.Destroy()
	wengine.Destroy()
}

// Test_LookupRequestExperimental : test new sec-ch headers
func Test_LookupRequestExperimental(t *testing.T) {

	var err error
	// var deviceLA *wurfl.Device
	// var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// set experimental headers
	wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, 1)

	// User-Agent
	UserAgent := "Mozilla/5.0 (Linux; Android 11; SM-M315F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.91 Mobile Safari/537.36"

	UserAgentDevice, err := wengine.LookupUserAgent(UserAgent)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}
	UserAgentDeviceId, _ := UserAgentDevice.GetDeviceID()

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
	if err != nil {
		t.Errorf("LookupRequest error : %s\n", err.Error())
	}
	reqDeviceId, _ := reqDevice.GetDeviceID()

	// now verify that device retrieved with headers is the the same as XUCBrowserDeviceDevice
	if UserAgentDeviceId != reqDeviceId {
		t.Errorf("Devices should be the same : %s, %s\n", UserAgentDeviceId, reqDeviceId)
	}

	UserAgentDevice.Destroy()
	reqDevice.Destroy()
	wengine.Destroy()
}

func Test_LookupWithImportantHeaderMap(t *testing.T) {
	var err error
	var deviceLA *wurfl.Device
	var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// do a LookupUserAgent() with UA

	// fmt.Println(wengine.ImportantHeaderNames)

	deviceLA, err = wengine.LookupUserAgent(UserAgent)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}
	LaDeviceId, _ := deviceLA.GetDeviceID()

	// create IHMap and lookup using headers

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["X-UCBrowser-Device-UA"] = XUCBrowserDeviceUA

	deviceIHM, err = wengine.LookupWithImportantHeaderMap(IHMap)
	if err != nil {
		t.Errorf("LookupWithImportantHeaderMap error : %s\n", err.Error())
	}
	IHMDeviceId, _ := deviceIHM.GetDeviceID()

	if LaDeviceId == IHMDeviceId {
		t.Errorf("Devices are the same, should be different : %s, %s\n", LaDeviceId, IHMDeviceId)
	}

	deviceLA.Destroy()
	deviceIHM.Destroy()
	wengine.Destroy()
}

// support issue #3267, INFUZE-1053
/*
original support function (does not compile as-is)
func lookUpUserAgent(ua string, capabilities []string) map[string]string {
	device, werr := wengine.LookupUserAgent(ua)
	defer device.Destroy()
	PrintError(werr)
	return device.GetCapabilities(capabilities)
}
*/

func lookUpUserAgent(wengine *wurfl.Wurfl, ua string, capabilities []string) map[string]string {
	device, _ := wengine.LookupUserAgent(ua)
	defer device.Destroy()
	return device.GetCapabilities(capabilities)
}

func TestJira_INFUZE1053(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	capabilitiesVector := []string{"max_image_width", "dual_orientation", "density_class", "is_ios"}
	capabilitiesMap := lookUpUserAgent(wengine, "Callpod Keeper for Android 1.0 (10.0.0/234) Dalvik/2.1.0 (Linux; U; Android 5.0.1; SAMSUNG-SGH-I337 Build/LRX22C)", capabilitiesVector)
	fmt.Println(capabilitiesMap)
	wengine.Destroy()
}

func Test_LookupDeviceIDWithImportantHeaderMap(t *testing.T) {
	var err error
	var deviceIDIHM *wurfl.Device
	var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// do a LookupUserAgent() with UA

	// fmt.Println(wengine.ImportantHeaderNames)

	// create IHMap and lookup using headers

	IHMap := make(map[string]string)
	IHMap["User-Agent"] = UserAgent
	IHMap["X-UCBrowser-Device-UA"] = XUCBrowserDeviceUA

	deviceIHM, err = wengine.LookupWithImportantHeaderMap(IHMap)
	if err != nil {
		t.Errorf("LookupWithImportantHeaderMap error : %s\n", err.Error())
	}
	IHMDeviceId, _ := deviceIHM.GetDeviceID()
	AdvBrow1 := deviceIHM.GetVirtualCapability("advertised_browser")

	// now lookup by deviceID and no header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceID(IHMDeviceId)
	if err != nil {
		t.Errorf("LookupDeviceID error : %s\n", err.Error())
	}
	AdvBrow2 := deviceIDIHM.GetVirtualCapability("advertised_browser")

	if AdvBrow1 == AdvBrow2 {
		t.Errorf("advertised_browser are the same, should be different : %s, %s\n", AdvBrow1, AdvBrow2)
	}

	// now lookup by deviceID and header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceIDWithImportantHeaderMap(IHMDeviceId, IHMap)
	if err != nil {
		t.Errorf("LookupDeviceIDWithImportantHeaderMap error : %s\n", err.Error())
	}
	AdvBrow3 := deviceIDIHM.GetVirtualCapability("advertised_browser")

	if AdvBrow1 != AdvBrow3 {
		t.Errorf("advertised_browser are different, should be the same: %s, %s\n", AdvBrow1, AdvBrow3)
	}

	deviceIHM.Destroy()
	wengine.Destroy()
}

func Test_LookupWithImportantHeaderMapCaseInsensitive(t *testing.T) {
	var err error
	var deviceLA *wurfl.Device
	var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// User-Agent
	UserAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"

	// do a LookupUserAgent() with UA

	// fmt.Println(wengine.ImportantHeaderNames)

	deviceLA, err = wengine.LookupUserAgent(UserAgent)
	if err != nil {
		t.Errorf("LookupuserAgent error : %s\n", err.Error())
	}
	LaDeviceId, _ := deviceLA.GetDeviceID()

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

	deviceIHM, err = wengine.LookupWithImportantHeaderMap(IHMap)
	if err != nil {
		t.Errorf("LookupWithImportantHeaderMap error : %s\n", err.Error())
	}
	IHMDeviceId, _ := deviceIHM.GetDeviceID()

	if LaDeviceId == IHMDeviceId {
		t.Errorf("Devices are the same, should be different : %s, %s\n", LaDeviceId, IHMDeviceId)
	}

	deviceLA.Destroy()
	deviceIHM.Destroy()
	wengine.Destroy()
}

func Test_LookupDeviceIDWithRequest(t *testing.T) {
	var err error
	var deviceIDIHM *wurfl.Device
	var deviceIHM *wurfl.Device

	wengine := fixtureCreateEngine(t)

	// User-Agent
	UserAgent := "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCBrowser/8.7.0.218 U2/1.0.0 Mobile"
	// X-UCBrowser-Device-UA
	XUCBrowserDeviceUA := "Mozilla/5.0 (Linux; U; Android 5.1.1; en-US; SM-J200G Build/LMY47X) AppleWebKit/528.5+ (KHTML, like Gecko) Version/3.1.2 Mobile Safari/525.20.1"

	// create http.Request and lookup using headers
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("X-UCBrowser-Device-UA", XUCBrowserDeviceUA)

	deviceIHM, err = wengine.LookupRequest(req)
	if err != nil {
		t.Errorf("LookupRequest error : %s\n", err.Error())
	}
	IHMDeviceId, _ := deviceIHM.GetDeviceID()
	AdvBrow1 := deviceIHM.GetVirtualCapability("advertised_browser")

	// now lookup by deviceID and no header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceID(IHMDeviceId)
	if err != nil {
		t.Errorf("LookupDeviceID error : %s\n", err.Error())
	}
	AdvBrow2 := deviceIDIHM.GetVirtualCapability("advertised_browser")

	if AdvBrow1 == AdvBrow2 {
		t.Errorf("advertised_browser are the same, should be different : %s, %s\n", AdvBrow1, AdvBrow2)
	}

	// now lookup by deviceID and header and check that an advertised vcap behaves correctly
	deviceIDIHM, err = wengine.LookupDeviceIDWithRequest(IHMDeviceId, req)
	if err != nil {
		t.Errorf("LookupDeviceIDWithRequest error : %s\n", err.Error())
	}
	AdvBrow3 := deviceIDIHM.GetVirtualCapability("advertised_browser")

	if AdvBrow1 != AdvBrow3 {
		t.Errorf("advertised_browser are different, should be the same: %s, %s\n", AdvBrow1, AdvBrow3)
	}

	deviceIHM.Destroy()
	wengine.Destroy()
}

func setTimeBackForFile(filename string, days int) {
	file, _ := os.Stat(filename)

	nowfile := file.ModTime()

	then := nowfile.AddDate(0, 0, days)

	_ = os.Chtimes(filename, then, then)
}

func TestWurfl_UpdaterRunonce(t *testing.T) {
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

	// set modification time to 1 month before
	setTimeBackForFile("/tmp/wurfl.zip", -30)

	info, _ := os.Stat("/tmp/wurfl.zip")
	mdt1 := info.ModTime()

	wengine, err = wurfl.Create("/tmp/wurfl.zip", nil, nil, wurfl.WurflEngineTargetHighAccuray, wurfl.WurflCacheProviderDoubleLru, "100000")
	if err != nil {
		t.Errorf("Create error : %s\n", err.Error())
	}

	// set env var SM_UPDATER_DATA_URL to your updater url (from scientiamobile Vault)

	Url := os.Getenv("SM_UPDATER_DATA_URL")
	if Url == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// set logger file
	_ = wengine.SetUpdaterLogPath("/tmp/wurfl-updater-log.txt")

	// set updater path
	uerr := wengine.SetUpdaterDataURL(Url)
	if uerr != nil {
		t.Errorf("SetUpdaterDataUrl returned : %s\n", uerr.Error())
	}

	// set updater user-agent
	expUa := fmt.Sprintf("golang_wurfl_test/%s", wurfl.Version)

	t.Logf("Specific golang binding User-Agent string: %s", expUa)

	uerr = wengine.SetUpdaterUserAgent(expUa)
	if uerr != nil {
		t.Errorf("SetUpdaterUserAgent returned : %s\n", uerr.Error())
	}

	ua := wengine.GetUpdaterUserAgent()
	if ua == "" {
		t.Errorf("GetUpdaterUserAgent returned empty string, expected %s", expUa)
	}

	// set timeout to defaults
	uerr = wengine.SetUpdaterDataURLTimeout(-1, -1)

	uerr = wengine.UpdaterRunonce()
	if uerr != nil {
		t.Errorf("UpdaterRunonce returned : %s\n", uerr.Error())
	}

	// check if the modification time of wurfl.zip has changed
	info, _ = os.Stat("/tmp/wurfl.zip")
	mdt2 := info.ModTime()

	if mdt1.Equal(mdt2) {
		t.Errorf("/tmp/wurfl.zip not downloaded\n")
	}

	wengine.Destroy()
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

	// set modification time to 1 month before
	setTimeBackForFile("/tmp/wurfl.zip", -30)

	info, _ := os.Stat("/tmp/wurfl.zip")
	mdt1 := info.ModTime()

	wengine, err = wurfl.Create("/tmp/wurfl.zip", nil, nil, wurfl.WurflEngineTargetHighAccuray, wurfl.WurflCacheProviderDoubleLru, "100000")
	if err != nil {
		t.Errorf("Create error : %s\n", err.Error())
	}

	_ = wengine.SetUpdaterLogPath("/tmp/wurfl-updater-log.txt")

	lastLoadTime := wengine.GetLastLoadTime()

	// set env var SM_UPDATER_DATA_URL to your updater url (from scientiamobile Vault)

	Url := os.Getenv("SM_UPDATER_DATA_URL")
	if Url == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// set updater path
	uerr := wengine.SetUpdaterDataURL(Url)
	if uerr != nil {
		t.Errorf("SetUpdaterDataUrl returned : %s\n", uerr.Error())
	}

	uerr = wengine.UpdaterStart()
	if uerr != nil {
		t.Errorf("UpdaterStart returned : %s\n", uerr.Error())
	}

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

	wengine.Destroy()
}

func TestWurfl_Getters(t *testing.T) {
	wengine := fixtureCreateEngine(t)

	s := wengine.GetLastLoadTime()
	if len(s) < 5 {
		t.Errorf("Wrong last load time %s", s)
	}

	i := wengine.GetInfo()
	if len(i) < 5 {
		t.Errorf("Wrong wurfl info %s", i)
	}

	a := wengine.GetInfo()
	if len(a) < 3 {
		t.Errorf("Wrong wurfl info %s", a)
	}

	e := wengine.GetEngineTarget()
	if len(e) < 6 {
		t.Errorf("Wrong engine target %s", e)
	}

	wengine.SetUserAgentPriority(wurfl.WurflUserAgentPriorityUsePlainUserAgent)

	p := wengine.GetUserAgentPriority()
	if len(p) < 4 {
		t.Errorf("Wrong Useragent Priority %s", p)
	}

	if wengine.HasVirtualCapability("pippo") != false {
		t.Errorf("HasVirtualCapability failure")
	}

	if wengine.HasVirtualCapability("is_ios") != true {
		t.Errorf("HasVirtualCapability failure")
	}

	if wengine.HasCapability("device_os") != true {
		t.Errorf("HasCapability failure")
	}

	wengine.Destroy()
}

func TestWurfl_GetAllDeviceIds(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	ids := wengine.GetAllDeviceIds()
	if len(ids) < 50000 {
		t.Errorf("Not all device ids have been loaded, they are %d", len(ids))
	}

	wengine.Destroy()
}

// Test_SetAttr :
func Test_SetAttr(t *testing.T) {

	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)

	t.Run("setattr", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, 10)

		if err != nil {
			t.Errorf("SetAttr returns error %s", err.Error())
		}

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)

		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != 10 {
			t.Errorf("Wrong attr value : %d instead of 10", attrValue)
		}
	})

	t.Run("setattr negative value", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrExtraHeadersExperimental, -10)

		if err != nil {
			t.Errorf("SetAttr returns error %s", err.Error())
		}

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)

		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != -10 {
			t.Errorf("Wrong attr value : %d instead of 10", attrValue)
		}
	})

	t.Run("setattr invalid attr", func(t *testing.T) {
		err = wengine.SetAttr(44, 10)

		if err == nil {
			t.Errorf("SetAttr doesn't return error but it should")
		}
	})

	wengine.Destroy()
}

func Test_SetAttr_FallbackCache(t *testing.T) {
	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)

	t.Run("setAttr fallback cache - default", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		if err != nil {
			t.Errorf("SetAttr returns error %s", err.Error())
		}

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDefault {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		}
	})

	t.Run("setAttr fallback cache - disabled", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
		if err != nil {
			t.Errorf("SetAttr returns error %s", err.Error())
		}

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDisabled {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
		}
	})

	t.Run("setAttr fallback cache - limited", func(t *testing.T) {
		err = wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		if err != nil {
			t.Errorf("SetAttr returns error %s", err.Error())
		}

		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheLimited {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		}

		// check that a new set overwrites the old one
		wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrCapabilityFallbackCache)
		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}
		if attrValue != wurfl.WurflAttrCapabilityFallbackCacheDefault {
			t.Errorf("GetAttr returns %d, but %d was expected", attrValue, wurfl.WurflAttrCapabilityFallbackCacheLimited)
		}

	})
	wengine.Destroy()
}

// Test_GetAttr :
func Test_GetAttr(t *testing.T) {

	var err error
	var attrValue int

	wengine := fixtureCreateEngine(t)

	t.Run("getattr", func(t *testing.T) {
		attrValue, err = wengine.GetAttr(wurfl.WurflAttrExtraHeadersExperimental)

		if err != nil {
			t.Errorf("GetAttr returns error %s", err.Error())
		}

		if attrValue != 1 {
			t.Errorf("Wrong attr value : %d instead of 1", attrValue)
		}

	})

	t.Run("getattr invalid attr", func(t *testing.T) {
		attrValue, err = wengine.GetAttr(44)

		if err == nil {
			t.Errorf("GetAttr doesn't return error but it should")
		}
	})

	wengine.Destroy()
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

func Benchmark_GetCapability_FallbackCacheDefault_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName := device.GetCapability("model_name")
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
		device.GetCapability("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheDefault_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv := device.GetCapability("is_smarttv")
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

	if fallbackMode != wurfl.WurflAttrCapabilityFallbackCacheDefault {
		b.Error("fallback cache mode should be DEFAULT")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		device.GetCapability("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheLimited_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName := device.GetCapability("model_name")
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
		device.GetCapability("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheLimited_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv := device.GetCapability("is_smarttv")
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
		device.GetCapability("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheDisabled_ModelName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	modelName := device.GetCapability("model_name")
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
		device.GetCapability("model_name")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheDisabled_IsSmarttv(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")
	isSmarttv := device.GetCapability("is_smarttv")
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
		device.GetCapability("is_smarttv")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheDisabled_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDisabled)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName := device.GetCapability("marketing_name")
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
		device.GetCapability("marketing_name")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheDefault_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheDefault)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName := device.GetCapability("marketing_name")
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
		device.GetCapability("marketing_name")
	}

	b.StopTimer()
}

func Benchmark_GetCapability_FallbackCacheLimited_MarketingName(b *testing.B) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	wengine.SetAttr(wurfl.WurflAttrCapabilityFallbackCache, wurfl.WurflAttrCapabilityFallbackCacheLimited)
	device, err := wengine.LookupDeviceID("samsung_sm_g990b_ver1")
	mktName := device.GetCapability("marketing_name")
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
		device.GetCapability("marketing_name")
	}

	b.StopTimer()
}

// Benchmark_GetCapability : test time in single const get_capability
func Benchmark_GetCapability(b *testing.B) {
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
		device.GetCapability("device_os")
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

func TestDevice_GetCapabilityAsInt(t *testing.T) {
	wengine := fixtureCreateEngine(nil)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	t.Run("Test GetCapabilityAsInt OK", func(t *testing.T) {
		cap, err := device.GetCapabilityAsInt("resolution_height")
		if err != nil {
			t.Error("error should have been nil, while is: " + err.Error())
		}
		if cap <= 0 {
			t.Error("expected a positive value for resolution_height")
		}
	})

	t.Run("Test GetCapabilityAsInt calling a non int capability, must return a not nil error", func(t *testing.T) {
		capname := "brand_name"
		// from 1.12.7.1 libwurfl returns error when asked for non numeric capabilities (ie: brand_name)
		_, err := device.GetCapabilityAsInt(capname)
		if err == nil {
			t.Error("error should have been not nil for non numeric capabilities calls")
		}
	})

	t.Run("Test GetCapabilityAsInt calling a capability using an empty string, must return a not nil error", func(t *testing.T) {
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetCapabilityAsInt("")
		if err == nil {
			t.Error("error should have been not nil for capabilities calls using empty string as name")
		}
	})
}

func TestDevice_GetVirtualCapabilityAsInt(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	defer wengine.Destroy()
	device, err := wengine.LookupDeviceID("google_pixel_5_ver1")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	defer device.Destroy()

	t.Run("Test GetVirtualCapabilityAsInt OK", func(t *testing.T) {
		cap, err := device.GetVirtualCapabilityAsInt("pixel_density")
		if err != nil {
			t.Error("error should have been nil, while is: " + err.Error())
		}
		if cap <= 0 {
			t.Error("expected a positive value for pixel_density")
		}
	})

	t.Run("Test GetVirtualCapabilityAsInt calling a non int virtual capability, must return a not nil error", func(t *testing.T) {
		capname := "form_factor"
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetVirtualCapabilityAsInt(capname)
		if err == nil {
			t.Error("error should have been not nil for non numeric virtual capabilities calls")
		}
	})

	t.Run("Test GetVirtualCapabilityAsInt calling a virtual capability using an empty string, must return a not nil error", func(t *testing.T) {
		// from 1.12.7.1 libwurfl returns error when asked for non numeric
		// virtual capabilities (ie: form_factor)
		_, err := device.GetVirtualCapabilityAsInt("")
		if err == nil {
			t.Error("error should have been not nil for virtual capabilities calls using empty string as name")
		}
	})
}

func TestDevice_GetRootID(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	defer wengine.Destroy()

	/**\
	* This device have no device roots in its fall back tree,
	* since no devices above it (itself included) are real devices (actual device roots),
	* in this case "" is expected.
	**/
	device, err := wengine.LookupDeviceID("generic_android_ver11_0_subff102_tablet")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	rootId := device.GetRootID()
	expectedRootId := ""
	if rootId != expectedRootId {
		t.Errorf("Expected rootId %s got %s", expectedRootId, rootId)
	}

	// generic has empty root id
	device, err = wengine.LookupDeviceID("generic")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	rootId = device.GetRootID()
	expectedRootId = ""
	if rootId != expectedRootId {
		t.Errorf("Expected rootId %s got %s", expectedRootId, rootId)
	}

	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1_subu3k10")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	rootId = device.GetRootID()
	expectedRootId = "natec_smart_tv_dongle_hd221_ver1"
	if rootId != expectedRootId {
		t.Errorf("Expected rootId %s got %s", expectedRootId, rootId)
	}

	//is an actual device root , root is itself
	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	rootId = device.GetRootID()
	expectedRootId = "natec_smart_tv_dongle_hd221_ver1"
	if rootId != expectedRootId {
		t.Errorf("Expected rootId %s got %s", expectedRootId, rootId)
	}
}

func TestDevice_GetParentID(t *testing.T) {
	wengine := fixtureCreateEngine(t)
	defer wengine.Destroy()

	device, err := wengine.LookupDeviceID("generic_android_ver11_0_subff102_tablet")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	parentId := device.GetParentID()
	expectedParentId := "generic_android_ver11_0_subff101_tablet"
	if parentId != expectedParentId {
		t.Errorf("Expected parentId %s got %s", expectedParentId, parentId)
	}

	// generic has empty parent id
	device, err = wengine.LookupDeviceID("generic")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	parentId = device.GetParentID()
	expectedParentId = ""
	if parentId != expectedParentId {
		t.Errorf("Expected parentId %s got %s", expectedParentId, parentId)
	}

	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1_subu3k10")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	parentId = device.GetParentID()
	expectedParentId = "natec_smart_tv_dongle_hd221_ver1"
	if parentId != expectedParentId {
		t.Errorf("Expected parentId %s got %s", expectedParentId, parentId)
	}

	//is an actual device root , but has a parent
	device, err = wengine.LookupDeviceID("natec_smart_tv_dongle_hd221_ver1")

	if err != nil {
		t.Error("error from LookupDeviceID should be the nil")
	}
	if device == nil {
		t.Error("device from LookupDeviceID should be the NOT nil")
	}

	parentId = device.GetParentID()
	expectedParentId = "generic_android_ver4_2"
	if parentId != expectedParentId {
		t.Errorf("Expected parentId %s got %s", expectedParentId, parentId)
	}
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

func (m *MockWurfl) Download(url string, folder string) error {
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
		expected := 2
		caps := s.GetAllCaps()
		if len(caps) != expected {
			t.Errorf("Expected %q, got %q", expected, caps)
		}
	})

	t.Run("Test HasCapability mock", func(t *testing.T) {
		expected := true
		hasCap := s.HasCapability("brand_name")
		if hasCap != expected {
			t.Errorf("Expected %t, got %t", expected, hasCap)
		}

		expected = false
		hasCap = s.HasCapability("not_exists")
		if hasCap != expected {
			t.Errorf("Expected %t, got %t", expected, hasCap)
		}
	})

	t.Run("Test HasVirtualCapability mock", func(t *testing.T) {
		expected := true
		hasVCap := s.HasVirtualCapability("is_ios")
		if hasVCap != expected {
			t.Errorf("Expected %t, got %t", expected, hasVCap)
		}

		expected = false
		hasVCap = s.HasVirtualCapability("not_exists")
		if hasVCap != expected {
			t.Errorf("Expected %t, got %t", expected, hasVCap)
		}
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
	url := os.Getenv("SM_UPDATER_DATA_URL")
	if url == "" {
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
	err = wurfl.Download(url, ".")
	if err != nil {
		t.Errorf("Wurfl.zip download failed: %v", err)
	}

	_, err = wurfl.Create("wurfl.zip", nil, capfilter, -1, wurfl.WurflCacheProviderLru, "100000")

	if err != nil {
		t.Errorf("Wurfl engine creation failed: %v\n", err)
	}

	// Remove local wurfl.zip
	err = os.Remove("wurfl.zip")

	if err != nil {
		t.Errorf("Remove wurfl.zip failure: %v\n", err)
	}
}

// TestDownload tests the Download function of the wurfl package. It creates a temporary
// directory, downloads the WURFL data file from the environment variable SM_UPDATER_DATA_URL,
// and checks for various error scenarios, including invalid URLs, invalid folders, and
// an empty URL.
func TestDownload(t *testing.T) {
	// get env var SM_UPDATER_DATA_URL url value (from scientiamobile Vault)
	url := os.Getenv("SM_UPDATER_DATA_URL")
	if url == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// create a temporary directory for the test, removed after the test
	tempDir, err := os.MkdirTemp("", "wurfl_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// tests is a slice of test cases for the Download function.
	// Each test case has a name, a URL to download from, a folder to download to,
	// and an expected error value.
	// The tests cover various scenarios, including a valid download, an invalid URL,
	// an invalid folder, and an empty URL.
	tests := []struct {
		name        string
		url         string
		folder      string
		expectedErr bool
	}{
		{
			name:        "Valid download",
			url:         url,
			folder:      tempDir,
			expectedErr: false,
		},
		{
			name:        "Invalid URL",
			url:         "https://invalid-url.com/wurfl.zip",
			folder:      tempDir,
			expectedErr: true,
		},
		{
			name:        "Invalid folder",
			url:         url,
			folder:      "/nonexistent/folder",
			expectedErr: true,
		},
		{
			name:        "Empty URL",
			url:         "",
			folder:      tempDir,
			expectedErr: true,
		},
		// Empty folder case: running it on TC will work, because unit tests are
		// run on dockers. It will fail on personal laptops, unless executed with
		// root privileges, because wurfl.zip will be downloaded under "/" dir
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wurfl.Download(tt.url, tt.folder)
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
	// get env var SM_UPDATER_DATA_URL url value (from scientiamobile Vault)
	url := os.Getenv("SM_UPDATER_DATA_URL")
	if url == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// Creates a temporary directory for the WURFL integration test and defers its removal.
	// The temporary directory is used to store the downloaded WURFL data file during the test.
	tempDir, err := os.MkdirTemp("", "wurfl_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Download() downloads the WURFL data file from the specified URL and saves it to the specified directory.
	// If the download is successful, the function returns nil. Otherwise, it returns an error.
	err = wurfl.Download(url, tempDir)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Check if the file was downloaded
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	if len(files) == 0 {
		t.Errorf("No files downloaded")
	}

	for _, file := range files {
		if file.Name() == "wurfl.zip" {
			wurflfile := tempDir + "/" + file.Name()
			// Try to create and load an engine with the freshly downloaded file
			_, err = wurfl.Create(wurflfile, nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
			if err != nil {
				t.Errorf("Create error : %s\n", err.Error())
			}
			return
		}
	}

	t.Errorf("wurfl.zip not found in downloaded files")
}
func TestGetStaticCap(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

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

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

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

func TestGetVirtualCap(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

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

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

	concurrency := 100
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			_, err := device.GetVirtualCap("is_smartphone")
			if err != nil {
				t.Errorf("Concurrent GetVirtualCap() failed: %v", err)
			}
		}()
	}

	wg.Wait()
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

func TestGetStaticCaps(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ValidCapabilities", func(t *testing.T) {
		caps := []string{"mobile_browser", "is_tablet", "pointing_method"}
		result, err := device.GetStaticCaps(caps)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
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
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d capabilities", len(result))
		}
	})

	t.Run("NonExistentCapability", func(t *testing.T) {
		caps := []string{"non_existent_cap"}
		result, err := device.GetStaticCaps(caps)
		if err == nil {
			t.Error("Expected error for non-existent capability, got nil")
		}
		if result == nil {
			t.Errorf("Expected empty map for non-existent capability, got nil")
		}
	})

	t.Run("MixedCapabilities", func(t *testing.T) {
		caps := []string{"mobile_browser", "non_existent_cap", "is_tablet"}
		result, err := device.GetStaticCaps(caps)
		if err == nil {
			t.Error("Expected error for mixed capabilities, got nil")
		}
		if result == nil {
			t.Errorf("Expected valid map for mixed capabilities, got nil")
		}
		// We expect wrong capability to be skipped
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})

	t.Run("DuplicateCaps", func(t *testing.T) {
		caps := []string{"mobile_browser", "mobile_browser", "is_tablet"}
		result, err := device.GetStaticCaps(caps)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// We expect one of the duplicate capabilities to be removed
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})
}

func TestGetVirtualCaps(t *testing.T) {
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 12_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0 Mobile/15E148 Safari/604.1"

	wengine := fixtureCreateEngine(t)

	defer wengine.Destroy()

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ValidCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "advertised_device_os", "pixel_density"}
		result, err := device.GetVirtualCaps(caps)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
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
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("Expected empty result, got %d capabilities", len(result))
		}
	})

	t.Run("NonExistentCap", func(t *testing.T) {
		caps := []string{"non_existent_cap"}
		result, err := device.GetVirtualCaps(caps)
		if err == nil {
			t.Error("Expected error for non-existent capability, got nil")
		}
		if result == nil {
			t.Errorf("Expected empty map for non-existent capability, got nil")
		}
	})

	t.Run("MixedCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "non_existent_cap", "pixel_density"}
		result, err := device.GetVirtualCaps(caps)
		if err == nil {
			t.Error("Expected error for mixed valid and non-existent capabilities, got nil")
		}
		if result == nil {
			t.Errorf("Expected empty map for non-existent capability, got nil")
		}
		// We expect wrong virtual capability to be skipped
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})

	t.Run("DuplicateCaps", func(t *testing.T) {
		caps := []string{"is_smartphone", "advertised_device_os", "is_smartphone"}
		result, err := device.GetVirtualCaps(caps)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// We expect one of the duplicate capabilities to be removed
		if len(result) != len(caps)-1 {
			t.Errorf("Expected %d capabilities, got %d", len(caps)-1, len(result))
		}
	})
}

func TestGetLastUpdated(t *testing.T) {
	// get env var SM_UPDATER_DATA_URL url value (from scientiamobile Vault)
	url := os.Getenv("SM_UPDATER_DATA_URL")
	if url == "" {
		t.Skip("SM_UPDATER_DATA_URL environment var not set")
	}

	// Download() downloads the WURFL data file from the specified URL and saves it to the specified directory.
	// If the download is successful, the function returns nil. Otherwise, it returns an error.
	err := wurfl.Download(url, ".")
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	wengine, err := wurfl.Create("./wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")
	if err != nil {
		t.Fatal(err)
	}
	defer wengine.Destroy()

	lastUpdated := wengine.GetLastUpdated()
	if lastUpdated == "" {
		t.Errorf("Got empty string as wurfl last updated time")
	}

	t.Logf("Last updated: %s", lastUpdated)

	// Remove the local wurfl.zip file
	err = os.Remove("./wurfl.zip")
	if err != nil {
		t.Fatalf("Failed to remove wurfl.zip: %v", err)
	}
}
