package models

import (
	"github.com/qframe/types/metrics"
)

type HttpResponse struct {
	Failed int
	Success int
	Errors []HttpError
}
func NewHttpResponse() HttpResponse {
	return HttpResponse{
		Failed: 0,
		Success: 0,
		Errors: []HttpError{},
	}
}
type HttpError struct {
	Datapoint qtypes_metrics.OpenTSDBMetric
	Error string
}

func NewHttpError(dp qtypes_metrics.OpenTSDBMetric, error string) HttpError {
	return HttpError{
		Datapoint: dp,
		Error: error,
	}
}
