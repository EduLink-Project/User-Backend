package interfaces

import "time"

type VC struct {
	Title         string
	Description   string
	StudentsLimit uint
	Instructorid  any
	Type          string
}

type RetrievedVC struct {
	ID            uint
	Title         string
	Description   string
	Instructor    string
	StudentsLimit uint
	Type          string
}

type DateTime struct {
	Date string
	Time string
}

type Sessions struct {
	Vcid               uint64
	QuestionsAllowance bool
	DateTime           DateTime
	Title              string
}

type RetrievedSessions struct {
	ID                 uint64
	Title              string
	Status             string
	DateTime           time.Time
	QuestionsAllowance bool
}

type Notifications struct {
	ID        uint64
	Type      string
	IsRead    bool
	Message   string
	Title     string
	CreatedAt time.Time
}
