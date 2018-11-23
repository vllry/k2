package main

import (
	"context"
	etcdclient "go.etcd.io/etcd/client"
	"time"
)

type dbConn struct {
	conn etcdclient.Client
	kApi etcdclient.KeysAPI
}

func newDb() (*dbConn, error) {
	cfg := etcdclient.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: etcdclient.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	conn, err := etcdclient.New(cfg)
	if err != nil {
		return nil, err
	}
	return &dbConn{
		conn: conn,
		kApi: etcdclient.NewKeysAPI(conn),
	}, nil
}

func (conn *dbConn) set(key string, value string) (*etcdclient.Response, error) {
	return conn.kApi.Set(context.Background(), key, value, nil)
}

func (conn *dbConn) get(key string) (string, error) {
	resp, err := conn.kApi.Get(context.Background(), key, nil)
	if err != nil {
		return "", err
	} else {
		return resp.Node.Value, nil
	}
}
