package appjwt

type CertResponse struct {
	Keys []JWK `json:"keys"`
}
