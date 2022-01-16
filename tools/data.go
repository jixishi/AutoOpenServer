package tools

import (
	"archive/zip"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

func DeCompressZip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
