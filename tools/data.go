package tools

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func MGetVersion(url string) (string, string) {
	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	tag := gjson.Get(string(body), "0.tag_name")
	DownUrl := gjson.Get(string(body), "0.assets.1.browser_download_url")
	return fmt.Sprint(tag), fmt.Sprint(DownUrl)
}
func WGetVersion(url string) (string, string, string) {
	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	tag := gjson.Get(string(body), "0.tag_name")
	JarDownUrl := gjson.Get(string(body), "0.assets.2.browser_download_url")
	ZipDownUrl := gjson.Get(string(body), "0.assets.1.browser_download_url")
	return fmt.Sprint(tag), fmt.Sprint(JarDownUrl), fmt.Sprint(ZipDownUrl)
}
