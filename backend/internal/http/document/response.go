package document

type UploadDocumentResponse struct {
	ID         string `json:"id"`
	CaseID     string `json:"case_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	FileName   string `json:"file_name"`
	MimeType   string `json:"mime_type"`
	FileSize   int64  `json:"file_size"`
	StorageKey string `json:"storage_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
