package pkgutil

import "github.com/arfan21/project-sprint-social-media-api/config"

func GetPort(ports ...string) string {
	if len(ports) > 0 {
		return ":" + ports[0]
	}
	port := config.Get().HttpPort
	if port != "" {
		return ":" + port
	}
	return ":8888"
}
