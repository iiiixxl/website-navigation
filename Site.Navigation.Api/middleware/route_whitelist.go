package middleware

import "sync"

var (
	anonymousRoutes = make(map[string]bool)
	mu              sync.RWMutex
)

func RegisterAnonymousRoute(method, path string) {
	mu.Lock()
	defer mu.Unlock()
	anonymousRoutes[method+" "+path] = true
}

func IsAnonymousRoute(method, path string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return anonymousRoutes[method+" "+path]
}
