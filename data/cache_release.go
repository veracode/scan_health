//go:build !cache

package data

func getCachedResponse(url string) []byte {
	return nil
}

func cacheResponse(url string, data []byte) {
	// Nothing to do
}
