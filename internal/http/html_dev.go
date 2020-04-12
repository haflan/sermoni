// +build !PRODUCTION

package http

import "io/ioutil"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// In dev mode, the file is read from the generated html file
func getWebsite() []byte {
	htmlData, err := ioutil.ReadFile("ui/dist/index.html")
	check(err)
	return htmlData
}
