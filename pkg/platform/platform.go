package platform

type Platform string

const (
	Unknown Platform = "Unknown"
	Linux   Platform = "Linux"
	Windows Platform = "Windows"
)

var currentPlatform Platform = Unknown

func GetPlatform() Platform {
	return currentPlatform
}
