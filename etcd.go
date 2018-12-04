package main

import (
	"context"
	etcdclient "go.etcd.io/etcd/client"
	"gopkg.in/yaml.v2"
	"strings"
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

func (conn *dbConn) setYaml(key string, structObject interface{}) (*etcdclient.Response, error) {
	bytes, err := yaml.Marshal(&structObject)
	if err != nil {
		return nil, err
	}
	resp, err := conn.set(key, string(bytes))
	return resp, err
}

func (conn *dbConn) get(key string) (string, bool, error) {
	resp, err := conn.kApi.Get(context.Background(), key, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			return "", false, nil
		} else {
			return "", false, err
		}
	} else {
		return resp.Node.Value, true, nil
	}
}

// Populates pass-by-reference struct from Yaml in the corresponding key.
func (conn *dbConn) getStruct(key string, structObject interface{}) (bool, error) {
	specYaml, found, err := conn.get(key)
	if err != nil {
		return false, err
	} else if !found {
		return false, nil
	}
	err = yaml.Unmarshal([]byte(specYaml), structObject)
	if err != nil {
		return false, err
	}
	return true, nil
}
