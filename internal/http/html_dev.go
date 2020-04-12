// +build !PRODUCTION

package http

import "io/ioutil"

// In production mode, the website is embedded in (generated) code 
// In dev mode it's more useful to read the html file on every request
func getWebsite() []byte {
	htmlData, err := ioutil.ReadFile("ui/dist/index.html")
	if err != nil {
		panic(err)
	}
	return htmlData
}
