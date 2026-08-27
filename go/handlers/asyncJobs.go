package handlers

import (
	"bytes"
	"context"
	"darwin2/jobs"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

var JobManager = jobs.NewManager(100, 4)

type capturedRequest struct {
	method  string
	path    string
	headers http.Header
	body    []byte
}

type jobOperation struct {
	jobType string
	name    string
	total   int64
	unit    string
	runner  jobs.Runner
}

func HandleStartJob(c echo.Context) error {
	body, err := io.ReadAll(http.MaxBytesReader(c.Response(), c.Request().Body, 2<<20))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to read request body"})
	}
	request := capturedRequest{
		method:  http.MethodPost,
		path:    c.Request().URL.Path,
		headers: c.Request().Header.Clone(),
		body:    body,
	}
	operation, err := buildJobOperation(c.Param("operation"), request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	job := JobManager.Start(operation.jobType, operation.name, operation.total, operation.unit, operation.runner)
	c.Response().Header().Set("Location", "/jobs/"+job.ID)
	return c.JSON(http.StatusAccepted, job)
}

func HandleListJobs(c echo.Context) error {
	limit := 20
	if raw := c.QueryParam("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "limit must be a positive integer"})
		}
		limit = parsed
	}
	return c.JSON(http.StatusOK, JobManager.List(c.QueryParam("type"), limit))
}

func HandleGetJob(c echo.Context) error {
	job, ok := JobManager.Get(c.Param("id"))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job not found"})
	}
	return c.JSON(http.StatusOK, job)
}

func HandleCancelJob(c echo.Context) error {
	job, ok := JobManager.Cancel(c.Param("id"))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Job not found"})
	}
	return c.JSON(http.StatusAccepted, job)
}

func buildJobOperation(name string, request capturedRequest) (jobOperation, error) {
	switch name {
	case "machine-learning":
		var payload MLRequest
		if err := json.Unmarshal(request.body, &payload); err != nil {
			return jobOperation{}, errors.New("invalid machine learning request")
		}
		if payload.SampleCount <= 0 {
			payload.SampleCount = 10
		}
		return legacyJob("machine-learning", "Web ML traffic", int64(payload.SampleCount), "samples", request, HandleMachineLearning), nil
	case "api-traffic-generation":
		values, err := url.ParseQuery(string(request.body))
		if err != nil {
			return jobOperation{}, errors.New("invalid API traffic request")
		}
		count, err := strconv.Atoi(values.Get("count"))
		if err != nil || count <= 0 {
			return jobOperation{}, errors.New("count must be a positive integer")
		}
		return legacyJob("api-traffic-generation", "API ML traffic", int64(count), "samples", request, HandleApiMachineLearning), nil
	case "web-scan":
		return legacyJob("web-scan", "Nikto web scan", 1, "scans", request, HandleWebScan), nil
	case "traffic-generation":
		return legacyJob("traffic-generation", "Nikto traffic generation", 100, "scans", request, HandleTrafficGenerator), nil
	case "selenium":
		var payload requestParams
		if err := json.Unmarshal(request.body, &payload); err != nil {
			return jobOperation{}, errors.New("invalid Selenium request")
		}
		if payload.LoopCount <= 0 {
			payload.LoopCount = 1
		}
		return seleniumJob(request, payload.LoopCount), nil
	case "fortiweb-create":
		return fortiWebJob(true, request), nil
	case "fortiweb-delete":
		return fortiWebJob(false, request), nil
	case "health-check":
		return legacyJob("health-check", "Configuration health check", 8, "checks", request, HandleHealthCheck), nil
	default:
		return jobOperation{}, fmt.Errorf("unknown job operation %q", name)
	}
}

func legacyJob(jobType, name string, total int64, unit string, request capturedRequest, handler echo.HandlerFunc) jobOperation {
	return jobOperation{
		jobType: jobType,
		name:    name,
		total:   total,
		unit:    unit,
		runner: func(ctx context.Context, reporter *jobs.Reporter) (string, error) {
			reporter.Progress(0, total, "Starting")
			result, err := executeHandler(jobs.WithReporter(ctx, reporter), request, handler)
			if err == nil {
				reporter.Progress(total, total, "Completed")
			}
			return result, err
		},
	}
}

func seleniumJob(request capturedRequest, loopCount int) jobOperation {
	return jobOperation{
		jobType: "selenium",
		name:    "Selenium browser actions",
		total:   int64(loopCount),
		unit:    "loops",
		runner: func(ctx context.Context, reporter *jobs.Reporter) (string, error) {
			results := make([]string, 0, loopCount)
			for i := 1; i <= loopCount; i++ {
				if err := ctx.Err(); err != nil {
					return strings.Join(results, "\n"), err
				}
				reporter.Progress(int64(i-1), int64(loopCount), fmt.Sprintf("Running browser loop %d of %d", i, loopCount))
				result, err := executeHandler(jobs.WithReporter(ctx, reporter), request, HandleSelenium)
				if err != nil {
					reporter.AddError(fmt.Sprintf("loop %d", i), err)
					results = append(results, fmt.Sprintf("Loop %d: %v", i, err))
				} else {
					results = append(results, fmt.Sprintf("Loop %d: %s", i, strings.TrimSpace(result)))
				}
				reporter.Progress(int64(i), int64(loopCount), fmt.Sprintf("Completed browser loop %d of %d", i, loopCount))
			}
			return strings.Join(results, "\n"), nil
		},
	}
}

type fortiWebTask struct {
	id      string
	label   string
	path    string
	handler echo.HandlerFunc
}

func fortiWebJob(create bool, original capturedRequest) jobOperation {
	tasks := createFortiWebTasks()
	jobType := "fortiweb-create"
	name := "Create FortiWeb objects"
	if !create {
		tasks = deleteFortiWebTasks()
		jobType = "fortiweb-delete"
		name = "Delete FortiWeb objects"
	}

	return jobOperation{
		jobType: jobType,
		name:    name,
		total:   int64(len(tasks)),
		unit:    "objects",
		runner: func(ctx context.Context, reporter *jobs.Reporter) (string, error) {
			results := make([]string, 0, len(tasks))
			for i, task := range tasks {
				if err := ctx.Err(); err != nil {
					return strings.Join(results, "\n"), err
				}
				reporter.Step(task.id, task.label, jobs.StatusRunning, "")
				reporter.Progress(int64(i), int64(len(tasks)), task.label)
				request := original
				request.path = task.path
				result, err := executeHandler(jobs.WithReporter(ctx, reporter), request, task.handler)
				if err == nil {
					var apiResult ApiResult
					if decodeErr := json.Unmarshal([]byte(result), &apiResult); decodeErr != nil {
						err = fmt.Errorf("invalid FortiWeb response: %w", decodeErr)
					} else if apiResult.Status != "success" {
						err = errors.New(apiResult.Message)
					}
				}
				if err != nil {
					reporter.Step(task.id, task.label, jobs.StatusFailed, err.Error())
					reporter.AddError(task.label, err)
					results = append(results, fmt.Sprintf("%s: %v", task.label, err))
				} else {
					reporter.Step(task.id, task.label, jobs.StatusSucceeded, "")
					results = append(results, fmt.Sprintf("%s: %s", task.label, strings.TrimSpace(result)))
				}
				reporter.Progress(int64(i+1), int64(len(tasks)), task.label)
			}
			return strings.Join(results, "\n"), nil
		},
	}
}

func executeHandler(ctx context.Context, request capturedRequest, handler echo.HandlerFunc) (string, error) {
	req := httptest.NewRequest(request.method, request.path, bytes.NewReader(request.body)).WithContext(ctx)
	req.Header = request.headers.Clone()
	recorder := httptest.NewRecorder()
	e := echo.New()
	ec := e.NewContext(req, recorder)

	if err := handler(ec); err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return strings.TrimSpace(recorder.Body.String()), fmt.Errorf("HTTP %d: %v", httpError.Code, httpError.Message)
		}
		return strings.TrimSpace(recorder.Body.String()), err
	}
	if recorder.Code < http.StatusOK || recorder.Code >= http.StatusMultipleChoices {
		return strings.TrimSpace(recorder.Body.String()), fmt.Errorf("HTTP %d: %s", recorder.Code, strings.TrimSpace(recorder.Body.String()))
	}
	return strings.TrimSpace(recorder.Body.String()), nil
}

func createFortiWebTasks() []fortiWebTask {
	return []fortiWebTask{
		{"createNewVirtualIP", "Create new Virtual IP", "/create-virtual-ip", HandleCreateNewVirtualIP},
		{"createNewServerPool", "Create new Server Pool", "/create-server-pool", HandleCreateNewServerPool},
		{"createNewMemberPool", "Create new Member Pool", "/create-member-pool", HandleCreateNewMemberPool},
		{"createNewVirtualServer", "Create new Virtual Server", "/create-virtual-server", HandleCreateNewVirtualServer},
		{"assignVIPToVirtualServer", "Assign Virtual IP to Virtual Server", "/assign-vip-to-virtual-server", HandleAssignVIPToVirtualServer},
		{"cloneSignatureProtection", "Clone Signature Protection", "/clone-signature-protection", HandleCloneSignatureProtection},
		{"cloneInlineProtection", "Clone Inline Protection", "/clone-inline-protection", HandleCloneInlineProtection},
		{"createNewXForwardedForRule", "Create X-Forwarded-For Rule", "/create-x-forwarded-for-rule", HandleCreateNewXForwardedForRule},
		{"configureProtectionProfile", "Configure Protection Profile", "/configure-protection-profile", HandleConfigureProtectionProfile},
		{"createNewPolicy", "Create new Policy", "/create-policy", HandleCreateNewPolicy},
	}
}

func deleteFortiWebTasks() []fortiWebTask {
	return []fortiWebTask{
		{"deletePolicy", "Delete Policy", "/delete-policy", HandleDeletePolicy},
		{"deleteInlineProtection", "Delete Inline Protection Profile", "/delete-inline-protection", HandleDeleteInlineProtection},
		{"deleteXForwardedForRule", "Delete X-Forwarded-For Rule", "/delete-x-forwarded-for-rule", HandleDeleteXForwardedForRule},
		{"deleteSignatureProtection", "Delete Signature Protection", "/delete-signature-protection", HandleDeleteSignatureProtection},
		{"deleteVirtualServer", "Delete Virtual Server", "/delete-virtual-server", HandleDeleteVirtualServer},
		{"deleteServerPool", "Delete Server Pool", "/delete-server-pool", HandleDeleteServerPool},
		{"deleteVirtualIP", "Delete Virtual IP", "/delete-virtual-ip", HandleDeleteVirtualIP},
	}
}
