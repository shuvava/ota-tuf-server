package encryption

import "github.com/shuvava/ota-tuf-server/pkg/data"

// Key represents a common methods of different keys.
type Key interface {
	// Type returns the type of key.
	Type() data.KeyType
}
