package auth

import ()

type authConfigUnit struct {
	IpAddr string
	MaxQps int
}

type authConfig struct {
	Auth map[string]*authConfigUnit
}
