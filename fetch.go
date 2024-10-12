package main

import (
	"fmt"
	"net/http"
)

func fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP get request error: %w", err)
	}
	return resp, nil
}
