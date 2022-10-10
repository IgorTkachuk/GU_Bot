package cache

type Repository interface {
	// GetIterator creates a new iterator for the cache.
	GetIterator() Iterator

	// Get returns value or not found error.
	Get(uuid []byte) ([]byte, error)

	// Set sets a key, value and expiration for a cache entry and stores it in the cache.
	Set(key []byte, val []byte, expireIn int) error

	// Del deletes an item in the cache by key and returns true or false if delete occurred.
	Del(key []byte) (affected bool)

	// EntryCount returns number of items currently in the cache.
	EntryCount() int64

	// HitCount is the metric that returns number of items a key was found in the cache.
	HitCount() int64

	// MissCount is the metric that returns number of items a miss occurred in the cache.
	MissCount() int64
}
