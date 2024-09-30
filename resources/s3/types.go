package s3_manager

type AwsS3infoInterface struct {
	Name         string            `json:"Name"`
	PublicAccess string            `json:"PublicAccess"`
	Tags         map[string]string `json:"Tags"`
	Notification string            `json:"Notification"`
}
