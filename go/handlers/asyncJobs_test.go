package handlers

import (
	"context"
	"darwin2/jobs"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestStartJobRejectsUnknownOperation(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/jobs/start/unknown", nil)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.SetPath("/jobs/start/:operation")
	c.SetParamNames("operation")
	c.SetParamValues("unknown")

	if err := HandleStartJob(c); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestJobHTTPHandlersExposeAndCancelJob(t *testing.T) {
	original := JobManager
	JobManager = jobs.NewManager(10, 1)
	t.Cleanup(func() { JobManager = original })

	started := make(chan struct{})
	created := JobManager.Start("test", "Test job", 1, "items", func(ctx context.Context, _ *jobs.Reporter) (string, error) {
		close(started)
		<-ctx.Done()
		return "", ctx.Err()
	})
	<-started

	e := echo.New()
	getRequest := httptest.NewRequest(http.MethodGet, "/jobs/"+created.ID, nil)
	getRecorder := httptest.NewRecorder()
	getContext := e.NewContext(getRequest, getRecorder)
	getContext.SetParamNames("id")
	getContext.SetParamValues(created.ID)
	if err := HandleGetJob(getContext); err != nil {
		t.Fatal(err)
	}
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d", getRecorder.Code)
	}

	cancelRequest := httptest.NewRequest(http.MethodDelete, "/jobs/"+created.ID, nil)
	cancelRecorder := httptest.NewRecorder()
	cancelContext := e.NewContext(cancelRequest, cancelRecorder)
	cancelContext.SetParamNames("id")
	cancelContext.SetParamValues(created.ID)
	if err := HandleCancelJob(cancelContext); err != nil {
		t.Fatal(err)
	}
	if cancelRecorder.Code != http.StatusAccepted {
		t.Fatalf("DELETE status = %d", cancelRecorder.Code)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		job, _ := JobManager.Get(created.ID)
		if job.Status == jobs.StatusCancelled {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("job was not cancelled")
}

func TestMachineLearningJobCountIsReadBeforeStart(t *testing.T) {
	request := capturedRequest{body: []byte(`{"sampleCount":3000}`)}
	operation, err := buildJobOperation("machine-learning", request)
	if err != nil {
		t.Fatal(err)
	}
	if operation.total != 3000 || operation.unit != "samples" {
		t.Fatalf("unexpected operation: %#v", operation)
	}
}

func TestExecuteHandlerCapturesResponse(t *testing.T) {
	request := capturedRequest{
		method:  http.MethodPost,
		path:    "/test",
		headers: http.Header{"Content-Type": []string{"application/json"}},
		body:    []byte(`{"value":"ok"}`),
	}
	result, err := executeHandler(context.Background(), request, func(c echo.Context) error {
		var body map[string]string
		if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
			return err
		}
		return c.String(http.StatusOK, strings.ToUpper(body["value"]))
	})
	if err != nil || result != "OK" {
		t.Fatalf("result = %q, err = %v", result, err)
	}
}
