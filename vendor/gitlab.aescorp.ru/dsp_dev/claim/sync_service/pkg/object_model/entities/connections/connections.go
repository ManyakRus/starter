package connections

import (
	table_connections "gitlab.aescorp.ru/dsp_dev/claim/sync_service/pkg/db/tables/table_connections"
)

// Connection
type Connection struct {
	table_connections.Table_Connection
}
