package responsedata

type AuthResponseData struct {
	AccessToken string
	TokenType   string
}

type CreateSessionResponseData struct {
	SessionId string
}

type GetNodeContainerStats struct {
	CpuPercent     float64
	CpuUsed        float64
	CpuLimit       float64
	RamPercent     float64
	RamUsed        float64
	RamLimit       float64
	StoragePercent float64
	StorageUsed    float64
	StorageLimit   float64
}

type SendSingleDockerCommandResponseData struct {
	Output string
}

type ListDirectoryResponseData struct {
	Files []File
}

type File struct {
	Filename  string `json:"filename"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Size      int64  `json:"size"`
}
