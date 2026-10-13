/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package manager

import (
	"context"
	"math/rand"
	"time"

	"github.com/TencentBlueKing/gopkg/cache"
	"github.com/TencentBlueKing/gopkg/cache/memory"
	"github.com/pkg/errors"

	"github.com/TencentBlueKing/bk-apigateway-sdks/v2/bkapi"
)

// PublicKeyProvider is the interface for public key provider.
type PublicKeyProvider interface {
	ProvidePublicKey(gatewayName string) (string, error)
}

// PublicKeySimpleProvider provides some predefined public keys.
type PublicKeySimpleProvider struct {
	publicKeys map[string]string
}

// ProvidePublicKey returns public key for given api name.
func (p *PublicKeySimpleProvider) ProvidePublicKey(gatewayName string) (string, error) {
	return p.publicKeys[gatewayName], nil
}

// NewPublicKeySimpleProvider creates a simple public key provider.
func NewPublicKeySimpleProvider(publicKeys map[string]string) *PublicKeySimpleProvider {
	return &PublicKeySimpleProvider{
		publicKeys: publicKeys,
	}
}

// PublicKeyMemoryCache will cache public key in memory.
type PublicKeyMemoryCache struct {
	cache memory.Cache
}

// ProvidePublicKey gets public key from cache.
func (c *PublicKeyMemoryCache) ProvidePublicKey(gatewayName string) (string, error) {
	return c.cache.GetString(context.Background(), cache.NewStringKey(gatewayName))
}

// NewPublicKeyMemoryCache creates a memory cache for public key.
func NewPublicKeyMemoryCache(
	config bkapi.Config,
	expiration time.Duration,
	clientFactory func(gatewayName string, config bkapi.Config) (*Manager, error),
) *PublicKeyMemoryCache {
	return &PublicKeyMemoryCache{
		cache: memory.NewCache(
			"public-key",
			func(ctx context.Context, key cache.Key) (any, error) {
				gatewayName := key.Key()
				manager, err := clientFactory(gatewayName, config)
				if err != nil {
					return nil, errors.WithMessagef(err, "failed to create manager for %s", gatewayName)
				}

				publicKey, err := manager.GetPublicKeyString(ctx)
				if err != nil {
					return nil, errors.WithMessagef(err, "failed to get public key for %s", gatewayName)
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

// NewDefaultPublicKeyMemoryCache creates a default memory cache for public key.
func NewDefaultPublicKeyMemoryCache(config bkapi.Config) *PublicKeyMemoryCache {
	return NewPublicKeyMemoryCache(config, 12*time.Hour, NewDefaultManager)
}
