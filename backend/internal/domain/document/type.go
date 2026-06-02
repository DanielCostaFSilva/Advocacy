package document

type DocumentType string

const (
	DocumentTypePetition        DocumentType = "petition"
	DocumentTypeContract        DocumentType = "contract"
	DocumentTypeEvidence        DocumentType = "evidence"
	DocumentTypeCourtDecision   DocumentType = "court_decision"
	DocumentTypePowerOfAttorney DocumentType = "power_of_attorney"
	DocumentTypeOther           DocumentType = "other"
)
