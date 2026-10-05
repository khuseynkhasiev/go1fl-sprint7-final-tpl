package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
		url   string
		city  string
	}{
		{count: 0, want: 0, url: "/cafe", city: "moscow"},
		{count: 1, want: 1, url: "/cafe", city: "moscow"},
		{count: 2, want: 2, url: "/cafe", city: "moscow"},
		{count: 100, want: 100, url: "/cafe", city: "moscow"},
	}

	for _, v := range requests {
		var lengthResponse int

		response := httptest.NewRecorder()

		fullUrl := v.url + "?count=" + strconv.Itoa(v.count) + "&city=" + v.city
		req := httptest.NewRequest("GET", fullUrl, nil)

		handler.ServeHTTP(response, req)

		responseLine := strings.TrimSpace(response.Body.String())
		if responseLine != "" {
			partsResponse := strings.Split(responseLine, ",")
			lengthResponse = len(partsResponse)
		}

		totalCityCafe := len(cafeList[v.city])
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, min(v.want, totalCityCafe), lengthResponse)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		want   int
		search string
		city   string
		url    string
	}{
		{want: 0, search: "фасоль", city: "moscow", url: "/cafe"},
		{want: 2, search: "кофе", city: "moscow", url: "/cafe"},
		{want: 1, search: "вилка", city: "moscow", url: "/cafe"},
	}

	for _, v := range requests {
		var lengthResponse int
		var partsResponse []string

		response := httptest.NewRecorder()
		fullUrl := v.url + "?search=" + v.search + "&city=" + v.city

		req := httptest.NewRequest("GET", fullUrl, nil)

		handler.ServeHTTP(response, req)

		responseLine := strings.TrimSpace(response.Body.String())

		if responseLine != "" {
			partsResponse = strings.Split(responseLine, ",")
			lengthResponse = len(partsResponse)
		}

		for _, line := range partsResponse {
			assert.True(t, strings.Contains(strings.ToLower(line), strings.ToLower(v.search)), "incorrect search")
		}

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.want, lengthResponse)
	}
}
