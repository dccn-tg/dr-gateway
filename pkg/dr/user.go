package dr

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/cyverse/go-irodsclient/irods/common"
	"github.com/cyverse/go-irodsclient/irods/connection"
	"github.com/cyverse/go-irodsclient/irods/fs"
	"github.com/cyverse/go-irodsclient/irods/message"
	"github.com/cyverse/go-irodsclient/irods/types"
	"golang.org/x/xerrors"
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

		if _users, err := listUsers(conn); err != nil {
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
						case "email":
							u.Email = m.Value
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

// list users is a copy from https://github.com/cyverse/go-irodsclient/blob/bf359d5f171e080a9ca6beac240bef0a6802dd8a/irods/fs/usergroup.go#L310
// with a fix to avoid `CAT_STATEMENT_TABLE_FULL` (due to the `continueIndex` in the `message.NewIRODSMessageQueryRequest`)
func listUsers(conn *connection.IRODSConnection) ([]*types.IRODSUser, error) {
	if conn == nil || !conn.IsConnected() {
		return nil, xerrors.Errorf("connection is nil or disconnected")
	}

	// lock the connection
	conn.Lock()
	defer conn.Unlock()

	users := []*types.IRODSUser{}

	continueQuery := true
	continueIndex := 0
	for continueQuery {
		query := message.NewIRODSMessageQueryRequest(common.MaxQueryRows, continueIndex, 0, 0)
		query.AddSelect(common.ICAT_COLUMN_USER_ID, 1)
		query.AddSelect(common.ICAT_COLUMN_USER_NAME, 1)
		query.AddSelect(common.ICAT_COLUMN_USER_TYPE, 1)
		query.AddSelect(common.ICAT_COLUMN_USER_ZONE, 1)

		condTypeVal := fmt.Sprintf("<> '%s'", types.IRODSUserRodsGroup)
		query.AddCondition(common.ICAT_COLUMN_USER_TYPE, condTypeVal)

		queryResult := message.IRODSMessageQueryResponse{}
		err := conn.Request(query, &queryResult, nil)
		if err != nil {
			return nil, xerrors.Errorf("failed to receive a user query result message: %w", err)
		}

		err = queryResult.CheckError()
		if err != nil {
			if types.GetIRODSErrorCode(err) == common.CAT_NO_ROWS_FOUND {
				// empty
				return users, nil
			}
			return nil, xerrors.Errorf("received a user query error: %w", err)
		}

		if queryResult.RowCount == 0 {
			break
		}

		if queryResult.AttributeCount > len(queryResult.SQLResult) {
			return nil, xerrors.Errorf("failed to receive user attributes - requires %d, but received %d attributes", queryResult.AttributeCount, len(queryResult.SQLResult))
		}

		pagenatedUsers := make([]*types.IRODSUser, queryResult.RowCount)

		for attr := 0; attr < queryResult.AttributeCount; attr++ {
			sqlResult := queryResult.SQLResult[attr]
			if len(sqlResult.Values) != queryResult.RowCount {
				return nil, xerrors.Errorf("failed to receive user rows - requires %d, but received %d attributes", queryResult.RowCount, len(sqlResult.Values))
			}

			for row := 0; row < queryResult.RowCount; row++ {
				value := sqlResult.Values[row]

				if pagenatedUsers[row] == nil {
					// create a new
					pagenatedUsers[row] = &types.IRODSUser{
						ID:   -1,
						Zone: "",
						Name: "",
						Type: types.IRODSUserRodsUser,
					}
				}

				switch sqlResult.AttributeIndex {
				case int(common.ICAT_COLUMN_USER_ID):
					userID, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return nil, xerrors.Errorf("failed to parse user id '%s': %w", value, err)
					}
					pagenatedUsers[row].ID = userID
				case int(common.ICAT_COLUMN_USER_ZONE):
					pagenatedUsers[row].Zone = value
				case int(common.ICAT_COLUMN_USER_NAME):
					pagenatedUsers[row].Name = value
				case int(common.ICAT_COLUMN_USER_TYPE):
					pagenatedUsers[row].Type = types.IRODSUserType(value)
				default:
					// ignore
				}
			}
		}

		users = append(users, pagenatedUsers...)

		continueIndex = queryResult.ContinueIndex
		if continueIndex == 0 {
			continueQuery = false
		}
	}

	return users, nil
}
