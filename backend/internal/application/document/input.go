package document

type UploadDocumentInput struct {
	CaseID      string
	Name        string
	Description string
	Type        string
	FileName    string
	MimeType    string
	FileSize    int64
	StorageKey  string
}
