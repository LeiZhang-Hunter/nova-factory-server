package milvus

import (
	"context"
	"errors"
	"testing"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"google.golang.org/grpc"
)

type fakeMilvusDatabaseClient struct {
	listDatabases []string
	listErr       error
	createErr     error
	useErr        error

	created []string
	used    []string
}

func (f *fakeMilvusDatabaseClient) ListDatabase(context.Context, milvusclient.ListDatabaseOption, ...grpc.CallOption) ([]string, error) {
	return f.listDatabases, f.listErr
}

func (f *fakeMilvusDatabaseClient) CreateDatabase(_ context.Context, opt milvusclient.CreateDatabaseOption, _ ...grpc.CallOption) error {
	f.created = append(f.created, opt.Request().GetDbName())
	return f.createErr
}

func (f *fakeMilvusDatabaseClient) UseDatabase(_ context.Context, opt milvusclient.UseDatabaseOption) error {
	f.used = append(f.used, opt.DbName())
	return f.useErr
}

func TestEnsureClientDatabaseSkipsDefault(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{}

	if err := ensureClientDatabase(context.Background(), fake, "default"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 0 || len(fake.used) != 0 {
		t.Fatalf("unexpected calls: created=%v used=%v", fake.created, fake.used)
	}
}

func TestEnsureClientDatabaseCreatesMissingDatabase(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{listDatabases: []string{"default"}}

	if err := ensureClientDatabase(context.Background(), fake, "nova"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 1 || fake.created[0] != "nova" {
		t.Fatalf("unexpected create calls: %v", fake.created)
	}
	if len(fake.used) != 1 || fake.used[0] != "nova" {
		t.Fatalf("unexpected use calls: %v", fake.used)
	}
}

func TestEnsureClientDatabaseUsesExistingDatabase(t *testing.T) {
	fake := &fakeMilvusDatabaseClient{listDatabases: []string{"default", "nova"}}

	if err := ensureClientDatabase(context.Background(), fake, "nova"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.created) != 0 {
		t.Fatalf("unexpected create calls: %v", fake.created)
	}
	if len(fake.used) != 1 || fake.used[0] != "nova" {
		t.Fatalf("unexpected use calls: %v", fake.used)
	}
}

func TestIsMilvusDatabaseAlreadyExistsError(t *testing.T) {
	if !isMilvusDatabaseAlreadyExistsError(errors.New("database already exists")) {
		t.Fatal("expected already-exists error to match")
	}
	if isMilvusDatabaseAlreadyExistsError(errors.New("database not found")) {
		t.Fatal("unexpected match for not-found error")
	}
}
