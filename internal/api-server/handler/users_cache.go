package handler

import (
	"context"
	"sync"
	"time"

	"github.com/Donders-Institute/dr-gateway/internal/api-server/config"
	"github.com/Donders-Institute/dr-gateway/pkg/dr"
	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
)

// UsersCache is an in-memory store for caching all DR users.
type UsersCache struct {

	// Config is the general API server configuration.
	Config config.Configuration

	// Context is the API server context.
	Context context.Context

	store []*dr.DRUser
	mutex sync.RWMutex
}

// init initializes the cache with first reload.
func (c *UsersCache) Init() {

	// first refresh
	c.refresh()

	// every 10 minutes??
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Infof("refreshing users cache")
				c.refresh()
				log.Infof("users cache refreshed")
			case <-c.Context.Done():
				log.Infof("users cache refresh stopped")
				ticker.Stop()
				return
			}
		}
	}()

	log.Infof("users cache initalized")
}

// refresh update the cache with up-to-data users.
func (c *UsersCache) refresh() {

	if users, err := dr.GetAllUsers(c.Config.Dr); err != nil {
		log.Errorf("cannot refresh cache: %s", err.Error())
		return
	} else {
		// new users
		d := []*dr.DRUser{}

		for c := range users {
			d = append(d, c)
		}

		// set store to new data map
		c.mutex.Lock()
		c.store = d
		c.mutex.Unlock()
	}
}

// GetUsers returns all users from cache
func (c *UsersCache) GetUsers() []*dr.DRUser {

	c.mutex.RLock()
	users := c.store
	c.mutex.RUnlock()

	return users
}
