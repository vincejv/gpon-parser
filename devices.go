package main

type OntDevice interface {
	GetOpticalInfo() *OpticalStats
	GetDeviceInfo() *DeviceStats
	GetGponUrl() string
	UpdateCachedPage()
}

type HG6245D_Globe struct {
}

type AN5506_Stock struct {
}

type ZTEF670L struct {
}

type ZLTG3000A struct {
}
