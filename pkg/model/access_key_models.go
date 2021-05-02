package model

import "strings"

// AccessKey contains information about project's access key
type AccessKey struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	Key       string `json:"key"`
	ProjectID string `json:"-"`
}

// String converts an object to string
func (p AccessKey) String() string {
	return toJSON(&p)
}

// AccessKeyCreateParams contains parameters for access key creation
type AccessKeyCreateParams struct {
	Label string `json:"label"`
}

// String converts an object to string
func (p AccessKeyCreateParams) String() string {
	return toJSON(&p)
}

// Normalize normalizes request's fields
func (p *AccessKeyCreateParams) Normalize() {
	p.Label = strings.TrimSpace(p.Label)
}

// AccessKeys is a list of AccessKey
type AccessKeys []*AccessKey
