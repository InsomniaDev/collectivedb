package api

type UpdateStruct struct {
	Key      string `json:"key,omitempty"`
	Database string `json:"database,omitempty"`
	Data     string `json:"data,omitempty"`
}

type DeleteStruct struct {
	Key      string `json:"key,omitempty"`
	Database string `json:"database,omitempty"`
}
