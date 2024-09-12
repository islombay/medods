package status

type Status struct {
	Code    int               `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Count   int64             `json:"count,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Error   map[string]string `json:"error,omitempty"`
}

type StatusError string

var (
	ErrInvalid   = StatusError("invalid")
	ErrNotFound  = StatusError("not_found")
	ErrBadValue  = StatusError("bad_value")
	ErrUUID      = StatusError("uuid")
	ErrDuplicate = StatusError("duplicate")
)

func (s Status) AddError(key string, value StatusError) Status {
	if s.Error == nil {
		s.Error = map[string]string{}
	}
	s.Error[key] = string(value)
	return s
}

func (s Status) AddData(data interface{}) Status {
	s.Data = data
	return s
}

func (s Status) AddDataMap(key, val string) Status {
	if s.Data == nil {
		s.Data = map[string]interface{}{}
	}
	s.Data.(map[string]interface{})[key] = val
	return s
}

func (s Status) AddCode(code int) Status {
	s.Code = code
	return s
}

func (s Status) AddCount(count int64) Status {
	s.Count = count
	return s
}
