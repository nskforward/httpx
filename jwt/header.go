package jwt

type Header struct {
	A SignAlg `json:"alg"`
	T string  `json:"typ"`
}
