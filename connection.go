package goutils

import "context"

// Connection is a struct that holds the connection information. It is used to connect to a database.
type Connection interface {
	// Connect to a database and return an error if it fails.
	// The name is the name of the connection defined in .env file. If the name is empty, the default connection will be used.
	Open(name ...string) error

	// Close the connection to the database.
	Close(name ...string) error

	// Get the connection to the database.
	Client(ctx ...context.Context) interface{}

	// Get the connection information.
	GetConfig(ctx ...context.Context) interface{}

	// Print all connection information from the .env file.
	Print(name ...string)
}

// When to get client with a connection name, must inject the key `conn_name` into Context,
// but the context doesn't allow to inject the key as string directly.
// Type `ctxConnNameKeyType` would be used in this case to replace type `string`.
// For example:
//
//	func GetClient(connName string) (Client, error) {
//		// goutils.CtxConnNameKey must be used instead of string "conn_name"
//		// `connName` is the connection name defined in .env file (e.g. REDIS_AGGS_URL, connName = "aggs")
//		ctx := context.WithValue(context.Background(), goutils.CtxConnNameKey, connName)
//		return Client(ctx)
//	}
const CtxConnNameKey ctxConnNameKeyType = "conn_name"

type ctxConnNameKeyType string
