package canalx

// Message 可以根据需要把其它字段也加入进来。
type Message[T any] struct {
	Data     []T    `json:"data"`
	Database string `json:"database"`
	Table    string `json:"table"`
	Type     string `json:"type"`
}
