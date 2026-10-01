package coursephone

import (
	"testing"
	"time"
)

func TestDeadlineDecision(t *testing.T) {
	deadline := time.Date(2026, time.September, 20, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name string
		now  time.Time
		want string
	}{
		{name: "learner can continue before cutoff", now: deadline.Add(-time.Minute), want: "course_open"},
		{name: "educator follow-up after cutoff", now: deadline.Add(time.Minute), want: "deadline_passed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deadlineDecision("go-infra-101", deadline, tt.now)
			if got.Status != tt.want {
				t.Fatalf("status = %q, want %q", got.Status, tt.want)
			}
		})
	}
}
