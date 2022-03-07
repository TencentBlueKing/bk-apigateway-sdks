package manager

import (
	"math/rand"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/bkapi"
	"github.com/TencentBlueKing/gopkg/cache"
	"github.com/TencentBlueKing/gopkg/cache/memory"
	"github.com/pkg/errors"
)

// PublicKeyProvider is the interface for public key provider.
type PublicKeyProvider interface {
	ProvidePublicKey(apiName string) (string, error)
}

// PublicKeyMemoryCache :
type PublicKeyMemoryCache struct {
	cache memory.Cache
}

// ProvidePublicKey :
func (c *PublicKeyMemoryCache) ProvidePublicKey(apiName string) (string, error) {
	return c.cache.GetString(cache.NewStringKey(apiName))
}

// NewPublicKeyMemoryCache :
func NewPublicKeyMemoryCache(
	config bkapi.ClientConfig,
	expiration time.Duration,
	clientFactory func(apiName string, config bkapi.ClientConfig) (*Manager, error),
) *PublicKeyMemoryCache {
	return &PublicKeyMemoryCache{
		cache: memory.NewCache(
			"public-key",
			false,
			func(key cache.Key) (interface{}, error) {
				apiName := key.Key()
				manager, err := clientFactory(apiName, config)
				if err != nil {
					return nil, errors.WithMessagef(err, "failed to create manager for %s", apiName)
				}

				publicKey, err := manager.GetPublicKeyString()
				if err != nil {
					return nil, errors.WithMessagef(err, "failed to get public key for %s", apiName)
				}

				return publicKey, nil
			},
			expiration,
			func() time.Duration {
				return time.Duration(rand.Intn(10000)) * time.Millisecond
			},
		),
	}
}

// NewDefaultPublicKeyMemoryCache :
func NewDefaultPublicKeyMemoryCache(
	config bkapi.ClientConfig, expiration time.Duration,
) *PublicKeyMemoryCache {
	return NewPublicKeyMemoryCache(config, expiration, NewDefaultManager)
}
