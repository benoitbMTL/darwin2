package jobs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"sync"
	"time"
)

type Status string

const (
	StatusQueued     Status = "queued"
	StatusRunning    Status = "running"
	StatusSucceeded  Status = "succeeded"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
	maxResultBytes          = 1 << 20
	maxPartialErrors        = 100
)

type PartialError struct {
	At      time.Time `json:"at"`
	Item    string    `json:"item,omitempty"`
	Message string    `json:"message"`
}

type Step struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	Status    Status    `json:"status"`
	Message   string    `json:"message,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Snapshot struct {
	ID                  string         `json:"id"`
	Type                string         `json:"type"`
	Name                string         `json:"name"`
	Status              Status         `json:"status"`
	Current             int64          `json:"current"`
	Total               int64          `json:"total"`
	Unit                string         `json:"unit"`
	Message             string         `json:"message,omitempty"`
	Result              string         `json:"result,omitempty"`
	Error               string         `json:"error,omitempty"`
	PartialErrors       []PartialError `json:"partialErrors"`
	Steps               []Step         `json:"steps,omitempty"`
	CreatedAt           time.Time      `json:"createdAt"`
	StartedAt           *time.Time     `json:"startedAt,omitempty"`
	FinishedAt          *time.Time     `json:"finishedAt,omitempty"`
	ElapsedMilliseconds int64          `json:"elapsedMilliseconds"`
	RatePerSecond       float64        `json:"ratePerSecond"`
	CancellationPending bool           `json:"cancellationPending"`
}

type Runner func(context.Context, *Reporter) (string, error)

type job struct {
	snapshot Snapshot
	cancel   context.CancelFunc
}

type Manager struct {
	mu          sync.RWMutex
	jobs        map[string]*job
	order       []string
	maxHistory  int
	workerSlots chan struct{}
	now         func() time.Time
	newID       func() string
}

func NewManager(maxHistory, maxConcurrent int) *Manager {
	if maxHistory <= 0 {
		maxHistory = 100
	}
	if maxConcurrent <= 0 {
		maxConcurrent = 4
	}
	return &Manager{
		jobs:        make(map[string]*job),
		maxHistory:  maxHistory,
		workerSlots: make(chan struct{}, maxConcurrent),
		now:         time.Now,
		newID:       randomID,
	}
}

func (m *Manager) Start(jobType, name string, total int64, unit string, runner Runner) Snapshot {
	ctx, cancel := context.WithCancel(context.Background())
	createdAt := m.now().UTC()
	j := &job{
		snapshot: Snapshot{
			ID:            m.newID(),
			Type:          jobType,
			Name:          name,
			Status:        StatusQueued,
			Total:         total,
			Unit:          unit,
			PartialErrors: []PartialError{},
			CreatedAt:     createdAt,
		},
		cancel: cancel,
	}

	m.mu.Lock()
	m.jobs[j.snapshot.ID] = j
	m.order = append([]string{j.snapshot.ID}, m.order...)
	m.pruneLocked()
	snapshot := m.snapshotLocked(j, createdAt)
	m.mu.Unlock()

	go m.run(ctx, j, runner)
	return snapshot
}

func (m *Manager) run(ctx context.Context, j *job, runner Runner) {
	select {
	case m.workerSlots <- struct{}{}:
		defer func() { <-m.workerSlots }()
	case <-ctx.Done():
		m.finish(j, StatusCancelled, "", context.Canceled)
		return
	}

	m.mu.Lock()
	now := m.now().UTC()
	j.snapshot.Status = StatusRunning
	j.snapshot.StartedAt = &now
	m.mu.Unlock()

	reporter := &Reporter{manager: m, job: j}
	result, err := runner(ctx, reporter)
	status := StatusSucceeded
	if errors.Is(ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		status = StatusCancelled
	} else if err != nil {
		status = StatusFailed
	}
	m.finish(j, status, result, err)
}

func (m *Manager) finish(j *job, status Status, result string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now().UTC()
	j.snapshot.Status = status
	if len(result) > maxResultBytes {
		result = result[:maxResultBytes] + "\n\n[Result truncated at 1 MiB]"
	}
	j.snapshot.Result = result
	j.snapshot.FinishedAt = &now
	j.snapshot.CancellationPending = false
	j.cancel()
	if err != nil && !errors.Is(err, context.Canceled) {
		j.snapshot.Error = err.Error()
	}
	if status == StatusSucceeded && j.snapshot.Total > 0 && j.snapshot.Current < j.snapshot.Total {
		j.snapshot.Current = j.snapshot.Total
	}
	m.pruneLocked()
}

func (m *Manager) Get(id string) (Snapshot, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	if !ok {
		return Snapshot{}, false
	}
	return m.snapshotLocked(j, m.now().UTC()), true
}

func (m *Manager) List(jobType string, limit int) []Snapshot {
	if limit <= 0 || limit > m.maxHistory {
		limit = m.maxHistory
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := m.now().UTC()
	result := make([]Snapshot, 0, limit)
	for _, id := range m.order {
		j := m.jobs[id]
		if jobType != "" && j.snapshot.Type != jobType {
			continue
		}
		result = append(result, m.snapshotLocked(j, now))
		if len(result) == limit {
			break
		}
	}
	return result
}

func (m *Manager) Cancel(id string) (Snapshot, bool) {
	m.mu.Lock()
	j, ok := m.jobs[id]
	if !ok {
		m.mu.Unlock()
		return Snapshot{}, false
	}
	if j.snapshot.Status == StatusQueued || j.snapshot.Status == StatusRunning {
		j.snapshot.CancellationPending = true
		j.snapshot.Message = "Cancellation requested"
		j.cancel()
	}
	snapshot := m.snapshotLocked(j, m.now().UTC())
	m.mu.Unlock()
	return snapshot, true
}

func (m *Manager) pruneLocked() {
	if len(m.order) <= m.maxHistory {
		return
	}
	kept := make([]string, 0, len(m.order))
	for _, id := range m.order {
		j := m.jobs[id]
		terminal := j.snapshot.Status == StatusSucceeded || j.snapshot.Status == StatusFailed || j.snapshot.Status == StatusCancelled
		if len(kept) >= m.maxHistory && terminal {
			delete(m.jobs, id)
			continue
		}
		kept = append(kept, id)
	}
	m.order = kept
}

func (m *Manager) snapshotLocked(j *job, now time.Time) Snapshot {
	snapshot := j.snapshot
	snapshot.PartialErrors = append([]PartialError(nil), j.snapshot.PartialErrors...)
	snapshot.Steps = append([]Step(nil), j.snapshot.Steps...)
	start := snapshot.CreatedAt
	if snapshot.StartedAt != nil {
		start = *snapshot.StartedAt
	}
	end := now
	if snapshot.FinishedAt != nil {
		end = *snapshot.FinishedAt
	}
	if end.Before(start) {
		end = start
	}
	duration := end.Sub(start)
	snapshot.ElapsedMilliseconds = duration.Milliseconds()
	if duration > 0 {
		snapshot.RatePerSecond = float64(snapshot.Current) / duration.Seconds()
	}
	return snapshot
}

type Reporter struct {
	manager *Manager
	job     *job
}

type reporterContextKey struct{}

func WithReporter(ctx context.Context, reporter *Reporter) context.Context {
	return context.WithValue(ctx, reporterContextKey{}, reporter)
}

func ReporterFrom(ctx context.Context) *Reporter {
	reporter, _ := ctx.Value(reporterContextKey{}).(*Reporter)
	return reporter
}

func ReportProgress(ctx context.Context, current, total int64, message string) {
	if reporter := ReporterFrom(ctx); reporter != nil {
		reporter.Progress(current, total, message)
	}
}

func ReportError(ctx context.Context, item string, err error) {
	if reporter := ReporterFrom(ctx); reporter != nil {
		reporter.AddError(item, err)
	}
}

func (r *Reporter) Progress(current, total int64, message string) {
	r.manager.mu.Lock()
	defer r.manager.mu.Unlock()
	if current >= 0 {
		r.job.snapshot.Current = current
	}
	if total >= 0 {
		r.job.snapshot.Total = total
	}
	r.job.snapshot.Message = message
}

func (r *Reporter) AddError(item string, err error) {
	if err == nil {
		return
	}
	r.manager.mu.Lock()
	defer r.manager.mu.Unlock()
	r.job.snapshot.PartialErrors = append(r.job.snapshot.PartialErrors, PartialError{
		At:      r.manager.now().UTC(),
		Item:    item,
		Message: err.Error(),
	})
	if len(r.job.snapshot.PartialErrors) > maxPartialErrors {
		r.job.snapshot.PartialErrors = r.job.snapshot.PartialErrors[len(r.job.snapshot.PartialErrors)-maxPartialErrors:]
	}
}

func (r *Reporter) Step(id, label string, status Status, message string) {
	r.manager.mu.Lock()
	defer r.manager.mu.Unlock()
	now := r.manager.now().UTC()
	for i := range r.job.snapshot.Steps {
		if r.job.snapshot.Steps[i].ID == id {
			r.job.snapshot.Steps[i].Label = label
			r.job.snapshot.Steps[i].Status = status
			r.job.snapshot.Steps[i].Message = message
			r.job.snapshot.Steps[i].UpdatedAt = now
			return
		}
	}
	r.job.snapshot.Steps = append(r.job.snapshot.Steps, Step{
		ID: id, Label: label, Status: status, Message: message, UpdatedAt: now,
	})
	sort.SliceStable(r.job.snapshot.Steps, func(i, j int) bool {
		return r.job.snapshot.Steps[i].UpdatedAt.Before(r.job.snapshot.Steps[j].UpdatedAt)
	})
}

func randomID() string {
	var value [12]byte
	if _, err := rand.Read(value[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405.000000000")))
	}
	return hex.EncodeToString(value[:])
}
