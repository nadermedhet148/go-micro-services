package config

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func DbConnect() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9999"}, // Default native port is 9000
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "admin",
			Password: "admin",
		},
		TLS: nil, // Uncomment and configure if you need TLS
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
