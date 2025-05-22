// Package wurfl - WURFL InFuze Golang module
// wurfl is a golang package wrapping the WURFL C API and encapsulating it in
// 2 golang types to provide a fast and intuitive interface.
// It is released for linux/macos platforms.
package wurfl

//
//#cgo darwin CFLAGS: -I/usr/local/include
//#cgo darwin LDFLAGS: -L/usr/local/lib/
//#cgo windows CFLAGS: -I"C:/Program Files/Scientiamobile/InFuze/dev/include"
//#cgo windows LDFLAGS: -L"C:/Program Files/Scientiamobile/InFuze/bin"
//#cgo LDFLAGS: -lwurfl
//#include <stdlib.h>
//#include <wurfl/wurfl.h>
//#include <stdio.h>
import "C"

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unsafe"
)

// Engine Target possible values
// DEPRECATED as of 1.9.5.0
const (
	WurflEngineTargetHighAccuray             = C.WURFL_ENGINE_TARGET_HIGH_ACCURACY
	WurflEngineTargetHighPerformance         = C.WURFL_ENGINE_TARGET_HIGH_PERFORMANCE
	WurflEngineTargetDefault                 = C.WURFL_ENGINE_TARGET_DEFAULT
	WurflEngineTargetFastDesktopBrowserMatch = C.WURFL_ENGINE_TARGET_FAST_DESKTOP_BROWSER_MATCH
)

// UserAgent priority possible values
// DEPRECATED as of 1.9.5.0
const (
	WurflUserAgentPriorityOverrideSideloadedBrowserUserAgent = C.WURFL_USERAGENT_PRIORITY_OVERRIDE_SIDELOADED_BROWSER_USERAGENT
	WurflUserAgentPriorityUsePlainUserAgent                  = C.WURFL_USERAGENT_PRIORITY_USE_PLAIN_USERAGENT
)

const (
	// WurflAttrExtraHeadersExperimental is deprecated since 1.12.5.0 and should not be used.
	WurflAttrExtraHeadersExperimental = C.WURFL_ATTR_EXTRA_HEADERS_EXPERIMENTAL
	// WurflAttrCapabilityFallbackCache attribute controls the capability fallback cache (needs libwurfl version 1.12.9.3 or above)
	WurflAttrCapabilityFallbackCache = C.WURFL_ATTR_CAPABILITY_FALLBACK_CACHE
)

// To control capability fallback cache, needs libwurfl version 1.12.9.3 or above
const (
	// WurflAttrCapabilityFallbackCacheDefault is the default setting for capability fallback cache
	WurflAttrCapabilityFallbackCacheDefault = C.WURFL_ATTR_CAPABILITY_FALLBACK_CACHE_DEFAULT
	// WurflAttrCapabilityFallbackCacheDisabled disables the capability fallback cache
	WurflAttrCapabilityFallbackCacheDisabled = C.WURFL_ATTR_CAPABILITY_FALLBACK_CACHE_DISABLED
	// WurflAttrCapabilityFallbackCacheLimited sets the capability fallback cache to limited mode
	WurflAttrCapabilityFallbackCacheLimited = C.WURFL_ATTR_CAPABILITY_FALLBACK_CACHE_LIMITED
)

// Cache Provider possible values
const (
	WurflCacheProviderDefault = -1
	WurflCacheProviderNone    = C.WURFL_CACHE_PROVIDER_NONE
	WurflCacheProviderLru     = C.WURFL_CACHE_PROVIDER_LRU
	// Deprecated: use WurflCacheProviderLru instead
	WurflCacheProviderDoubleLru = C.WURFL_CACHE_PROVIDER_DOUBLE_LRU
)

// Match type
const (
	WurflMatchTypeExact           = C.WURFL_MATCH_TYPE_EXACT
	WurflMatchTypeConclusive      = C.WURFL_MATCH_TYPE_CONCLUSIVE
	WurflMatchTypeRecovery        = C.WURFL_MATCH_TYPE_RECOVERY
	WurflMatchTypeCatchall        = C.WURFL_MATCH_TYPE_CATCHALL
	WurflMatchTypeHighPerformance = C.WURFL_MATCH_TYPE_HIGHPERFORMANCE
	WurflMatchTypeNone            = C.WURFL_MATCH_TYPE_NONE
	WurflMatchTypeCached          = C.WURFL_MATCH_TYPE_CACHED
)

// Wurfl enumerator type
const (
	WurflEnumStaticCapabilities    = C.WURFL_ENUM_STATIC_CAPABILITIES
	WurflEnumVirtualCapabilities   = C.WURFL_ENUM_VIRTUAL_CAPABILITIES
	WurflEnumMandatoryCapabilities = C.WURFL_ENUM_MANDATORY_CAPABILITIES
	WurflEnumWurflID               = C.WURFL_ENUM_WURFLID
)

// Wurfl updater frequency
const (
	WurflUpdaterFrequencyDaily  = C.WURFL_UPDATER_FREQ_DAILY
	WurflUpdaterFrequencyWeekly = C.WURFL_UPDATER_FREQ_WEEKLY
)

// HeaderQuality represents the header quality value
type HeaderQuality int

func (hq HeaderQuality) String() string {
	switch hq {
	case HeaderQualityNone:
		return "None"
	case HeaderQualityBasic:
		return "Basic"
	case HeaderQualityFull:
		return "Full"
	}
	return "Unknown"
}

const (
	// HeaderQualityNone no User Agent Client Hints are present.
	HeaderQualityNone HeaderQuality = C.WURFL_ENUM_UACH_NONE
	// HeaderQualityBasic only some of the headers needed for a successful WURFL detection are present.
	HeaderQualityBasic HeaderQuality = C.WURFL_ENUM_UACH_BASIC
	// HeaderQualityFull all the headers needed for a successful WURFL detection are present.
	HeaderQualityFull HeaderQuality = C.WURFL_ENUM_UACH_FULL
)

// Wurfl represents internal wurfl infuze handle
type Wurfl struct {
	Wurfl                       C.wurfl_handle
	ImportantHeaderNames        []string
	importantHeaderCStringNames []*C.char
	capsCStringcache            map[string]*C.char
}

// Device represent internal matched device handle
type Device struct {
	Device           C.wurfl_device_handle
	Wurfl            C.wurfl_handle
	capsCStringcache map[string]*C.char
}

// WurflHandler defines API methods for the Wurfl Infuze handle
type WurflHandler interface {
	GetAPIVersion() string
	SetLogPath(LogFile string) error
	IsUserAgentFrozen(ua string) bool
	LookupDeviceIDWithImportantHeaderMap(DeviceID string, IHMap map[string]string) (DeviceHandler, error)
	LookupWithImportantHeaderMap(IHMap map[string]string) (DeviceHandler, error)
	LookupDeviceIDWithRequest(DeviceID string, r *http.Request) (DeviceHandler, error)
	LookupRequest(r *http.Request) (DeviceHandler, error)
	LookupUserAgent(ua string) (DeviceHandler, error)
	GetAllDeviceIds() []string
	LookupDeviceID(DeviceID string) (DeviceHandler, error)
	Destroy()
	GetAllVCaps() []string
	GetAllCaps() []string
	GetInfo() string
	GetHeaderQuality(r *http.Request) (HeaderQuality, error)
	GetLastLoadTime() string
	HasCapability(cap string) bool
	HasVirtualCapability(vcap string) bool
	SetAttr(attr int, value int) error
	GetAttr(attr int) (int, error)
	GetLastUpdated() string
}

// DeviceHandler defines API methods for the Wurfl Device handle
type DeviceHandler interface {
	GetMatchType() int
	GetVirtualCapabilities(caps []string) map[string]string
	GetVirtualCaps(caps []string) (map[string]string, error)
	GetVirtualCapability(vcap string) string
	GetVirtualCap(vcap string) (string, error)
	GetVirtualCapabilityAsInt(vcsp string) (int, error)
	GetCapabilities(caps []string) map[string]string
	GetStaticCaps(caps []string) (map[string]string, error)
	GetCapability(cap string) string
	GetStaticCap(cap string) (string, error)
	GetCapabilityAsInt(cap string) (int, error)
	IsRoot() bool
	GetRootID() string
	GetParentID() string
	GetDeviceID() (string, error)
	GetNormalizedUserAgent() (string, error)
	GetOriginalUserAgent() (string, error)
	GetUserAgent() (string, error)
	Destroy()
}

// Updater defines API methods for the Updater
type Updater interface {
	SetUpdaterDataURL(DataURL string) error
	SetUpdaterDataFrequency(Frequency int) error
	SetUpdaterDataURLTimeout(ConnectionTimeout int, DataTransferTimeout int) error
	SetUpdaterLogPath(LogFile string) error
	UpdaterRunonce() error
	UpdaterStart() error
	UpdaterStop() error
	SetUpdaterUserAgent(userAgent string) error
	GetUpdaterUserAgent() string
}

// Version is the current version of this package.
const Version = "1.31.0"

// APIVersion returns version of internal InFuze API without an initialized engine
func APIVersion() string {
	return C.GoString(C.wurfl_get_api_version())
}

// Download downloads the WURFL data file from the specified URL and saves it to the specified folder.
// If the download is successful, it returns nil. Otherwise, it returns an error.
func Download(url string, folder string) error {
	cURL := C.CString(url)
	cFolder := C.CString(folder)
	defer C.free(unsafe.Pointer(cURL))
	defer C.free(unsafe.Pointer(cFolder))
	cerr := C.wurfl_download(cURL, cFolder)
	if cerr != C.WURFL_OK {
		errMsg := C.GoString(C.wurfl_get_error_string(cerr))
		return fmt.Errorf("WurflDownload failed: %s", errMsg)
	}
	return nil
}

// Create the wurfl engine. Parameters :
// Wurflxml : path to the wurfl.xml/zip file
// Patches : slice of paths of patches files to load
// CapFilter : list of capabilities used; allow to init engine without loading all 500+ caps
//
//	Note : Capability filtering is discouraged and will be deprecated in future versions
//
// DEPRECATED: EngineTarget : As of 1.9.5.0 has no effect anymore
// CacheProvider : WurflCacheProviderLru
// CacheExtraConfig : size of single lru cache in the form "100000"
func Create(Wurflxml string, Patches []string, CapFilter []string, EngineTarget int, CacheProvider int, CacheExtraConfig string) (*Wurfl, error) {
	w := &Wurfl{}

	w.Wurfl = C.wurfl_create()

	if w.Wurfl == nil {
		// error in create : no way to get the error as the is no engine instance yet
		// in libwurfl. We can only return a generic memory allocation error
		return nil, cErrorToGoError(C.WURFL_ERROR_UNABLE_TO_ALLOCATE_MEMORY)
	}

	// setting cache if specified
	if CacheProvider != WurflCacheProviderDefault {
		ccacheec := C.CString(CacheExtraConfig)

		cp := C.wurfl_cache_provider(CacheProvider)
		C.wurfl_set_cache_provider(w.Wurfl, cp, ccacheec)
		C.free(unsafe.Pointer(ccacheec))
	}

	// setting wurfl.xml
	wxml := C.CString(Wurflxml)
	if ret := C.wurfl_set_root(w.Wurfl, wxml); ret != C.WURFL_OK {
		C.free(unsafe.Pointer(wxml))
		w.Destroy()
		return nil, cErrorToGoError(ret)
	}

	// setting patches
	for i := 0; i < len(Patches); i++ {
		cpatch := C.CString(Patches[i])
		if ret := C.wurfl_add_patch(w.Wurfl, cpatch); ret != C.WURFL_OK {
			C.free(unsafe.Pointer(cpatch))
			w.Destroy()
			return nil, cErrorToGoError(ret)
		}
		C.free(unsafe.Pointer(cpatch))
	}

	// filter capabilities in engine
	for i := 0; i < len(CapFilter); i++ {
		ccap := C.CString(CapFilter[i])
		if ret := C.wurfl_add_requested_capability(w.Wurfl, ccap); ret != C.WURFL_OK {
			C.free(unsafe.Pointer(ccap))
			w.Destroy()
			return nil, cErrorToGoError(ret)
		}
		C.free(unsafe.Pointer(ccap))
	}

	// loading engine
	if C.wurfl_load(w.Wurfl) != C.WURFL_OK {
		// we prefer wurfl handle based error message as it is richer than the standard one
		err := checkHandleError(w.Wurfl)
		w.Destroy()
		return nil, err
	}

	// prepare important headers slice
	ihe := C.wurfl_get_important_header_enumerator(w.Wurfl)
	if ihe == nil { // Check if enumerator creation failed
		err := checkHandleError(w.Wurfl)
		w.Destroy()
		return nil, err
	}
	defer C.wurfl_important_header_enumerator_destroy(ihe)

	for i := 0; C.wurfl_important_header_enumerator_is_valid(ihe) != 0; i++ {
		// get the header name
		headerName := C.wurfl_important_header_enumerator_get_value(ihe)
		// convert header name to go string
		gheaderName := C.GoString(headerName)
		// create a C string copy from the go string
		cheaderName := C.CString(gheaderName) // This CString needs to be managed (freed in Destroy)
		// append to slice
		w.ImportantHeaderNames = append(w.ImportantHeaderNames, gheaderName)
		w.importantHeaderCStringNames = append(w.importantHeaderCStringNames, cheaderName)
		// advance
		C.wurfl_important_header_enumerator_move_next(ihe)
	}

	// initialize caps/vcaps CString cache for faster calls to libwurfl

	caps := w.GetAllCaps()
	vcaps := w.GetAllCaps()

	w.capsCStringcache = make(map[string]*C.char, len(caps)+len(vcaps))

	for c := range caps {
		w.capsCStringcache[caps[c]] = C.CString(caps[c])
	}

	for v := range vcaps {
		w.capsCStringcache[vcaps[v]] = C.CString(vcaps[v])
	}

	return w, nil
}

// Destroy the wurfl engine
func (w *Wurfl) Destroy() {
	if w.Wurfl != nil {

		// deallocate important headers C strings
		for _, importantHeaderName := range w.importantHeaderCStringNames {
			if importantHeaderName != nil {
				C.free(unsafe.Pointer(importantHeaderName))
			}
		}

		// now free the caps/vcaps CStrings cache

		for _, v := range w.capsCStringcache {
			if v != nil {
				C.free(unsafe.Pointer(v))
			}
		}
		w.capsCStringcache = nil // Clear the map
		C.wurfl_destroy(w.Wurfl)
		w.Wurfl = nil
	}
}

// SetAttr : set engine attributes
func (w *Wurfl) SetAttr(attr int, value int) error {
	cattr := C.wurfl_attr(attr)
	cvalue := C.int(value)
	if C.wurfl_set_attr(w.Wurfl, cattr, cvalue) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}

	// now reload all important header in wurfl struct since they are different
	// if the attr is WurflAttrExtraHeadersExperimental
	if attr == WurflAttrExtraHeadersExperimental {
		ihe := C.wurfl_get_important_header_enumerator(w.Wurfl)

		w.ImportantHeaderNames = nil
		// deallocate important headers C strings
		for _, importantHeaderName := range w.importantHeaderCStringNames {
			C.free(unsafe.Pointer(importantHeaderName))
		}
		w.importantHeaderCStringNames = nil

		for i := 0; C.wurfl_important_header_enumerator_is_valid(ihe) != 0; i++ {
			// get the header name
			headerName := C.wurfl_important_header_enumerator_get_value(ihe)
			// convert header name to go string
			gheaderName := C.GoString(headerName)
			// create a C string copy from the go string
			cheaderName := C.CString(gheaderName)
			// append to slice
			w.ImportantHeaderNames = append(w.ImportantHeaderNames, gheaderName)
			w.importantHeaderCStringNames = append(w.importantHeaderCStringNames, cheaderName)
			// advance
			C.wurfl_important_header_enumerator_move_next(ihe)
		}
		C.wurfl_important_header_enumerator_destroy(ihe)
	}

	return nil
}

// GetAttr : get engine attributes
func (w *Wurfl) GetAttr(attr int) (int, error) {
	cattr := C.wurfl_attr(attr)
	var cvalue C.int
	if C.wurfl_get_attr(w.Wurfl, cattr, &cvalue) != C.WURFL_OK {
		return 0, checkHandleError(w.Wurfl)
	}
	return int(cvalue), nil
}

// SetLogPath - set path of main libwurfl log file (updater has a separate log file)
func (w *Wurfl) SetLogPath(LogFile string) error {
	// wurfl_error wurfl_set_log_path(wurfl_handle hwurfl, const char *log_path);
	clog := C.CString(LogFile)
	ret := C.wurfl_set_log_path(w.Wurfl, clog)
	C.free(unsafe.Pointer(clog))
	if ret != C.WURFL_OK {
		return cErrorToGoError(ret)
	}
	return nil
}

// SetUpdaterDataURL - set your scientiamobile WURFL Snapshot URL
func (w *Wurfl) SetUpdaterDataURL(DataURL string) error {

	apiVersion := w.GetAPIVersion()
	// we set useragent only if API version is >= 1.13.0.0 otherwise it will overwrite the libwurfl one
	if CompareVersions(apiVersion, "1.13.0.0") >= 0 {
		golangUA := "infuze_golang/" + Version
		cgolangUA := C.CString(golangUA)
		cret := C.wurfl_updater_set_useragent(w.Wurfl, cgolangUA)
		C.free(unsafe.Pointer(cgolangUA))
		if cret != C.WURFL_OK {
			return cErrorToGoError(cret)
		}
	}

	cdata := C.CString(DataURL)

	ret := C.wurfl_updater_set_data_url(w.Wurfl, cdata)
	C.free(unsafe.Pointer(cdata))

	if ret != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// SetUpdaterUserAgent - set the UserAgent used in calling the WURFL Snapshot server
func (w *Wurfl) SetUpdaterUserAgent(userAgent string) error {
	cdata := C.CString(userAgent)
	ret := C.wurfl_updater_set_useragent(w.Wurfl, cdata)
	C.free(unsafe.Pointer(cdata))
	if ret != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// GetUpdaterUserAgent - gets the UserAgent used in calling the WURFL Snapshot server
func (w *Wurfl) GetUpdaterUserAgent() string {
	ua := C.wurfl_updater_get_useragent(w.Wurfl)
	uaValue := C.GoString(ua)
	return uaValue
}

// SetUpdaterDataFrequency - set frequency of update checks
func (w *Wurfl) SetUpdaterDataFrequency(Frequency int) error {
	//     LIBWURFLAPI wurfl_error wurfl_updater_set_data_frequency(wurfl_handle hwurfl, wurfl_updater_frequency freq);
	cfreq := C.wurfl_updater_frequency(Frequency)
	if C.wurfl_updater_set_data_frequency(w.Wurfl, cfreq) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// SetUpdaterDataURLTimeout - set connection and data transfer timeouts (in millisecs) for updater
// http call. 0 for no timeout, -1 for defaults
func (w *Wurfl) SetUpdaterDataURLTimeout(ConnectionTimeout int, DataTransferTimeout int) error {
	// wurfl_error wurfl_updater_set_data_url_timeouts(wurfl_handle hwurfl, int connection_timeout, int data_transfer_timeout);
	cConn := C.int(ConnectionTimeout)
	cData := C.int(DataTransferTimeout)
	if C.wurfl_updater_set_data_url_timeouts(w.Wurfl, cConn, cData) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// SetUpdaterLogPath - set path of updater log file
func (w *Wurfl) SetUpdaterLogPath(LogFile string) error {
	//     LIBWURFLAPI wurfl_error wurfl_updater_set_log_path(wurfl_handle hwurfl, const char* log_path);
	clog := C.CString(LogFile)
	ret := C.wurfl_updater_set_log_path(w.Wurfl, clog)
	C.free(unsafe.Pointer(clog))
	if ret != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// UpdaterRunonce - Update the wurfl if needed and terminate
func (w *Wurfl) UpdaterRunonce() error {
	//     LIBWURFLAPI wurfl_error wurfl_updater_runonce(wurfl_handle hwurfl);
	if C.wurfl_updater_runonce(w.Wurfl) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// UpdaterStart - Start the updater, a thread that performs periodic check and update of the wurfl.zip file
// when a new wurfl.zip is available it is downloaded and engine is switched to use the new wurfl.zip file immediately
func (w *Wurfl) UpdaterStart() error {
	//     LIBWURFLAPI wurfl_error wurfl_updater_start(wurfl_handle hwurfl);
	if C.wurfl_updater_start(w.Wurfl) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// UpdaterStop - stop the updater
func (w *Wurfl) UpdaterStop() error {
	if C.wurfl_updater_stop(w.Wurfl) != C.WURFL_OK {
		return checkHandleError(w.Wurfl)
	}
	return nil
}

// GetAPIVersion returns version of internal InFuze API
func (w *Wurfl) GetAPIVersion() string {
	return C.GoString(C.wurfl_get_api_version())
}

// GetAllVCaps return all virtual capabilities names
func (w *Wurfl) GetAllVCaps() []string {
	var result []string

	eh := C.wurfl_enum_create(w.Wurfl, WurflEnumVirtualCapabilities)
	defer C.wurfl_enum_destroy(eh)

	for C.wurfl_enum_is_valid(eh) != 0 {
		cvcapname := C.wurfl_enum_get_name(eh)
		vcapname := C.GoString(cvcapname)
		result = append(result, vcapname)
		C.wurfl_enum_move_next(eh)
	}

	return result
}

// GetAllCaps return all static capabilities names
func (w *Wurfl) GetAllCaps() []string {
	var result []string

	eh := C.wurfl_enum_create(w.Wurfl, WurflEnumStaticCapabilities)
	defer C.wurfl_enum_destroy(eh)

	for C.wurfl_enum_is_valid(eh) != 0 {
		ccapname := C.wurfl_enum_get_name(eh)
		capname := C.GoString(ccapname)
		result = append(result, capname)
		C.wurfl_enum_move_next(eh)
	}

	return result
}

// GetInfo - get wurfl.xml info
func (w *Wurfl) GetInfo() string {
	return C.GoString(C.wurfl_get_wurfl_info(w.Wurfl))
}

// GetLastLoadTime - get last wurfl.xml load time
func (w *Wurfl) GetLastLoadTime() string {
	return C.GoString(C.wurfl_get_last_load_time_as_string(w.Wurfl))
}

// GetLastUpdated - get last wurfl.xml update time
func (w *Wurfl) GetLastUpdated() string {
	return C.GoString(C.wurfl_get_last_updated(w.Wurfl))
}

// GetEngineTarget - Returns a string representing the currently set WURFL Engine Target.
// DEPRECATED: will always return default value
func (w *Wurfl) GetEngineTarget() string {
	return "DEFAULT"
}

// SetUserAgentPriority - Sets which UA wurfl is using
// DEPRECATED. Since 1.9.5.0 has no effect anymore
func (w *Wurfl) SetUserAgentPriority(prio int) {
	return
}

// GetUserAgentPriority - Tells if WURFL is using the plain user agent or the sideloaded browser user agent for device detection
// DEPRECATED: will always return default value
func (w *Wurfl) GetUserAgentPriority() string {
	return "OVERRIDE SIDELOADED BROWSER USERAGENT"
}

// HasCapability - returns true if the static capability exists in wurfl.zip
func (w *Wurfl) HasCapability(cap string) bool {
	ccap := C.CString(cap)
	ret := C.wurfl_has_capability(w.Wurfl, ccap)
	C.free(unsafe.Pointer(ccap))
	if ret == 0 {
		return false
	}
	return true
}

// HasVirtualCapability - returns true if the virtual cap is available
func (w *Wurfl) HasVirtualCapability(vcap string) bool {
	cvcap := C.CString(vcap)
	ret := C.wurfl_has_virtual_capability(w.Wurfl, cvcap)
	C.free(unsafe.Pointer(cvcap))
	if ret == 0 {
		return false
	}
	return true
}

// LookupDeviceID : lookup by wurfl_ID and return Device handle
func (w *Wurfl) LookupDeviceID(DeviceID string) (*Device, error) {
	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	wDeviceID := C.CString(DeviceID)

	d.Device = C.wurfl_get_device(w.Wurfl, wDeviceID)
	C.free(unsafe.Pointer(wDeviceID))
	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// LookupUserAgent : lookup up useragent and return Device handle
func (w *Wurfl) LookupUserAgent(ua string) (*Device, error) {
	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	wua := C.CString(ua)

	d.Device = C.wurfl_lookup_useragent(w.Wurfl, wua)
	C.free(unsafe.Pointer(wua))
	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// LookupRequest : Lookup using Request headers and return Device handle
func (w *Wurfl) LookupRequest(r *http.Request) (*Device, error) {

	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	// create important headers object to pass to lookup

	cih := C.wurfl_important_header_create(w.Wurfl)
	if cih == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	defer C.wurfl_important_header_destroy(cih)

	// use important header names loaded during create
	for i, importantHeaderName := range w.ImportantHeaderNames {
		// retrieve header value from passed request, if any
		headerValue := r.Header.Get(importantHeaderName)
		if len(headerValue) != 0 {
			// create C strings from header value
			cheaderValue := C.CString(headerValue)

			// add this header to cih
			// for header names we use a set of preallocated CStrings with headernames
			C.wurfl_important_header_set(cih, w.importantHeaderCStringNames[i], cheaderValue)
			C.free(unsafe.Pointer(cheaderValue))
		}
	}

	d.Device = C.wurfl_lookup_with_important_header(w.Wurfl, cih)

	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// LookupDeviceIDWithRequest : lookup by wurfl_ID and request headers and return Device handle
func (w *Wurfl) LookupDeviceIDWithRequest(DeviceID string, r *http.Request) (*Device, error) {
	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	wDeviceID := C.CString(DeviceID)

	// create important headers object to pass to lookup

	cih := C.wurfl_important_header_create(w.Wurfl)
	if cih == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	defer C.wurfl_important_header_destroy(cih)

	// use important header names loaded during create
	for i, importantHeaderName := range w.ImportantHeaderNames {
		// retrieve header value from passed request, if any
		headerValue := r.Header.Get(importantHeaderName)
		if len(headerValue) != 0 {
			// create C strings from header name and value
			cheaderValue := C.CString(headerValue)

			// add this header to cih
			C.wurfl_important_header_set(cih, w.importantHeaderCStringNames[i], cheaderValue)
			C.free(unsafe.Pointer(cheaderValue))
		}
	}

	d.Device = C.wurfl_get_device_with_important_header(w.Wurfl, wDeviceID, cih)
	C.free(unsafe.Pointer(wDeviceID))
	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// LookupWithImportantHeaderMap : Lookup using header values found in IHMap.
// IHMap must be filled with Wurfl.ImportantHeaderNames and values
func (w *Wurfl) LookupWithImportantHeaderMap(IHMap map[string]string) (*Device, error) {
	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	// create important headers object to pass to lookup

	cih := C.wurfl_important_header_create(w.Wurfl)
	if cih == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	defer C.wurfl_important_header_destroy(cih)
	// fill it with IHMap entries
	for importantHeaderName, headerValue := range IHMap {
		// create C strings from header name and value
		cheaderName := C.CString(importantHeaderName)
		cheaderValue := C.CString(headerValue)

		// add this header to WURFL importtant headers object
		C.wurfl_important_header_set(cih, cheaderName, cheaderValue)
		C.free(unsafe.Pointer(cheaderName))
		C.free(unsafe.Pointer(cheaderValue))
	}

	d.Device = C.wurfl_lookup_with_important_header(w.Wurfl, cih)

	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// LookupDeviceIDWithImportantHeaderMap : Lookup deviceID using header values found in IHMap.
// IHMap must be filled with Wurfl.ImportantHeaderNames and values
func (w *Wurfl) LookupDeviceIDWithImportantHeaderMap(DeviceID string, IHMap map[string]string) (*Device, error) {
	d := &Device{}
	// copy wurfl handle into device handle for error handling
	d.Wurfl = w.Wurfl
	// copy the caps cache
	d.capsCStringcache = w.capsCStringcache

	cDeviceID := C.CString(DeviceID)

	// create important headers object to pass to lookup

	cih := C.wurfl_important_header_create(w.Wurfl)
	if cih == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	defer C.wurfl_important_header_destroy(cih)

	// fill it with IHMap entries

	for i, importantHeaderName := range w.ImportantHeaderNames {
		// retrieve header from IHMap
		header := IHMap[importantHeaderName]
		if len(header) != 0 {

			// create a C string from header
			cheader := C.CString(header)

			// add this header to cih
			C.wurfl_important_header_set(cih, w.importantHeaderCStringNames[i], cheader)
			C.free(unsafe.Pointer(cheader))
		}
	}

	d.Device = C.wurfl_get_device_with_important_header(w.Wurfl, cDeviceID, cih)
	C.free(unsafe.Pointer(cDeviceID))
	if d.Device == nil {
		return nil, checkHandleError(w.Wurfl)
	}
	return d, nil
}

// IsUserAgentFrozen : returns true if a UserAgent is frozen
func (w *Wurfl) IsUserAgentFrozen(ua string) bool {
	wua := C.CString(ua)
	ret := C.wurfl_is_ua_frozen(w.Wurfl, wua)
	C.free(unsafe.Pointer(wua))
	if ret == 0 {
		return false
	}
	return true
}

// GetHeaderQuality returns an indicator of how many sec-ch-ua headers are present in the request
func (w *Wurfl) GetHeaderQuality(r *http.Request) (HeaderQuality, error) {
	// create important headers object to pass to lookup
	cih := C.wurfl_important_header_create(w.Wurfl)
	if cih == nil {
		return HeaderQualityNone, checkHandleError(w.Wurfl)

	}
	defer C.wurfl_important_header_destroy(cih)

	// use important header names loaded during create
	for i, importantHeaderName := range w.ImportantHeaderNames {
		// retrieve header value from passed request, if any
		headerValue := r.Header.Get(importantHeaderName)
		if len(headerValue) != 0 {
			// create C strings from header name and value
			cheaderValue := C.CString(headerValue)

			// add this header to cih
			C.wurfl_important_header_set(cih, w.importantHeaderCStringNames[i], cheaderValue)
			C.free(unsafe.Pointer(cheaderValue))
		}
	}

	hq := C.wurfl_important_header_uach_quality(cih)
	return HeaderQuality(hq), nil
}

/*
 * Device methods
 */

// GetUserAgent Get default UserAgent of matched device (might be different from UA passed to lookup)
func (d *Device) GetUserAgent() (string, error) {
	cua := C.wurfl_device_get_useragent(d.Device)
	if cua == nil {
		return "", checkHandleError(d.Wurfl)
	}
	ua := C.GoString(cua)
	return ua, nil
}

// GetOriginalUserAgent Get the original userAgent of matched device (the one passed to lookup)
func (d *Device) GetOriginalUserAgent() (string, error) {
	oua := C.wurfl_device_get_original_useragent(d.Device)
	if oua == nil {
		return "", checkHandleError(d.Wurfl)
	}
	ua := C.GoString(oua)
	return ua, nil
}

// GetNormalizedUserAgent Get the Normalized (processed by wurfl api) userAgent ( Only for internal use/tooling)
func (d *Device) GetNormalizedUserAgent() (string, error) {
	nua := C.wurfl_device_get_normalized_useragent(d.Device)
	if nua == nil {
		return "", checkHandleError(d.Wurfl)
	}
	ua := C.GoString(nua)
	return ua, nil
}

// GetDeviceID Get wurfl_id string from device handle
func (d *Device) GetDeviceID() (string, error) {
	cdeviceid := C.wurfl_device_get_id(d.Device)
	if cdeviceid == nil {
		return "", checkHandleError(d.Wurfl)
	}
	deviceid := C.GoString(cdeviceid)
	return deviceid, nil
}

// GetRootID - Retrieve the root device id of this device.
func (d *Device) GetRootID() string {
	return C.GoString(C.wurfl_device_get_root_id(d.Device))
}

// GetParentID - Retrieve the parent device id of this device.
func (d *Device) GetParentID() string {
	return C.GoString(C.wurfl_device_get_parent_id(d.Device))
}

// IsRoot - true if device is device root
func (d *Device) IsRoot() bool {
	if C.wurfl_device_is_actual_device_root(d.Device) == 0 {
		return false
	}
	return true
}

// GetCapability Get a single Capability
// Deprecated: GetCapability is deprecated. Use GetStaticCap instead.
func (d *Device) GetCapability(cap string) string {
	ccap, found := d.capsCStringcache[cap]
	if !found {
		// non existing capability?
		ccap = C.CString(cap)
		defer C.free(unsafe.Pointer(ccap))
	}

	ccapvalue := C.wurfl_device_get_capability(d.Device, ccap)
	capvalue := C.GoString(ccapvalue)

	return capvalue
}

// GetStaticCap Get a single static cap using new C.wurfl_device_get_static_cap()
// that returns error
func (d *Device) GetStaticCap(cap string) (string, error) {
	ccap, found := d.capsCStringcache[cap]
	if !found {
		// non existing capability?
		ccap = C.CString(cap)
		defer C.free(unsafe.Pointer(ccap))
	}
	retCode := C.wurfl_error(0)
	ccapvalue := C.wurfl_device_get_static_cap(d.Device, ccap, &retCode)
	if retCode != C.WURFL_OK {
		return "", checkHandleError(d.Wurfl)
	}
	capvalue := C.GoString(ccapvalue)

	return capvalue, nil
}

// GetCapabilityAsInt gets a single static capability value that has a int type
// It returns an error if the requested static capability is not a numeric one (ie: brand_name)
func (d *Device) GetCapabilityAsInt(cap string) (int, error) {
	ccap, found := d.capsCStringcache[cap]
	if !found {
		// non existing capability?
		ccap = C.CString(cap)
		C.free(unsafe.Pointer(ccap))
	}
	cErr := C.wurfl_error(0)
	ccapvalue := C.wurfl_device_get_static_cap_as_int(d.Device, ccap, &cErr)
	// libwurfl currently returns zero if any error occurs
	if cErr != C.WURFL_OK {
		return 0, checkHandleError(d.Wurfl)
	}

	return int(ccapvalue), nil

}

// GetCapabilities Get a list of Static Capabilities
// Deprecated: GetCapabilities is deprecated. Use GetStaticCaps instead.
func (d *Device) GetCapabilities(caps []string) map[string]string {
	result := make(map[string]string, len(caps))

	for i := 0; i < len(caps); i++ {
		ccap, found := d.capsCStringcache[caps[i]]
		if !found {
			// non existing capability?
			ccap = C.CString(caps[i])
			defer C.free(unsafe.Pointer(ccap))
		}

		ccapvalue := C.wurfl_device_get_capability(d.Device, ccap)
		capvalue := C.GoString(ccapvalue)
		result[caps[i]] = capvalue
	}

	return result
}

// GetStaticCaps Get a list of Static Capabilities
func (d *Device) GetStaticCaps(caps []string) (map[string]string, error) {
	var errMsg *C.char
	result := make(map[string]string, len(caps))

	for i := 0; i < len(caps); i++ {
		ccap, found := d.capsCStringcache[caps[i]]
		if !found {
			// non existing capability?
			ccap = C.CString(caps[i])
			defer C.free(unsafe.Pointer(ccap))
		}

		retCode := C.wurfl_error(0)
		ccapvalue := C.wurfl_device_get_static_cap(d.Device, ccap, &retCode)
		if retCode != C.WURFL_OK {
			// Just save error message for now, and continue with next capability
			errMsg = C.wurfl_get_error_message(d.Wurfl)
			continue
		}
		capvalue := C.GoString(ccapvalue)
		result[caps[i]] = capvalue
	}

	if errMsg != nil {
		return result, checkHandleError(d.Wurfl)
	}

	return result, nil
}

// GetVirtualCapability Get Virtual Capability
// Deprecated: GetVirtualCapability is deprecated. Use GetVirtualCap instead.
func (d *Device) GetVirtualCapability(vcap string) string {
	cvcap, found := d.capsCStringcache[vcap]
	if !found {
		// non existing capability?
		cvcap = C.CString(vcap)
		defer C.free(unsafe.Pointer(cvcap))
	}

	cvcapvalue := C.wurfl_device_get_virtual_capability(d.Device, cvcap)
	vcapvalue := C.GoString(cvcapvalue)

	return vcapvalue
}

// GetVirtualCap Get Virtual Cap with new C.wurfl_device_get_virtual_cap()
// that manages errors
func (d *Device) GetVirtualCap(vcap string) (string, error) {
	cvcap, found := d.capsCStringcache[vcap]
	if !found {
		// non existing capability?
		cvcap = C.CString(vcap)
		defer C.free(unsafe.Pointer(cvcap))
	}
	retCode := C.wurfl_error(0)
	cvcapvalue := C.wurfl_device_get_virtual_cap(d.Device, cvcap, &retCode)
	if retCode != C.WURFL_OK {
		return "", checkHandleError(d.Wurfl)
	}
	vcapvalue := C.GoString(cvcapvalue)

	return vcapvalue, nil
}

// GetVirtualCapabilityAsInt gets a single virtual capability value that has a int type
// It returns an error if the requested virtual capability is not a numeric one (ie: brand_name)
func (d *Device) GetVirtualCapabilityAsInt(vcap string) (int, error) {
	// the "C" vcap name
	cvcap, found := d.capsCStringcache[vcap]
	if !found {
		// non existing capability?
		cvcap = C.CString(vcap)
		defer C.free(unsafe.Pointer(cvcap))
	}
	cErr := C.wurfl_error(0)
	ccapvalue := C.wurfl_device_get_virtual_cap_as_int(d.Device, cvcap, &cErr)
	// libwurfl currently returns zero if any error occurs
	if cErr != C.WURFL_OK {
		return 0, checkHandleError(d.Wurfl)
	}

	return int(ccapvalue), nil

}

// GetVirtualCapabilities Get a list of Virtual Capabilities
// Deprecated: GetVirtualCapabilities is deprecated. Use GetVirtualCaps instead.
func (d *Device) GetVirtualCapabilities(caps []string) map[string]string {
	result := make(map[string]string)

	for i := 0; i < len(caps); i++ {
		ccap, found := d.capsCStringcache[caps[i]]
		if !found {
			// non existing capability?
			ccap = C.CString(caps[i])
			defer C.free(unsafe.Pointer(ccap))
		}

		ccapvalue := C.wurfl_device_get_virtual_capability(d.Device, ccap)
		capvalue := C.GoString(ccapvalue)
		result[caps[i]] = capvalue
	}

	return result
}

// GetVirtualCaps Get a list of Virtual Capabilities
func (d *Device) GetVirtualCaps(caps []string) (map[string]string, error) {
	var errMsg *C.char
	result := make(map[string]string, len(caps))

	for i := 0; i < len(caps); i++ {
		ccap, found := d.capsCStringcache[caps[i]]
		if !found {
			// non existing capability?
			ccap = C.CString(caps[i])
			defer C.free(unsafe.Pointer(ccap))
		}

		retCode := C.wurfl_error(0)
		ccapvalue := C.wurfl_device_get_virtual_cap(d.Device, ccap, &retCode)
		if retCode != C.WURFL_OK {
			// Just save error message for now, and continue with next capability
			errMsg = C.wurfl_get_error_message(d.Wurfl)
			continue
		}
		capvalue := C.GoString(ccapvalue)
		result[caps[i]] = capvalue
	}

	if errMsg != nil {
		return result, checkHandleError(d.Wurfl)
	}

	return result, nil
}

// GetMatchType Get type of Match occurred in lookup
func (d *Device) GetMatchType() int {

	cmtype := C.wurfl_device_get_match_type(d.Device)
	mtype := int(cmtype)

	return mtype
}

// Destroy device handle, should be called when when device attributes
// are not needed anymore
func (d *Device) Destroy() {
	if d.Device != nil {
		C.wurfl_device_destroy(d.Device)
		d.Device = nil
	}
}

// GetAllDeviceIds returns a slice containing all wurfl_id present in wurfl.zip
func (w *Wurfl) GetAllDeviceIds() []string {

	eh := C.wurfl_enum_create(w.Wurfl, WurflEnumWurflID)
	elen := C.wurfl_enum_len(eh)
	var result = make([]string, 0, elen)

	for C.wurfl_enum_is_valid(eh) != 0 {
		cdid := C.wurfl_enum_get_name(eh)
		did := C.GoString(cdid)
		if len(did) != 0 && did != "" {
			// add this id to the slice
			result = append(result, did)
		}
		C.wurfl_enum_move_next(eh)
	}
	C.wurfl_enum_destroy(eh)

	return result
}

/*
 *
 * Project : WURFL InFuze Golang module
 *
 * Author(s): Paul Stephen Borile
 *
 * Date: Aug 16 2016
 *
 * Copyright (c) ScientiaMobile, Inc.
 * http://www.scientiamobile.com
 *
 * This software package is the property of ScientiaMobile Inc. and is licensed
 * commercially according to a contract between the Licensee and ScientiaMobile Inc. (Licensor).
 * If you represent the Licensee, please refer to the licensing agreement which has been signed
 * between the two parties. If you do not represent the Licensee, you are not authorized to use
 * this software in any way.
 */

// Methods used for benchmarking : since they use cgo we cannot put them in test packages
// and They need to be public otherwise we won't be able to test them.

// GoStringToCStringAndFree converts a Go string to a C string and frees the memory.
func GoStringToCStringAndFree(capname string) *C.char {
	ccap := C.CString(capname)
	C.free(unsafe.Pointer(ccap))
	return ccap
}

// GoStringToCStringUsingMap returns a C string pointer for the given capability name using a cached map.
func (w *Wurfl) GoStringToCStringUsingMap(capname string) *C.char {
	return w.capsCStringcache[capname]
}

// CompareVersions Returns 0 if v1 == v2, -1 if v1 < v2, and 1 if v1 > v2.
func CompareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// Compare each part
	for i := 0; i < 4; i++ {
		// Convert string parts to integers for comparison
		num1, err := strconv.Atoi(parts1[i])
		if err != nil {
			fmt.Printf("Error converting version part '%s' to integer: %s\n", parts1[i], err)
			return 0
		}

		num2, err := strconv.Atoi(parts2[i])
		if err != nil {
			fmt.Printf("Error converting version part '%s' to integer: %s\n", parts2[i], err)
			return 0
		}

		// Compare parts
		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	// If all parts are equal
	return 0
}
