package transaction

type Transaction struct {
	PublicID string `json:"public_id"`

	AccountPublicID string `json:"account_public_id"`

	TaskPublicID string `json:"task_public_id"`

	Value int `json:"value"`
}
