package sfn_manager

type (
	SfnInfoInterface struct {
		InputBucket  string `json:"InputBucket"`
		OutputBucket string `json:"OutputBucket"`
	}
)
