package main

import (
	"fmt"
	"os"

	wurfl "github.com/WURFL/golang-wurfl"
)

func main() {

	// Replace this with your own WURFL Snapshot URL
	wurflUpdaterURL := "https://data.scientiamobile.com/xxxxx/wurfl.zip"

	fmt.Println("Downloading WURFL file ...")

	if err := wurfl.Download(wurflUpdaterURL, "."); err != nil {
		fmt.Printf("Error downloading WURFL file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("WURFL file downloaded successfully")

	wengine, err := wurfl.Create("./wurfl.zip", nil, nil, -1, wurfl.WurflCacheProviderLru, "100000")

	if err != nil {
		fmt.Printf("Error creating WURFL engine: %v", err)
		os.Exit(1)
	}
	defer wengine.Destroy()

	fmt.Println("Engine loaded, version : ", wengine.GetAPIVersion(), "wurfl info ", wengine.GetInfo())

	// start the updater : will keep the wurfl.zip file updated to the last version
	// if err := wengine.SetUpdaterDataURL(wurflUpdaterURL); err != nil {
	// 	fmt.Printf("Error setting updater data URL: %v\n", err)
	// }
	// wengine.SetUpdaterDataFrequency(wurfl.WurflUpdaterFrequencyDaily)
	// wengine.UpdaterStart()

	c := wengine.GetAllCaps()
	fmt.Println("Capabilities available = ", c)

	ua := "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36"

	ihmap := make(map[string]string)
	ihmap["User-Agent"] = ua
	ihmap["Accept-Encoding"] = "gzip, deflate, br, zstd"
	ihmap["Sec-CH-UA-Platform"] = "Android"
	ihmap["Sec-CH-UA-Platform-Version"] = `"13.0.0"`
	ihmap["Sec-CH-UA"] = `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`
	ihmap["Sec-CH-UA-Mobile"] = "?1"
	ihmap["Sec-Ch-Ua-Model"] = `"SM-S135DL"`
	ihmap["Sec-Ch-Ua-Full-Version"] = `"126.0.6478.71"`
	ihmap["Sec-CH-UA-Full-Version-List"] = `"Not/A)Brand";v="8.0.0.0", "Chromium";v="126.0.6478.71", "Google Chrome";v="126.0.6478.71"`

	device, err := wengine.LookupWithImportantHeaderMap(ihmap)
	if err != nil {
		fmt.Printf("Error in lookup: %v", err)
		os.Exit(1)
	}
	defer device.Destroy()

	// obtain the wurfl_id for this device, a unique device identifier
	deviceid, err := device.GetDeviceID()
	if err != nil {
		fmt.Printf("Error in GetDeviceID: %v", err)
	}
	fmt.Println(deviceid)

	// static capabilities are stored in the wurfl.zip
	cap, _ := device.GetStaticCap("model_name")
	fmt.Printf("model_name = %s\n", cap)

	// virtual capabilities are computed on the fly
	vcap, _ := device.GetVirtualCap("complete_device_name")
	fmt.Printf("complete_device_name = %s\n", vcap)

	caps := []string{
		"marketing_name",
		"brand_name",
		"device_os",
	}
	// get a list of static caps, returns a map
	caplist, err := device.GetStaticCaps(caps)
	if err != nil {
		fmt.Printf("Error in GetDeviceID: %v", err)
	}
	fmt.Printf("marketing_name = %s\n", caplist["marketing_name"])
}
