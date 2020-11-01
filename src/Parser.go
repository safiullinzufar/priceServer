package src

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func getPrice(url *string) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", *url, nil)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("http request not sent")
	}
	buf,_ := ioutil.ReadAll(resp.Body)
	reg, _ := regexp.Compile("avito.item.price = '([0-9]+)'")
	arr := reg.FindStringIndex(string(buf))
	if len(arr) < 2 {
		log.Println(resp.StatusCode, *url)
		return "", errors.New("can not find price")
	}
	slicePrice := string(buf[arr[0]:arr[1]])
	reg, _ = regexp.Compile("([0-9]+)")
	arr = reg.FindStringIndex(slicePrice)
	if len(arr) < 2 {
		log.Println(resp.StatusCode, *url)
		return "", errors.New("can not find price")
	}
	price := slicePrice[arr[0]:arr[1]]
	return price, nil
}

