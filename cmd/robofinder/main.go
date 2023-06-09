package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func getRobotsURL(domain string) string {
	return domain + "/robots.txt"
}

func newClient(userAgent string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {

			return nil
		},
	}
}
func getURLResponse(url string, userAgent string) (*http.Response, error) {
	client := newClient(userAgent)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cookie", "some_cookie=some_value")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getURLBody(url string, userAgent string) (io.Reader, error) {
	client := newClient(userAgent)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cookie", "some_cookie=some_value")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	return resp.Body, nil
}

func getURLsFromReader(reader io.Reader) ([]string, error) {
	urls := make([]string, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Disallow:") || strings.HasPrefix(line, "Allow:") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				path := fields[1]
				urls = append(urls, path)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func getURLsDirect(domain string) ([]string, error) {
	url := getRobotsURL(domain)
	body, err := getURLBody(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)  (compatible; Kaya/1.0)")
	if err != nil {
		return nil, err
	}
	defer body.(io.ReadCloser).Close()
	return getURLsFromReader(body)
}

func getURLsArchive(domain string, archive string) ([]string, error) {
	url := archive + domain + "/robots.txt&output=json&filter=statuscode:200&fl=timestamp,original&collapse=digest"
	body, err := getURLBody(url, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	defer body.(io.ReadCloser).Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var records [][]string
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	for _, record := range records {
		if len(record) > 1 {
			timestamp := record[0]
			original := record[1]
			time.Sleep(500 * time.Millisecond)
			resp, err := getURLResponse("https://web.archive.org/web/"+timestamp+"/"+original, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
			if err != nil {
				return nil, err
			}
			if resp.StatusCode == http.StatusNotFound {
				continue
			}
			body := resp.Body
			defer body.(io.ReadCloser).Close()
			paths, err := getURLsFromReader(body)
			if err != nil {
				return nil, err
			}
			urls = append(urls, paths...)
		}
	}
	return urls, nil
}

func printPaths(urls []string, host string) {
	seen := make(map[string]bool)
	for _, url := range urls {
		if !seen[url] {
			fmt.Println(host + url)
			seen[url] = true
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: robofinder <domain>")
		os.Exit(1)
	}
	domain := os.Args[1]
	urls, err := getURLsDirect(domain)
	if err != nil {
		fmt.Println("Error:", err)
	}

	archiveURLs, err := getURLsArchive(domain, "https://web.archive.org/cdx/search/cdx?url=")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		urls = append(urls, archiveURLs...)
		printPaths(urls, domain)
	}
}
