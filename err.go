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
	ErrInvalidHandle                       = errors.New("wurfl: (sentinel) invalid handle")
	ErrAlreadyLoad                         = errors.New("wurfl: (sentinel) already loaded")
	ErrFileNotFound                        = errors.New("wurfl: (sentinel) file not found")
	ErrUnexpectedEndOfFile                 = errors.New("wurfl: (sentinel) unexpected end of file")
	ErrInputOutputFailure                  = errors.New("wurfl: (sentinel) input output failure")
	ErrDeviceNotFound                      = errors.New("wurfl: (sentinel) device not found")
	ErrCapabilityNotFound                  = errors.New("wurfl: (sentinel) capability not found")
	ErrInvalidCapabilityValue              = errors.New("wurfl: (sentinel) invalid capability value")
	ErrVirtualCapabilityNotFound           = errors.New("wurfl: (sentinel) virtual capability not found")
	ErrCantLoadCapabilityNotFound          = errors.New("wurfl: (sentinel) can't load capability not found")
	ErrCantLoadVirtualCapabilityNotFound   = errors.New("wurfl: (sentinel) can't load virtual capability not found")
	ErrEmptyID                             = errors.New("wurfl: (sentinel) empty id")
	ErrCapabilityGroupNotFound             = errors.New("wurfl: (sentinel) capability group not found")
	ErrCapabilityGroupMismatch             = errors.New("wurfl: (sentinel) capability group mismatch")
	ErrDeviceAlreadyDefined                = errors.New("wurfl: (sentinel) device already defined")
	ErrUseragentAlreadyDefined             = errors.New("wurfl: (sentinel) useragent already defined")
	ErrDeviceHierarchyCircularReference    = errors.New("wurfl: (sentinel) device hierarchy circular reference")
	ErrUnknown                             = errors.New("wurfl: (sentinel) unknown error")
	ErrInvalidUseragentPriority            = errors.New("wurfl: (sentinel) invalid useragent priority")
	ErrInvalidParameter                    = errors.New("wurfl: (sentinel) invalid parameter")
	ErrInvalidCacheSize                    = errors.New("wurfl: (sentinel) invalid cache size")
	ErrXMLConsistency                      = errors.New("wurfl: (sentinel) xml consistency error")
	ErrInternal                            = errors.New("wurfl: (sentinel) internal error")
	ErrVirtualCapabilityNotAvailable       = errors.New("wurfl: (sentinel) virtual capability not available")
	ErrMissingUseragent                    = errors.New("wurfl: (sentinel) missing useragent")
	ErrXMLParse                            = errors.New("wurfl: (sentinel) xml parse error")
	ErrUpdaterInvalidDataURL               = errors.New("wurfl: (sentinel) updater invalid data url")
	ErrUpdaterInvalidLicense               = errors.New("wurfl: (sentinel) updater invalid license")
	ErrUpdaterNetworkError                 = errors.New("wurfl: (sentinel) updater network error")
	ErrEngineNotInitialized                = errors.New("wurfl: (sentinel) engine not initialized")
	ErrUpdaterAlreadyRunning               = errors.New("wurfl: (sentinel) updater already running")
	ErrUpdaterNotRunning                   = errors.New("wurfl: (sentinel) updater not running")
	ErrUpdaterTooManyRequests              = errors.New("wurfl: (sentinel) updater too many requests")
	ErrUpdaterCmdlineDownloaderUnavailable = errors.New("wurfl: (sentinel) updater cmdline downloader unavailable")
	ErrUpdaterTimedout                     = errors.New("wurfl: (sentinel) updater timed out")
	ErrRootNotSet                          = errors.New("wurfl: (sentinel) root not set")
	ErrWrongEngineTarget                   = errors.New("wurfl: (sentinel) wrong engine target")
	ErrCannotFilterStaticCap               = errors.New("wurfl: (sentinel) cannot filter static cap")
	ErrUnableToAllocateMemory              = errors.New("wurfl: (sentinel) unable to allocate memory")
	ErrEngineNotLoaded                     = errors.New("wurfl: (sentinel) engine not loaded")
	ErrUpdaterCannotStartThread            = errors.New("wurfl: (sentinel) updater cannot start thread")
	ErrEnumEmptySet                        = errors.New("wurfl: (sentinel) enum empty set")
	ErrUpdaterWrongDataFormat              = errors.New("wurfl: (sentinel) updater wrong data format")
	ErrUpdaterInvalidUseragent             = errors.New("wurfl: (sentinel) updater invalid useragent")
	ErrPermissionDenied                    = errors.New("wurfl: (sentinel) permission denied")
	ErrNotZipFile                          = errors.New("wurfl: (sentinel) not zip file")
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
