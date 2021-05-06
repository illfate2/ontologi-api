package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/illfate2/ontology-api/internal/config"
	"github.com/illfate2/ontology-api/internal/ontology"
)

func TestClassRepo(t *testing.T) {
	conn, closeF := config.MustGetDgraphConn(context.TODO(), config.DefaultDgraphURI)
	defer closeF()
	classRepo := NewClass(conn)
	err := MigrateSchema(context.Background(), conn)
	require.NoError(t, err)

	id, err := classRepo.AddSubclassToParent(context.Background(), ontology.Class{
		ID:     "0x4e21",
		Name:   "Fruit",
		IsRoot: true,
	}, ontology.Class{
		Name:   "SubFruit3",
		DType:  []string{ontology.ClassDType},
		IsRoot: false,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	all, err := classRepo.FindAll(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, all)
}

func TestDelete(t *testing.T) {
	conn, closeF := config.MustGetDgraphConn(context.TODO(), config.DefaultDgraphURI)
	defer closeF()
	classRepo := NewClass(conn)
	err := MigrateSchema(context.Background(), conn)
	require.NoError(t, err)
	t.Log(classRepo.Remove(context.Background(), "0xea8e"))
}
