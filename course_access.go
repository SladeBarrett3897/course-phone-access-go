package coursephone

import (
	"context"
	"time"
)

type CodeRequest struct {
	Phone  string `json:"phone"`
	Locale string `json:"locale,omitempty"`
}

type VerifyRequest struct {
	Phone      string    `json:"phone"`
	Code       string    `json:"code"`
	CourseID   string    `json:"course_id"`
	Deadline   time.Time `json:"deadline"`
	LearnerNow time.Time `json:"learner_now"`
}

type AccessResult struct {
	CourseID string `json:"course_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

func (c *Client) SendSignupCode(ctx context.Context, input CodeRequest) error {
	_, err := c.post(ctx, "/v1/auth/phone/send_code", struct {
		Phone   string `json:"phone"`
		Purpose string `json:"purpose"`
		Locale  string `json:"locale,omitempty"`
	}{Phone: input.Phone, Purpose: "signup", Locale: input.Locale})
	if err != nil {
		return err
	}

	_, err = c.post(ctx, "/v1/sms/otp", struct {
		To string `json:"to"`
	}{To: input.Phone})
	return err
}

func (c *Client) VerifyCourseAccess(ctx context.Context, input VerifyRequest) (AccessResult, error) {
	_, err := c.post(ctx, "/v1/auth/phone/verify", struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
		Login bool   `json:"login"`
	}{Phone: input.Phone, Code: input.Code, Login: true})
	if err != nil {
		return AccessResult{}, err
	}
	return deadlineDecision(input.CourseID, input.Deadline, input.LearnerNow), nil
}

func deadlineDecision(courseID string, deadline, now time.Time) AccessResult {
	if now.After(deadline) {
		return AccessResult{CourseID: courseID, Status: "deadline_passed", Message: "contact your educator"}
	}
	return AccessResult{CourseID: courseID, Status: "course_open", Message: "continue course"}
}
