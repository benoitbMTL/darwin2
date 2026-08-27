package jobs

import (
	"context"
	"errors"
	"testing"
	"time"
)

func waitForTerminal(t *testing.T, manager *Manager, id string) Snapshot {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		snapshot, ok := manager.Get(id)
		if !ok {
			t.Fatalf("job %s disappeared", id)
		}
		if snapshot.Status == StatusSucceeded || snapshot.Status == StatusFailed || snapshot.Status == StatusCancelled {
			return snapshot
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("job %s did not finish", id)
	return Snapshot{}
}

func TestSuccessfulJobReportsProgress(t *testing.T) {
	manager := NewManager(10, 1)
	job := manager.Start("test", "Test", 2, "items", func(_ context.Context, reporter *Reporter) (string, error) {
		reporter.Progress(1, 2, "first")
		reporter.Progress(2, 2, "done")
		return "result", nil
	})

	finished := waitForTerminal(t, manager, job.ID)
	if finished.Status != StatusSucceeded || finished.Current != 2 || finished.Result != "result" {
		t.Fatalf("unexpected snapshot: %#v", finished)
	}
}

func TestJobCanBeCancelled(t *testing.T) {
	manager := NewManager(10, 1)
	started := make(chan struct{})
	job := manager.Start("test", "Test", 1, "items", func(ctx context.Context, _ *Reporter) (string, error) {
		close(started)
		<-ctx.Done()
		return "", ctx.Err()
	})
	<-started
	if _, ok := manager.Cancel(job.ID); !ok {
		t.Fatal("cancel did not find job")
	}
	finished := waitForTerminal(t, manager, job.ID)
	if finished.Status != StatusCancelled {
		t.Fatalf("status = %s, want cancelled", finished.Status)
	}
}

func TestFailedJobKeepsPartialErrors(t *testing.T) {
	manager := NewManager(10, 1)
	job := manager.Start("test", "Test", 3, "items", func(_ context.Context, reporter *Reporter) (string, error) {
		reporter.Progress(1, 3, "first")
		reporter.AddError("sample 2", errors.New("request failed"))
		return "", errors.New("stopped")
	})

	finished := waitForTerminal(t, manager, job.ID)
	if finished.Status != StatusFailed || finished.Error != "stopped" || len(finished.PartialErrors) != 1 {
		t.Fatalf("unexpected snapshot: %#v", finished)
	}
}

func TestWorkerLimitKeepsAdditionalJobQueued(t *testing.T) {
	manager := NewManager(10, 1)
	releaseFirst := make(chan struct{})
	firstStarted := make(chan struct{})
	first := manager.Start("test", "First", 1, "items", func(_ context.Context, _ *Reporter) (string, error) {
		close(firstStarted)
		<-releaseFirst
		return "", nil
	})
	<-firstStarted

	secondStarted := make(chan struct{})
	second := manager.Start("test", "Second", 1, "items", func(_ context.Context, _ *Reporter) (string, error) {
		close(secondStarted)
		return "", nil
	})
	secondSnapshot, _ := manager.Get(second.ID)
	if secondSnapshot.Status != StatusQueued {
		t.Fatalf("second status = %s, want queued", secondSnapshot.Status)
	}

	close(releaseFirst)
	waitForTerminal(t, manager, first.ID)
	select {
	case <-secondStarted:
	case <-time.After(time.Second):
		t.Fatal("second job did not start after worker became available")
	}
	waitForTerminal(t, manager, second.ID)
}
