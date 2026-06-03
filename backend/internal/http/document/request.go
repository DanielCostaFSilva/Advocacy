package document

const MaxDocumentSize = 20 * 1024 * 1024

var allowedMimeTypes = map[string]bool{
	"application/pdf":                                          true,
	"image/jpeg":                                               true,
	"image/png":                                                true,
	"application/msword":                                       true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}
