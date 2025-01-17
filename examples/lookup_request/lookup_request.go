package main

import (
	"fmt"
	"log"
	"net/http"
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
		fmt.Printf("Error creating WURFL engine: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Engine loaded, version : ", wengine.GetAPIVersion(), "wurfl info ", wengine.GetInfo())

	// start the updater : will keep the wurfl.zip file updated to the last version

	if err := wengine.SetUpdaterDataURL(wurflUpdaterURL); err != nil {
		fmt.Printf("Error setting updater data URL: %v\n", err)
	}
	wengine.SetUpdaterDataFrequency(wurfl.WurflUpdaterFrequencyDaily)
	wengine.UpdaterStart()

	// start server and process reqs

	fmt.Printf("Starting server on port 8080\n")

	// I need to inject the wurfl enginee into the handler
	// so that it can be used to lookup the device
	// for each request

	http.HandleFunc("/detect", makeHandler(wengine))

	log.Fatal(http.ListenAndServe(":8080", nil))

	wengine.Destroy()
}

// makeHandler creates a http.HandlerFunc injecting the wurfl engine
func makeHandler(wengine *wurfl.Wurfl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		device, err := wengine.LookupRequest(r)
		if err != nil {
			fmt.Fprintf(w, "Error in lookup: %v", err)
		}

		ado, _ := device.GetVirtualCap("advertised_device_os")
		isbot, _ := device.GetVirtualCap("is_robot")
		ab, _ := device.GetVirtualCap("advertised_browser")
		device.Destroy()

		fmt.Fprintf(w, "is_robot = %s\n", isbot)
		fmt.Fprintf(w, "advertised_device_os = %s\n", ado)
		fmt.Fprintf(w, "advertised_browser = %s\n", ab)
		// close response write and send out the response
		w.WriteHeader(http.StatusOK)

	}
}
