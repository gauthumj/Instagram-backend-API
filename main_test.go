package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPost(t *testing.T) {
	req, err := http.NewRequest("GET", "/posts/?id=61617db63ef00974f1839c80", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPosts)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `map[Caption:hi deer Timestamp:2021-10-09 17:02:06.0156016 +0530 IST m=+13.217800801 
	URL:iefbwiblauiebfveialraureasdcsdc Userid:ObjectID("61615dbff197baad0f8369f1") _id:ObjectID("61617db63ef00974f1839c80")]`

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreatePost(t *testing.T) {
	var jsonStr = []byte(`{"Caption":"hello","ImageURL":"test123.png","Userid":"61615dbff197baad0f8369f1"}`)
	req, err := http.NewRequest("POST", "posts", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createPost)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `Post created`
	if len(strings.TrimSpace(rr.Body.String())) != len(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"Name":"test","Email":"testuser@gmail.com","Password":"testpass"}`)
	req, err := http.NewRequest("POST", "users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `User created successfully!`
	if len(strings.TrimSpace(rr.Body.String())) != len(expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "users/?id=61615dbff197baad0f8369f1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[map[Email:abc@123 Name:nemesis Pass:fa237f9e776c81f5550a2f0f0f32d055628e6b99e7c1cfca48e15d919e7bb4820108 _id:ObjectID("616187a2563ecd91b0878b5c")]]
	`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPostsByUser(t *testing.T) {
	req, err := http.NewRequest("GET", "posts/users/?id=61615dbff197baad0f8369f1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userPosts)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `[map[Caption:hi deer Timestamp:2021-10-09 19:40:11.562941 +0530 IST m=+15.514119801 URL:iefbwiblauiebfveialraureasdcsdc Userid:ObjectID("61615dbff197baad0f8369f1") _id:ObjectID("6161a2c35c3bfd54914f8359")]]

	`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}