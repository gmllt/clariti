package utils

import (
	"sync"
)

// JSONResponsePool is a pool for reusing JSON response maps to reduce allocations
var JSONResponsePool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{}, 8) // Pre-allocate with capacity
	},
}

// StringMapPool is a pool for reusing string maps
var StringMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string, 4)
	},
}

// GetJSONResponse gets a JSON response map from the pool
func GetJSONResponse() map[string]interface{} {
	m := JSONResponsePool.Get().(map[string]interface{})
	// Clear the map for reuse
	for k := range m {
		delete(m, k)
	}
	return m
}

// PutJSONResponse returns a JSON response map to the pool
func PutJSONResponse(m map[string]interface{}) {
	if len(m) > 16 { // Don't pool very large maps
		return
	}
	JSONResponsePool.Put(m)
}

// GetStringMap gets a string map from the pool
func GetStringMap() map[string]string {
	m := StringMapPool.Get().(map[string]string)
	// Clear the map for reuse
	for k := range m {
		delete(m, k)
	}
	return m
}

// PutStringMap returns a string map to the pool
func PutStringMap(m map[string]string) {
	if len(m) > 8 { // Don't pool very large maps
		return
	}
	StringMapPool.Put(m)
}
