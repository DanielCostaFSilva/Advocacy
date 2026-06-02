package document

type UploadDocumentOutput struct {
	ID          string
	CaseID      string
	Name        string
	Description string
	Type        string
	FileName    string
	MimeType    string
	FileSize    int64
	StorageKey  string
}
