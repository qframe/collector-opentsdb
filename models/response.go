package models

import 	"github.com/qnib/qframe-types"

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
	Datapoint qtypes.OpenTSDBMetric
	Error string
}

func NewHttpError(dp qtypes.OpenTSDBMetric, error string) HttpError {
	return HttpError{
		Datapoint: dp,
		Error: error,
	}
}
