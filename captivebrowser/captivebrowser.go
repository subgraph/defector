
package main

import (
	"flag"
	"fmt"
        "io/ioutil"
        "regexp"
        "net/http"
	"os"

	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-webkit/webkit"
)

var url string
var debugMode bool
var isCaptivePortalOnStartup bool
var authenticated bool

// TODO: User-supplied captive portal test option
func init() {
	flag.StringVar(&url, "u", "http://clients3.google.com/generate_204", "Captive portal test URL")
	flag.BoolVar(&debugMode, "d", false, "Debug mode")
}

func httpFetch(experimentalUrl string) (*http.Response, error) {
        userAgent := "Mozilla/5.0 (Linux; Android 5.1; Elite 5 Build/LMY47D) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/39.0.0.0 Mobile Safari/537.36"
	client := &http.Client{}
        request, _ := http.NewRequest("GET", experimentalUrl, nil)
        request.Header.Set("User-Agent", userAgent)
        response, err := client.Do(request)
        if err != nil {
                return response, err
        }
        return response, nil
}

func httpContentMatch(response *http.Response, controlResult string, fuzzy bool) (bool) {
        defer response.Body.Close()
        responseContent, err := ioutil.ReadAll(response.Body)
        if err != nil {
                fmt.Println(err)
        }
        if fuzzy {
                regex, _ := regexp.Compile(controlResult)
                match := regex.FindStringIndex(string(responseContent))
                if match != nil {
                        return true
                }
        } else {
                if string(responseContent) == controlResult {
                        return true
                }
        }
        return false
}

func httpStatusCodeMatch(experimentalCode int, controlCode int) (bool) {
        return experimentalCode == controlCode
}

func DetectCaptivePortal() (bool) {
	response, err := httpFetch(url)
	if err != nil {
		panic(err)
	}
	statusMatch := httpStatusCodeMatch(204, response.StatusCode)
	contentMatch := httpContentMatch(response, "", false)
	return !(statusMatch && contentMatch)
}

func LaunchBrowser() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Connect to Captive Portal")
	window.SetDefaultSize(1000, 680)
	window.SetKeepAbove(true)
	window.Connect("destroy", gtk.MainQuit)
	
	vbox := gtk.NewVBox(false, 1)
	menuBox := gtk.NewHBox(false, 1)
	vbox.PackStart(menuBox, false, false, 0)
	toolbar := gtk.NewToolbar()
	toolbar.SetStyle(gtk.TOOLBAR_ICONS)
	
	backButton := gtk.NewToolButtonFromStock(gtk.STOCK_GO_BACK)
	backButton.SetSensitive(false)
	menuBox.PackStart(backButton, false, false, 0)
	
	forwardButton := gtk.NewToolButtonFromStock(gtk.STOCK_GO_FORWARD)
	menuBox.PackStart(forwardButton, false, false, 0)
	forwardButton.SetSensitive(false)
	
	connectButton := gtk.NewToolButtonFromStock(gtk.STOCK_CONNECT)
	menuBox.PackStart(connectButton, false, false, 0)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)

	webview := webkit.NewWebView()
	webview.SetMaintainsBackForwardList(true)
	swin.Add(webview)

	vbox.Add(swin)
	statusBox := gtk.NewHBox(false, 1)

	statusBar := gtk.NewStatusbar()
	contextId := statusBar.GetContextId("captive-browser")
	statusBox.PackStart(statusBar, true, true, 0)
	vbox.PackStart(statusBox, false, true, 0)
	connectButton.Connect("clicked", func() {
		statusBar.Push(contextId, "Connecting to captive portal...")
		webview.LoadUri(url)
		backButton.SetSensitive(true)
        })
	backButton.Connect("clicked", func() {
		if webview.CanGoBack() {
			statusBar.Push(contextId,
				"Connecting to previous page...")
			webview.GoBack()
		}
	})
	forwardButton.Connect("clicked", func() {
		if webview.CanGoForward() {
			statusBar.Push(contextId, "Connecting to next page...")
			webview.GoForward()
		}
	})
	webview.Connect("load-committed", func() {
		statusBar.Push(contextId, webview.GetUri())
		forwardButton.SetSensitive(webview.CanGoForward())
		backButton.SetSensitive(webview.CanGoBack())
	})
	webview.Connect("load-finished", func() {
		isCaptivePortal := DetectCaptivePortal()
		// TODO: Send a notification that captive portal authentication was successful
		if !isCaptivePortal && isCaptivePortalOnStartup {
			authenticated = true
			fmt.Println("Authentication successful")
		}
	})
	window.Add(vbox)
	window.SetSizeRequest(1000, 600)
	window.ShowAll()
	gtk.Main()
}

func main () {
	// TODO: in real life we want to exit with a notification that a captive portal is not available, though the 
	// dialog can present a retry option ("Retry vs. Cancel") -- these should continue until 'Cancel' is pressed.
        flag.Parse()
	isCaptivePortalOnStartup = DetectCaptivePortal()
	if debugMode {
		isCaptivePortalOnStartup = true
	}
	if isCaptivePortalOnStartup {
		LaunchBrowser()
	} else {
		os.Exit(1)
	}
	if authenticated {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
