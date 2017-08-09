package models

import "time"

type Annotation struct {
	tsuid 		string `json:"tsuid"`
	description string `json:"description"`
	notes 		string `json:"notes"`
	custom 		map[string]string	`json:"custom"`
	startTime 	time.Time
	endTime 	time.Time
}

func NewAnnotation(tsuid, description,notes string, startTime time.Time) Annotation {
	return Annotation{
		tsuid: tsuid,
		description: description,
		notes: notes,
		startTime: startTime,
		custom: map[string]string{},
	}
}

func (a *Annotation) AddTag(k, v string) {
	a.custom[k] = v
}
