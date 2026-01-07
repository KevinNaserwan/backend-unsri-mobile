package service

import (
	"fmt"
	"strings"
)

// SimpleGetServiceURL is a simple, working implementation
func SimpleGetServiceURL(path string) (string, error) {
	// Direct service mapping - guaranteed to work
	switch {
	case strings.HasPrefix(path, "/api/v1/auth"):
		return "http://auth-service:8081" + strings.TrimPrefix(path, "/api/v1/auth"), nil
	case strings.HasPrefix(path, "/api/v1/users"):
		return "http://user-service:8082" + strings.TrimPrefix(path, "/api/v1/users"), nil
	case strings.HasPrefix(path, "/api/v1/schedule"):
		return "http://schedule-service:8083" + strings.TrimPrefix(path, "/api/v1/schedule"), nil
	case strings.HasPrefix(path, "/api/v1/attendance"):
		return "http://attendance-service:8084" + strings.TrimPrefix(path, "/api/v1/attendance"), nil
	case strings.HasPrefix(path, "/api/v1/qr"):
		return "http://qr-service:8085" + strings.TrimPrefix(path, "/api/v1/qr"), nil
	case strings.HasPrefix(path, "/api/v1/broadcast"):
		return "http://broadcast-service:8086" + strings.TrimPrefix(path, "/api/v1/broadcast"), nil
	case strings.HasPrefix(path, "/api/v1/notifications"):
		return "http://notification-service:8087" + strings.TrimPrefix(path, "/api/v1/notifications"), nil
	case strings.HasPrefix(path, "/api/v1/calendar"):
		return "http://calendar-service:8088" + strings.TrimPrefix(path, "/api/v1/calendar"), nil
	case strings.HasPrefix(path, "/api/v1/courses"):
		return "http://course-service:8089" + strings.TrimPrefix(path, "/api/v1/courses"), nil
	case strings.HasPrefix(path, "/api/v1/locations"):
		return "http://location-service:8090" + strings.TrimPrefix(path, "/api/v1/locations"), nil
	case strings.HasPrefix(path, "/api/v1/access"):
		return "http://access-service:8091" + strings.TrimPrefix(path, "/api/v1/access"), nil
	case strings.HasPrefix(path, "/api/v1/quick-actions"):
		return "http://quick-actions-service:8092" + strings.TrimPrefix(path, "/api/v1/quick-actions"), nil
	case strings.HasPrefix(path, "/api/v1/files"):
		return "http://file-storage-service:8093" + strings.TrimPrefix(path, "/api/v1/files"), nil
	case strings.HasPrefix(path, "/api/v1/search"):
		return "http://search-service:8094" + strings.TrimPrefix(path, "/api/v1/search"), nil
	case strings.HasPrefix(path, "/api/v1/reports"):
		return "http://report-service:8095" + strings.TrimPrefix(path, "/api/v1/reports"), nil
	case strings.HasPrefix(path, "/api/v1/master-data"):
		return "http://master-data-service:8096" + strings.TrimPrefix(path, "/api/v1/master-data"), nil
	case strings.HasPrefix(path, "/api/v1/leave"):
		return "http://leave-service:8097" + strings.TrimPrefix(path, "/api/v1/leave"), nil
	default:
		return "", fmt.Errorf("no service found for path: %s", path)
	}
}
