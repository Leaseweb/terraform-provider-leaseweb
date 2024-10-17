package contracts

// PublicCloudService gets data associated with public_cloud.
type PublicCloudService interface {
	// ValidateContractTerm checks if the passed combination of contractTerm & contractType is valid.
	ValidateContractTerm(contractTerm int64, contractType string) error
}
