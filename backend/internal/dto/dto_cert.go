package dto

import "time"

// CertSubmitRequest 提交资质请求。
type CertSubmitRequest struct {
	UserID       uint64     `json:"user_id" binding:"required"`
	CertType     string     `json:"cert_type" binding:"required,max=50"`
	CertNo       string     `json:"cert_no" binding:"max=50"`
	IssueOrg     string     `json:"issue_org" binding:"max=100"`
	IssueDate    *time.Time `json:"issue_date"`
	ValidUntil   *time.Time `json:"valid_until"`
	CertPhotoURL string     `json:"cert_photo_url" binding:"max=255"`
}

// CertReviewRequest 审核资质请求。
type CertReviewRequest struct {
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}
