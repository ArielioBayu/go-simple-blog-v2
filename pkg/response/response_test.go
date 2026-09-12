package response_test

import (
	"encoding/json"
	"testing"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
)

func TestPaginationResponse(t *testing.T) {
	resp := response.PaginationResponse{
		Status:  200,
		Message: "success get all post",
		Pagination: &response.Pagination{
			Limit:  10,
			Offset: 0,
		},
		Data: []string{"post1", "post2"},
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if result["status"] != float64(200) {
		t.Errorf("expected status 200, got %v", result["status"])
	}
	if result["message"] != "success get all post" {
		t.Errorf("expected message 'success get all post', got %v", result["message"])
	}
	if result["pagination"] == nil {
		t.Error("expected pagination not to be nil")
	}
	if result["data"] == nil {
		t.Error("expected data not to be nil")
	}
}

func TestDataResponse(t *testing.T) {
	resp := response.DataResponse{
		Status:  200,
		Message: "success get post",
		Data:    map[string]any{"id": 1, "title": "test"},
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if result["status"] != float64(200) {
		t.Errorf("expected status 200, got %v", result["status"])
	}
	if result["message"] != "success get post" {
		t.Errorf("expected message 'success get post', got %v", result["message"])
	}
	if result["data"] == nil {
		t.Error("expected data not to be nil")
	}
}

func TestMessageResponse(t *testing.T) {
	resp := response.MessageResponse{
		Status:  201,
		Message: "Success Create Post",
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if result["status"] != float64(201) {
		t.Errorf("expected status 201, got %v", result["status"])
	}
	if result["message"] != "Success Create Post" {
		t.Errorf("expected message 'Success Create Post', got %v", result["message"])
	}
}
