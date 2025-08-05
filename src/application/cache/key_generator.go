package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type KeyGenerator[T any, ID comparable] interface {
	GenerateKey(entity T) string
	GenerateKeyByID(id ID) string
	GenerateQueryKey(opts interface{}) string
	GetPrefix() string
}

type DefaultKeyGenerator[T any, ID comparable] struct {
	prefix    string
	getID     func(T) ID
	stringify func(ID) string
}

func NewKeyGenerator[T any, ID comparable](
	prefix string,
	getID func(T) ID,
	stringify func(ID) string,
) KeyGenerator[T, ID] {
	return &DefaultKeyGenerator[T, ID]{
		prefix:    prefix,
		getID:     getID,
		stringify: stringify,
	}
}

func (g *DefaultKeyGenerator[T, ID]) GenerateKey(entity T) string {
	return g.GenerateKeyByID(g.getID(entity))
}

func (g *DefaultKeyGenerator[T, ID]) GenerateKeyByID(id ID) string {
	return fmt.Sprintf("%s:%s", g.prefix, g.stringify(id))
}

func (g *DefaultKeyGenerator[T, ID]) GenerateQueryKey(opts interface{}) string {
	data, _ := json.Marshal(opts)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%s:query:%x", g.prefix, hash[:8])
}

func (g *DefaultKeyGenerator[T, ID]) GetPrefix() string {
	return g.prefix
}
