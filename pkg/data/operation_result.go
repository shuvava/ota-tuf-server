package data

// OperationResult ????
type OperationResult struct {
	Target     TargetFilename          `json:"target"`
	Hashes     map[HashMethod]Checksum `json:"hashes"`
	Length     int64                   `json:"length"`
	ResultCode int                     `json:"resultCode"`
	ResultText string                  `json:"resultText"`
}

// IsSuccess returns true if operation successful
func (o OperationResult) IsSuccess() bool {
	return o.ResultCode == 0 || o.ResultCode == 1
}
