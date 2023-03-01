package dr

import (
	"strconv"
	"sync"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/types"
)

type Config struct {
	IrodsHost           string
	IrodsPort           int
	IrodsZone           string
	IrodsUser           string
	IrodsPass           string
	OrganisationalUnits []string
}

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

func NewAccount(config Config) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSAccount(
		config.IrodsHost,
		config.IrodsPort,
		config.IrodsUser,
		config.IrodsZone,
		types.AuthSchemeNative,
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	// sslConfig, err := types.CreateIRODSSSLConfig(
	// 	"/opt/irods/ssl/icat-prod.pem",
	// 	32,
	// 	"AES-256-CBC",
	// 	8,
	// 	16,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// account.SSLConfiguration = sslConfig
	// account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	// account.ClientServerNegotiation = true

	return account, nil
}

func NewProxyAccount(config Config, user string) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSProxyAccount(
		config.IrodsHost,
		config.IrodsPort,
		user,
		config.IrodsZone,
		config.IrodsUser,
		config.IrodsZone,
		types.AuthSchemeNative,
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	// sslConfig, err := types.CreateIRODSSSLConfig(
	// 	"/opt/irods/ssl/icat-prod.pem",
	// 	32,
	// 	"AES-256-CBC",
	// 	8,
	// 	16,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// account.SSLConfiguration = sslConfig
	// account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	// account.ClientServerNegotiation = true

	return account, nil
}

func GetAllCollections(config Config) (chan *DRCollection, error) {
	acc, err := NewAccount(config)
	if err != nil {
		return nil, err
	}

	cpaths := make(chan string, 10000)
	go func() {

		defer close(cpaths)

		conn := connection.NewIRODSConnection(acc, 300*time.Second, "dr-gateway")
		if err := conn.Connect(); err != nil {
			log.Errorf("connection failure: %s", err)
			return
		}
		defer conn.Disconnect()

		for _, ou := range config.OrganisationalUnits {
			if _colls, err := fs.SearchCollectionsByMeta(conn, "organisationalUnit", ou); err != nil {
				log.Errorf(err.Error())
			} else {
				for _, c := range _colls {
					cpaths <- c.Path
				}
			}
		}
	}()

	return getCollections(acc, cpaths), nil
}

func FindCollectionsByMeta(config Config, key, value string) (chan *DRCollection, error) {
	//acc, err := dr.NewProxyAccount(os.Args[1])
	acc, err := NewAccount(config)
	if err != nil {
		return nil, err
	}

	cpaths := make(chan string, 10000)
	go func() {

		defer close(cpaths)

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
				cpaths <- c.Path
			}
		}
	}()

	return getCollections(acc, cpaths), nil
}

func getCollections(acc *types.IRODSAccount, cpaths chan string) chan *DRCollection {

	colls := make(chan *DRCollection, 1)

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
			for p := range cpaths {
				c := new(DRCollection)
				c.Path = p
				if meta, err := fs.ListCollectionMeta(conn, p); err != nil {
					log.Errorf(err.Error())
				} else {
					for _, m := range meta {

						switch m.Name {
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
