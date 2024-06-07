package dr

import (
	"strconv"
	"sync"
	"time"

	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

type CollectionState int

const (
	UnknownState CollectionState = iota
	Editable
	ReviewableInternal
	ReviewableExternal
	Archived
	Published
)

func (c CollectionState) String() string {
	switch c {
	case Editable:
		return "EDITABLE"
	case ReviewableInternal:
		return "REVIEWABLE_INTERNAL"
	case ReviewableExternal:
		return "REVIEWABLE_EXTERNAL"
	case Archived:
		return "ARCHIVED"
	case Published:
		return "PUBLISHED"
	default:
		return "UNKNOWN"
	}
}

func NewCollectionState(stat string) CollectionState {
	switch stat {
	case "EDITABLE":
		return Editable
	case "REVIEWABLE_INTERNAL":
		return ReviewableInternal
	case "REVIEWABLE_EXTERNAL":
		return ReviewableExternal
	case "ARCHIVED":
		return Archived
	case "PUBLISHED":
		return Published
	default:
		return UnknownState
	}
}

type CollectionType int

const (
	UnknownType CollectionType = iota
	DAC
	RDC
	DSC
)

func (t CollectionType) String() string {
	switch t {
	case DAC:
		return "DATA_ACQUISITION"
	case RDC:
		return "RESEARCH_DOCUMENTATION"
	case DSC:
		return "DATA_SHARING"
	default:
		return "UNKNOWN"
	}
}

func NewCollectionType(t string) CollectionType {
	switch t {
	case "DATA_ACQUISITION":
		return DAC
	case "RESEARCH_DOCUMENTATION":
		return RDC
	case "DATA_SHARING":
		return DSC
	default:
		return UnknownType
	}
}

type DRCollection struct {
	Identifier         string
	IdentifierDOI      string
	ProjectID          string
	Path               string
	Type               CollectionType
	State              CollectionState
	OrganisationalUnit string
	QuotaInBytes       int64
	SizeInBytes        int64
	NumberOfFiles      int64
}

type ouPath struct {
	ou   string
	path string
}

func GetAllCollections(config Config) (chan *DRCollection, error) {

	accounts, err := NewServiceAccounts(config)
	if err != nil {
		return nil, err
	}

	cpaths := make(chan ouPath, 10000)
	go func() {

		defer close(cpaths)

		for _, ou := range config.OrganisationalUnits {

			acc := accounts[ou.Name]

			conn := connection.NewIRODSConnection(acc, 300*time.Second, "dr-gateway")
			if err := conn.Connect(); err != nil {
				log.Errorf("connection failure: %s", err)
				return
			}
			defer conn.Disconnect()

			if _colls, err := fs.SearchCollectionsByMeta(conn, "organisationalUnit", ou.Name); err != nil {
				log.Errorf("%s\n", err.Error())
			} else {
				for _, c := range _colls {
					cpaths <- ouPath{
						ou.Name,
						c.Path,
					}
				}
			}
		}
	}()

	return getCollections(accounts, cpaths), nil
}

func FindCollectionsByMeta(config Config, key, value string) (chan *DRCollection, error) {

	accounts, err := NewServiceAccounts(config)
	if err != nil {
		return nil, err
	}

	cpaths := make(chan ouPath, 10000)
	go func() {

		defer close(cpaths)

		for _, ou := range config.OrganisationalUnits {
			acc := accounts[ou.Name]

			conn := connection.NewIRODSConnection(acc, 300*time.Second, "dr-gateway")
			if err := conn.Connect(); err != nil {
				log.Errorf("connection failure: %s", err)
				return
			}
			defer conn.Disconnect()

			if _colls, err := fs.SearchCollectionsByMeta(conn, key, value); err != nil {
				log.Errorf(err.Error())
			} else {
				for _, c := range _colls {
					cpaths <- ouPath{
						ou.Name,
						c.Path,
					}
				}
			}
		}
	}()

	return getCollections(accounts, cpaths), nil
}

func getCollections(accounts map[string]*types.IRODSAccount, cpaths chan ouPath) chan *DRCollection {

	colls := make(chan *DRCollection, 1)

	// workers to get collection metadata
	wg := sync.WaitGroup{}
	// start concurrent workers to get project resources from the filer.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// OU-specific connection map.
			// connections will be initialized later whenever needed.
			// initialized connections is kept in the map for reuse.
			conns := make(map[string]*connection.IRODSConnection)

			// closing up all connections
			defer func() {
				for _, conn := range conns {
					conn.Disconnect()
				}
			}()

		collLoop:
			for p := range cpaths {
				c := new(DRCollection)
				c.Path = p.path

				if _, ok := conns[p.ou]; !ok {
					conn := connection.NewIRODSConnection(accounts[p.ou], 300*time.Second, "dr-gateway")
					if err := conn.Connect(); err != nil {
						log.Errorf(err.Error())
						continue collLoop
					} else {
						conns[p.ou] = conn
					}
				}

				if meta, err := fs.ListCollectionMeta(conns[p.ou], p.path); err != nil {
					log.Errorf(err.Error())
				} else {
					for _, m := range meta {

						switch m.Name {
						case "deleted":
							if m.Value == "true" {
								// skip the deleted collection
								continue collLoop
							}
						case "identifierDOI":
							c.IdentifierDOI = m.Value
						case "collectionIdentifier":
							c.Identifier = m.Value
						case "projectId":
							c.ProjectID = m.Value
						case "state":
							c.State = NewCollectionState(m.Value)
						case "type":
							c.Type = NewCollectionType(m.Value)
						case "organisationalUnit":
							c.OrganisationalUnit = m.Value
						case "quotaInBytes":
							if s, err := strconv.Atoi(m.Value); err == nil {
								c.QuotaInBytes = int64(s)
							}
						case "sizeInBytes":
							if s, err := strconv.Atoi(m.Value); err == nil {
								c.SizeInBytes = int64(s)
							}
						case "numberOfFiles":
							if s, err := strconv.Atoi(m.Value); err == nil {
								c.NumberOfFiles = int64(s)
							}
						default:
						}
					}
					colls <- c
				}
			}
		}()
	}

	// wait for all workers to finish the work, and close the colls channel
	go func() {
		wg.Wait()
		close(colls)
	}()

	return colls
}
