package document

import (
	"fmt"
	"time"

	documentDomain "legalflow/internal/domain/document"
)

const DownloadURLExpiration = 30 * time.Minute

type DownloadURLProvider interface {
	Generate(document *documentDomain.Document) string
	ExpiresAt() time.Time
}

type TemporaryDownloadURLProvider struct {
	expiresAt time.Time
}

func NewTemporaryDownloadURLProvider() *TemporaryDownloadURLProvider {
	return &TemporaryDownloadURLProvider{
		expiresAt: time.Now().Add(DownloadURLExpiration),
	}
}

func (p *TemporaryDownloadURLProvider) Generate(document *documentDomain.Document) string {
	return fmt.Sprintf("https://temporary-download.local/documents/%s", document.ID.String())
}

func (p *TemporaryDownloadURLProvider) ExpiresAt() time.Time {
	return p.expiresAt
}
