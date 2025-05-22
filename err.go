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
import "C" // Still needed here for C.WURFL_ERROR_LAST and other C constants if used directly

import (
	"errors"
	"fmt"
)

// Go error variables corresponding to wurfl_error enum
// These are sentinel values. Their messages are placeholders and typically won't be
// seen by the user directly if the error helpers in wurfl.go are used, as those
// will wrap these sentinels with messages fetched from the C library.
var (
	ErrInvalidHandle                       = errors.New("invalid handle")
	ErrAlreadyLoad                         = errors.New("wurfl_load has already been invoked on the specific wurfl_handle")
	ErrFileNotFound                        = errors.New("file not found during wurfl_load or remote data file update")
	ErrUnexpectedEndOfFile                 = errors.New("unexpected end of file or parsing error during wurfl_load")
	ErrInputOutputFailure                  = errors.New("error reading stream during wurfl_load or updater accessing local updated data file")
	ErrDeviceNotFound                      = errors.New("specified device is missing")
	ErrCapabilityNotFound                  = errors.New("specified capability is missing")
	ErrInvalidCapabilityValue              = errors.New("invalid capability value")
	ErrVirtualCapabilityNotFound           = errors.New("specified virtual capability is missing")
	ErrCantLoadCapabilityNotFound          = errors.New("specified capability is missing")
	ErrCantLoadVirtualCapabilityNotFound   = errors.New("specified virtual capability is missing")
	ErrEmptyID                             = errors.New("missing id in searching device")
	ErrCapabilityGroupNotFound             = errors.New("specified capability is missing in its group")
	ErrCapabilityGroupMismatch             = errors.New("specified capability mismatch in its group")
	ErrDeviceAlreadyDefined                = errors.New("specified device is already defined")
	ErrUseragentAlreadyDefined             = errors.New("specified user agent is already defined")
	ErrDeviceHierarchyCircularReference    = errors.New("circular reference in device hierarchy")
	ErrUnknown                             = errors.New("unknown error")
	ErrInvalidUseragentPriority            = errors.New("specified override sideloaded browser user agent configuration not valid")
	ErrInvalidParameter                    = errors.New("invalid parameter")
	ErrInvalidCacheSize                    = errors.New("specified an invalid cache size, 0 or a negative value.")
	ErrXMLConsistency                      = errors.New("WURFL data file is out of date or wrong - some needed device_id/capability is missing")
	ErrInternal                            = errors.New("internal error. For further info enable log using wurfl_set_log_path()")
	ErrVirtualCapabilityNotAvailable       = errors.New("the requested virtual capability has not been licensed")
	ErrMissingUseragent                    = errors.New("an XML device definition without mandatory UA has been detected")
	ErrXMLParse                            = errors.New("the XML data file is malformed")
	ErrUpdaterInvalidDataURL               = errors.New("updater data URL is missing or invalid (note: only .zip and .gz formats allowed)")
	ErrUpdaterInvalidLicense               = errors.New("client license is invalid or expired")
	ErrUpdaterNetworkError                 = errors.New("updater(curl) request returned an http error")
	ErrEngineNotInitialized                = errors.New("prerequisite for executing an update is that the engine has been initialized (i.e. wurfl_load() has been called)")
	ErrUpdaterAlreadyRunning               = errors.New("wurfl_updater_start() can be called just once when the updater is not running")
	ErrUpdaterNotRunning                   = errors.New("wurfl_updater_stop() can be called just once when the updater is running")
	ErrUpdaterTooManyRequests              = errors.New("Updater encountered HTTP 429 error, too many requests")
	ErrUpdaterCmdlineDownloaderUnavailable = errors.New("Curl executable not found. Please check path")
	ErrUpdaterTimedout                     = errors.New("Curl operation timed out.")
	ErrRootNotSet                          = errors.New("set_root() must be called before any load() / reload() and update attempt")
	ErrWrongEngineTarget                   = errors.New("set_engine_target() was called with a wrong/unrecognized parameter")
	ErrCannotFilterStaticCap               = errors.New("cannot filter static capability")
	ErrUnableToAllocateMemory              = errors.New("unable to allocate memory")
	ErrEngineNotLoaded                     = errors.New("engine not loaded, cannot return information")
	ErrUpdaterCannotStartThread            = errors.New("create of updater thread failed, enable updater logs for deatails")
	ErrEnumEmptySet                        = errors.New("enum set is empty")
	ErrUpdaterWrongDataFormat              = errors.New("unknown/unrecognized format in wurfl_updater_reload_root_from_memory()")
	ErrUpdaterInvalidUseragent             = errors.New("user_agent string passed to upd_set_useragent() is not valid")
	ErrPermissionDenied                    = errors.New("can't access file/directory - permission denied")
	ErrNotZipFile                          = errors.New("file isn't a zip")
)

// wurflGoErrors maps C wurfl_error codes to Go error values by index.
// The order and length must exactly match the wurfl_error enum in wurfl.h,
// up to WURFL_ERROR_LAST.
var wurflGoErrors []error

func init() {
	// C.WURFL_ERROR_LAST is the total number of error codes, including WURFL_OK (0).
	// The slice length should be C.WURFL_ERROR_LAST to accommodate indices from 0 to C.WURFL_ERROR_LAST-1.
	wurflGoErrors = make([]error, C.WURFL_ERROR_LAST)

	// wurflGoErrors[C.WURFL_OK] (index 0) is already nil by default, which represents no error.
	// Assign statically defined Go sentinel errors to the slice, indexed by their C enum value.
	wurflGoErrors[C.WURFL_ERROR_INVALID_HANDLE] = ErrInvalidHandle
	wurflGoErrors[C.WURFL_ERROR_ALREADY_LOAD] = ErrAlreadyLoad
	wurflGoErrors[C.WURFL_ERROR_FILE_NOT_FOUND] = ErrFileNotFound
	wurflGoErrors[C.WURFL_ERROR_UNEXPECTED_END_OF_FILE] = ErrUnexpectedEndOfFile
	wurflGoErrors[C.WURFL_ERROR_INPUT_OUTPUT_FAILURE] = ErrInputOutputFailure
	wurflGoErrors[C.WURFL_ERROR_DEVICE_NOT_FOUND] = ErrDeviceNotFound
	wurflGoErrors[C.WURFL_ERROR_CAPABILITY_NOT_FOUND] = ErrCapabilityNotFound
	wurflGoErrors[C.WURFL_ERROR_INVALID_CAPABILITY_VALUE] = ErrInvalidCapabilityValue
	wurflGoErrors[C.WURFL_ERROR_VIRTUAL_CAPABILITY_NOT_FOUND] = ErrVirtualCapabilityNotFound
	wurflGoErrors[C.WURFL_ERROR_CANT_LOAD_CAPABILITY_NOT_FOUND] = ErrCantLoadCapabilityNotFound
	wurflGoErrors[C.WURFL_ERROR_CANT_LOAD_VIRTUAL_CAPABILITY_NOT_FOUND] = ErrCantLoadVirtualCapabilityNotFound
	wurflGoErrors[C.WURFL_ERROR_EMPTY_ID] = ErrEmptyID
	wurflGoErrors[C.WURFL_ERROR_CAPABILITY_GROUP_NOT_FOUND] = ErrCapabilityGroupNotFound
	wurflGoErrors[C.WURFL_ERROR_CAPABILITY_GROUP_MISMATCH] = ErrCapabilityGroupMismatch
	wurflGoErrors[C.WURFL_ERROR_DEVICE_ALREADY_DEFINED] = ErrDeviceAlreadyDefined
	wurflGoErrors[C.WURFL_ERROR_USERAGENT_ALREADY_DEFINED] = ErrUseragentAlreadyDefined
	wurflGoErrors[C.WURFL_ERROR_DEVICE_HIERARCHY_CIRCULAR_REFERENCE] = ErrDeviceHierarchyCircularReference
	wurflGoErrors[C.WURFL_ERROR_UNKNOWN] = ErrUnknown
	wurflGoErrors[C.WURFL_ERROR_INVALID_USERAGENT_PRIORITY] = ErrInvalidUseragentPriority
	wurflGoErrors[C.WURFL_ERROR_INVALID_PARAMETER] = ErrInvalidParameter
	wurflGoErrors[C.WURFL_ERROR_INVALID_CACHE_SIZE] = ErrInvalidCacheSize
	wurflGoErrors[C.WURFL_ERROR_XML_CONSISTENCY] = ErrXMLConsistency
	wurflGoErrors[C.WURFL_ERROR_INTERNAL] = ErrInternal
	wurflGoErrors[C.WURFL_ERROR_VIRTUAL_CAPABILITY_NOT_AVAILABLE] = ErrVirtualCapabilityNotAvailable
	wurflGoErrors[C.WURFL_ERROR_MISSING_USERAGENT] = ErrMissingUseragent
	wurflGoErrors[C.WURFL_ERROR_XML_PARSE] = ErrXMLParse
	wurflGoErrors[C.WURFL_ERROR_UPDATER_INVALID_DATA_URL] = ErrUpdaterInvalidDataURL
	wurflGoErrors[C.WURFL_ERROR_UPDATER_INVALID_LICENSE] = ErrUpdaterInvalidLicense
	wurflGoErrors[C.WURFL_ERROR_UPDATER_NETWORK_ERROR] = ErrUpdaterNetworkError
	wurflGoErrors[C.WURFL_ERROR_ENGINE_NOT_INITIALIZED] = ErrEngineNotInitialized
	wurflGoErrors[C.WURFL_ERROR_UPDATER_ALREADY_RUNNING] = ErrUpdaterAlreadyRunning
	wurflGoErrors[C.WURFL_ERROR_UPDATER_NOT_RUNNING] = ErrUpdaterNotRunning
	wurflGoErrors[C.WURFL_ERROR_UPDATER_TOO_MANY_REQUESTS] = ErrUpdaterTooManyRequests
	wurflGoErrors[C.WURFL_ERROR_UPDATER_CMDLINE_DOWNLOADER_UNAVAILABLE] = ErrUpdaterCmdlineDownloaderUnavailable
	wurflGoErrors[C.WURFL_ERROR_UPDATER_TIMEDOUT] = ErrUpdaterTimedout
	wurflGoErrors[C.WURFL_ERROR_ROOT_NOT_SET] = ErrRootNotSet
	wurflGoErrors[C.WURFL_ERROR_WRONG_ENGINE_TARGET] = ErrWrongEngineTarget
	wurflGoErrors[C.WURFL_ERROR_CANNOT_FILTER_STATIC_CAP] = ErrCannotFilterStaticCap
	wurflGoErrors[C.WURFL_ERROR_UNABLE_TO_ALLOCATE_MEMORY] = ErrUnableToAllocateMemory
	wurflGoErrors[C.WURFL_ERROR_ENGINE_NOT_LOADED] = ErrEngineNotLoaded
	wurflGoErrors[C.WURFL_ERROR_UPDATER_CANNOT_START_THREAD] = ErrUpdaterCannotStartThread
	wurflGoErrors[C.WURFL_ERROR_ENUM_EMPTY_SET] = ErrEnumEmptySet
	wurflGoErrors[C.WURFL_ERROR_UPDATER_WRONG_DATA_FORMAT] = ErrUpdaterWrongDataFormat
	wurflGoErrors[C.WURFL_ERROR_UPDATER_INVALID_USERAGENT] = ErrUpdaterInvalidUseragent
	wurflGoErrors[C.WURFL_ERROR_PERMISSION_DENIED] = ErrPermissionDenied
	wurflGoErrors[C.WURFL_ERROR_NOT_ZIP_FILE] = ErrNotZipFile
}

// wurflError is a custom error type that wraps a sentinel error
// and provides a dynamic message from the C library.
type wurflError struct {
	sentinel error         // The base sentinel error (e.g., ErrDeviceNotFound) for errors.Is
	msg      string        // The actual message from the C library
	code     C.wurfl_error // The C error code, for reference
}

// Error returns the message fetched from the C library.
func (e *wurflError) Error() string {
	if e.msg == "" {
		// Fallback, though msg should always be populated by the constructors.
		return fmt.Sprintf("wurfl: uninitialized or unknown error (code %d)", e.code)
	}
	return e.msg
}

// Unwrap returns the underlying sentinel error, allowing errors.Is to work.
func (e *wurflError) Unwrap() error {
	return e.sentinel
}

// cErrorToGoError converts a C.wurfl_error (returned directly by a C func or via pointer) to a Go error.
// It uses the pre-initialized wurflGoErrors slice from err.go.
func cErrorToGoError(cErr C.wurfl_error) error {
	if cErr == C.WURFL_OK {
		return nil
	}

	// Get the actual error message from the C library for this code.
	actualCMsgChars := C.wurfl_get_error_string(cErr)
	var actualCMsg string
	if actualCMsgChars != nil {
		actualCMsg = C.GoString(actualCMsgChars)
	} else {
		// This case should ideally not happen for valid C error codes.
		actualCMsg = fmt.Sprintf("wurfl: undefined error message for code %d", cErr)
	}

	errCodeInt := int(cErr)
	// Check bounds and if the error is mapped in the slice from err.go
	if errCodeInt > 0 && errCodeInt < len(wurflGoErrors) && wurflGoErrors[errCodeInt] != nil {
		baseSentinelErr := wurflGoErrors[errCodeInt]
		// Return our custom error type, which will use actualCMsg for .Error()
		return &wurflError{sentinel: baseSentinelErr, msg: actualCMsg, code: cErr}
	}

	// Fallback for unmapped/new error codes - create a standard Go error.
	return fmt.Errorf("%s (code %d)", actualCMsg, cErr)
}

// checkHandleError checks the error state on a WURFL handle after an operation
// that doesn't directly return a wurfl_error but indicates failure via return value (e.g., NULL)
// and sets the error state on the handle.
func checkHandleError(handle C.wurfl_handle) error {
	if handle == nil {
		// This implies the Wurfl object itself (or its C handle) is nil.
		// We should return the specific sentinel for this, with its C message.
		cMsgChars := C.wurfl_get_error_string(C.WURFL_ERROR_INVALID_HANDLE)
		var msg string
		if cMsgChars != nil {
			msg = C.GoString(cMsgChars)
		} else {
			msg = "wurfl: (fallback) invalid handle" // Should match sentinel if C call fails
		}
		return &wurflError{sentinel: ErrInvalidHandle, msg: msg, code: C.WURFL_ERROR_INVALID_HANDLE}
	}

	errCode := C.wurfl_get_error_code(handle)
	if errCode == C.WURFL_OK {
		return nil
	}

	// Prioritize the specific runtime error message from the handle.
	specificRuntimeMsgChars := C.wurfl_get_error_message(handle)
	var finalMsg string
	if specificRuntimeMsgChars != nil {
		finalMsg = C.GoString(specificRuntimeMsgChars)
	}

	// If the specific message from the handle is empty, fall back to the generic string for the error code.
	if finalMsg == "" {
		genericMsgChars := C.wurfl_get_error_string(errCode)
		if genericMsgChars != nil {
			finalMsg = C.GoString(genericMsgChars)
		} else {
			// Ultimate fallback if both message sources fail.
			finalMsg = fmt.Sprintf("wurfl: undefined error message for code %d", errCode)
		}
	}

	// Try to get a pre-defined sentinel Go error for this code.
	errCodeInt := int(errCode)
	if errCodeInt > 0 && errCodeInt < len(wurflGoErrors) && wurflGoErrors[errCodeInt] != nil {
		baseSentinelErr := wurflGoErrors[errCodeInt]
		return &wurflError{sentinel: baseSentinelErr, msg: finalMsg, code: errCode}
	}

	// No sentinel error found (unmapped code) - create a standard Go error.
	return fmt.Errorf("%s (code %d)", finalMsg, errCode)
}
