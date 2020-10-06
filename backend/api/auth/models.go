package auth

type RequestJwt struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
	IssuedAt    int64  `json:"issuedAt"`
}
