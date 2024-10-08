package model

type RecordID string

type RecordType string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserID string

type RatingValue int

type Rating struct {
	RecordID   string      `json:`
	RecordType string      `json:`
	UserID     UserID      `json:`
	Value      RatingValue `json:`
}