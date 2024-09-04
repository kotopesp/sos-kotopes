package refreshsession

type RefreshSession struct {
	Fingerprint string `json:"fingerprint" form:"fingerprint" validate:"required,max=200"`
}
