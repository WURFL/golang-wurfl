# Overview
golang-wurfl is a Go language binding to the WURFL InFuze C API.

# Getting Started
```
go get -u github.com/WURFL/golang-wurfl
```

# Usage
``` go
package main

import (
	"fmt"
	"os"

	wurfl "github.com/scientiamobile/golang-wurfl"
)

func main() {

	// Replace this with your own WURFL updater URL
	wurflUpdaterURL := "https://data.scientiamobile.com/xxxxx/wurfl.zip"

	fmt.Println("Downloading WURFL file ...")

	err := wurfl.Download(wurflUpdaterURL, ".")

	if err != nil {
		fmt.Printf("Error downloading WURFL file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("WURFL file downloaded successfully")

	wengine, err := wurfl.Create("./wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")

	if err != nil {
		fmt.Printf("Error creating WURFL engine: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Engine loaded, version : ", wengine.GetAPIVersion(), "wurfl info ", wengine.GetInfo())

	c := wengine.GetAllCaps()
	fmt.Println("Capabilities available = ", c)

	ua := "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36"

	device, err := wengine.LookupUserAgent(ua)
	if err != nil {
		fmt.Printf("Error in lookup: %v", err)
		os.Exit(1)
	}

	deviceid, err := device.GetDeviceID()
	if err != nil {
		fmt.Printf("Error in GetDeviceID: %v", err)
		os.Exit(1)
	}
	fmt.Println(deviceid)

	cap, _ := device.GetStaticCap("model_name")
	fmt.Printf("model_name = %s\n", cap)

	vcap, _ := device.GetVirtualCap("is_android")
	fmt.Printf("is_android = %s\n", vcap)

	if wengine.IsUserAgentFrozen(ua) {
		fmt.Printf("UA %s is frozen. Sec-Ch-Ua headers are necessary for correct device identification.\n", ua)
	}

	device.Destroy()
	wengine.Destroy()
}

```