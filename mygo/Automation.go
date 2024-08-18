package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BaseAutomationEvent represents the structure of the event data sent to the API
type BaseAutomationEvent struct {
	ID   string                 `json:"id"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data"`
}

// Point represents the structure of the response from the API
type Point struct {
	X    int                    `json:"x,omitempty"`
	Y    int                    `json:"y,omitempty"`
	Data map[string]interface{} `json:"data"`
}

// APIClient defines the client for the AutomationController API
type APIClient struct {
	BaseURL string
}

// CallAPI is a helper function to send POST requests to the AutomationController API
func (client *APIClient) CallAPI(endpoint string, event BaseAutomationEvent) (*Point, error) {
	url := client.BaseURL + endpoint
	jsonData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	var point Point
	if err := json.NewDecoder(resp.Body).Decode(&point); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	return &point, nil
}

// Test calls the /test endpoint
func (client *APIClient) Test() (*Point, error) {
	event := BaseAutomationEvent{
		Data: make(map[string]interface{}),
	}
	return client.CallAPI("/test", event)
}

// MouseMove calls the /mousemove endpoint
func (client *APIClient) MouseMove(x, y int) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"x": x,
			"y": y,
		},
	}
	return client.CallAPI("/mousemove", event)
}

// Click calls the /click endpoint
func (client *APIClient) Click(x, y int) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"x": x,
			"y": y,
		},
	}
	return client.CallAPI("/click", event)
}

// MoveWheel calls the /movewheel endpoint
func (client *APIClient) MoveWheel(amt int) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"amt": amt,
		},
	}
	return client.CallAPI("/movewheel", event)
}

// SendText calls the /sendtext endpoint
func (client *APIClient) SendText(x, y int, text string) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"x":    x,
			"y":    y,
			"text": text,
		},
	}
	return client.CallAPI("/sendtext", event)
}

// Type calls the /type endpoint
func (client *APIClient) Type(text string) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"text": text,
		},
	}
	return client.CallAPI("/type", event)
}

// FindText calls the /findtext endpoint
func (client *APIClient) FindText(text string) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"text": text,
		},
	}
	return client.CallAPI("/findtext", event)
}

// GetText calls the /gettext endpoint
func (client *APIClient) GetText(x, y, w, h int) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"x": x,
			"y": y,
			"w": w,
			"h": h,
		},
	}
	return client.CallAPI("/gettext", event)
}

// FindImage calls the /findimage endpoint
func (client *APIClient) FindImage(img string) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"img": img,
		},
	}
	return client.CallAPI("/findimage", event)
}

// Screenshot calls the /screenshot endpoint
func (client *APIClient) Screenshot() (*Point, error) {
	event := BaseAutomationEvent{
		Data: make(map[string]interface{}),
	}
	return client.CallAPI("/screenshot", event)
}

// GetMouseColor calls the /getmousecolor endpoint
func (client *APIClient) GetMouseColor(x, y int) (*Point, error) {
	event := BaseAutomationEvent{
		Data: map[string]interface{}{
			"x": x,
			"y": y,
		},
	}
	return client.CallAPI("/getmousecolor", event)
}

// GetMouse calls the /getmouse endpoint
func (client *APIClient) GetMouse() (*Point, error) {
	event := BaseAutomationEvent{
		Data: make(map[string]interface{}),
	}
	return client.CallAPI("/getmouse", event)
}
