package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var ERROR_NOT_FOUND = errors.New("not found")
var ERROR_TIMEOUT = errors.New("timeout")

func ResponseBodyToHTML(body io.ReadCloser) (string, error) {
	str, err := ioutil.ReadAll(body)

	if err != nil {
		return "", err
	}

	return string(str), nil
}

func MakeMockRequest(path string) (*http.Response, error) {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)

	return &http.Response{
		Body: ioutil.NopCloser(r),
	}, nil
}

func FormatTestError(expect, have interface{}) string {
	return fmt.Sprintf("\nexpect: %v, \nhave: %v", expect, have)
}

func NewError(fn string, err error) error {
	return fmt.Errorf("%s: %w", fn, err)
}

func NewRequest(url, method string, params url.Values, headers map[string]string) (*http.Response, error) {
	encoded := params.Encode()

	req, _ := http.NewRequest(method, url, strings.NewReader(encoded))

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := http.Client{}

	return client.Do(req)
}

func TimeoutRoutine(timeout chan error) {
	strTimeout := os.Getenv("GOMANGA_TIMEOUT")

	if strTimeout == "" {
		strTimeout = "10"
	}

	seconds, _ := strconv.Atoi(strTimeout)
	time.Sleep(time.Duration(seconds) * time.Second)
	timeout <- ERROR_TIMEOUT
}

func StrToInt(str string) (int, error) {
	part := strings.Split(str, ".")[0]
	return strconv.Atoi(part)
}

func GetTitleWithGreatestSimilarity(text string, titles []string) (title string) {
	var greatest float64 = 0

	for _, value := range titles {
		similarity := cosineSimilarity(text, value)

		if similarity > greatest {
			greatest = similarity
			title = value
		}
	}
	return
}

func cosineSimilarity(s1 string, s2 string) float64 {
	set1 := make(map[string]int)
	set2 := make(map[string]int)

	for _, char := range strings.Split(s1, "") {
		set1[char]++
	}

	for _, char := range strings.Split(s2, "") {
		set2[char]++
	}

	numerator := 0
	for key, value := range set1 {
		numerator += value * set2[key]
	}

	sum1 := 0
	for _, value := range set1 {
		sum1 += value * value
	}
	norm1 := math.Sqrt(float64(sum1))

	sum2 := 0
	for _, value := range set2 {
		sum2 += value * value
	}
	norm2 := math.Sqrt(float64(sum2))

	return float64(numerator) / (norm1 * norm2)
}
