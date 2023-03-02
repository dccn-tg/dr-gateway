package dr

import (
	"sync"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
)

type DRUser struct {
	DisplayName         string
	IdentityProvider    string
	Email               string
	OrganizationalUnits []string
}

func GetAllUsers(config Config) (chan *DRUser, error) {
	acc, err := NewAccount(config)
	if err != nil {
		return nil, err
	}

	uids := make(chan string, 10000)
	go func() {

		defer close(uids)

		conn := connection.NewIRODSConnection(acc, 300*time.Second, "dr-gateway")
		if err := conn.Connect(); err != nil {
			log.Errorf("connection failure: %s", err)
			return
		}
		defer conn.Disconnect()

		if _users, err := fs.ListUsers(conn); err != nil {
			log.Errorf("%s\n", err.Error())
		} else {
			for _, u := range _users {
				uids <- u.Name
			}
		}
	}()

	return getUsers(acc, uids), nil
}

func getUsers(acc *types.IRODSAccount, uids chan string) chan *DRUser {

	users := make(chan *DRUser, 1)

	// workers to get collection metadata
	wg := sync.WaitGroup{}
	// start concurrent workers to get project resources from the filer.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn := connection.NewIRODSConnection(acc, 300*time.Second, "dr-gateway")
			if err := conn.Connect(); err != nil {
				log.Errorf(err.Error())
				return
			}
			defer conn.Disconnect()
			for id := range uids {
				u := new(DRUser)
				if meta, err := fs.ListUserMeta(conn, id); err != nil {
					log.Errorf(err.Error())
				} else {
					for _, m := range meta {

						switch m.Name {
						case "customDisplayName":
							u.DisplayName = m.Value
						case "displayName":
							if u.DisplayName == "" {
								u.DisplayName = m.Value
							}
						case "identityProvider":
							u.IdentityProvider = m.Value
						case "organisationalUnit":
							u.OrganizationalUnits = append(u.OrganizationalUnits, m.Value)
						default:
						}
					}
					users <- u
				}
			}
		}()
	}

	// wait for all workers to finish the work, and close the users channel
	go func() {
		wg.Wait()
		close(users)
	}()

	return users
}
