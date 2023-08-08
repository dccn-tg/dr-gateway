package handler

import (
	"context"
	"sync"
	"time"

	"github.com/dccn-tg/dr-gateway/internal/api-server/config"
	"github.com/dccn-tg/dr-gateway/pkg/dr"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

// CollectionsCache is an in-memory store for caching all DR collections.
type CollectionsCache struct {

	// Config is the general API server configuration.
	Config config.Configuration

	// Context is the API server context.
	Context context.Context

	store []*dr.DRCollection
	mutex sync.RWMutex
}

// init initializes the cache with first reload.
func (c *CollectionsCache) Init() {

	// first refresh
	c.refresh()

	// every 10 minutes??
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infof("refreshing collections cache")
				c.refresh()
				log.Infof("collections cache refreshed")
			case <-c.Context.Done():
				log.Infof("collections cache refresh stopped")
				ticker.Stop()
				return
			}
		}
	}()

	log.Infof("collections cache initalized")
}

// refresh update the cache with up-to-data collections.
func (c *CollectionsCache) refresh() {

	if colls, err := dr.GetAllCollections(c.Config.Dr); err != nil {
		log.Errorf("cannot refresh cache: %s", err.Error())
		return
	} else {
		// new collections
		d := []*dr.DRCollection{}

		for c := range colls {
			d = append(d, c)
		}

		// set store to new data map
		c.mutex.Lock()
		c.store = d
		c.mutex.Unlock()
	}
}

// GetCollections returns all collections from cache
func (c *CollectionsCache) GetCollections() []*dr.DRCollection {

	c.mutex.RLock()
	colls := c.store
	c.mutex.RUnlock()

	return colls
}
