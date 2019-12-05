package types

type TesseractConfig struct {
	Inbound struct {
		Action             string `json:"action"`     // Action can be { encrypt, decrypt, split, join }
		Encryption         string `json:"encryption"` // Encryption can be { password, pgp, aes, sha256, sha512 }
		PublicKeyFileName  string `json:"publicKeyFileName"`
		PrivateKeyFileName string `json:"privateKeyFileName"`
		Passphrase         string `json:"passphrase"`
	} `json:"inbound"`
	Outbound struct {
		Action     string `json:"action"`
		Encryption string `json:"encryption"`
	} `json:"outbound"`
}
