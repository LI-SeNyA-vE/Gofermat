package model

type Configs struct {
	ConfigFlag
	SecretKey
	HashKey
}

type SecretKey struct {
	SecretKey string
}

type HashKey struct {
	HashKey string
}

type ConfigFlag struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	SecretKeyForJWT      string
	SecretKeyForPassword string
}
