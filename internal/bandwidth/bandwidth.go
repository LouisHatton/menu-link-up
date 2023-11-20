package bandwidth

type MonthlyBandwidth struct {
	ID                    string `json:"id"`
	UserId                string `json:"userId"`
	Month                 int    `json:"month"`
	Year                  int    `json:"year"`
	BytesTransferred      int64  `json:"bytesTransferred"`
	BytesTransferredLimit int64  `json:"bytesTransferredLimit"`
	BytesUploaded         int64  `json:"bytesUploaded"`
	BytesUploadedLimit    int64  `json:"bytesUploadedLimit"`
}
