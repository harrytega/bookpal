package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"test-project/internal/models"
	"test-project/internal/test"
	"test-project/internal/util/db"
)

func TestOr(t *testing.T) {
	age := 42
	filter := UserFilter{
		Name: Name{
			PublicName: PublicName{
				First: "Max",
			},
			MiddleName: "Gustav",
			Lastname:   "Muster",
		},
		Country: "Austria",
		City:    "Vienna",
		Scopes:  []string{"app", "user_info"},
		Age:     &age,
	}

	qms := []qm.QueryMod{
		qm.Where("id = ?", 123),
		qm.Where("username = ?", "max.muster@example.org"),
		db.WhereJSON("users", "profile", filter),
	}
	sql, args := buildOrQuery(t, qms)

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}

func TestOrSingle(t *testing.T) {
	q := qm.Where("username = ?", "max.muster@example.org")
	qms := db.CombineWithOr([]qm.QueryMod{q})
	require.Len(t, qms, 1)
	assert.Equal(t, q, qms[0])
}

func TestOrEmpty(t *testing.T) {
	qms := db.CombineWithOr([]qm.QueryMod{})
	assert.Empty(t, qms)

	qms = db.CombineWithOr(nil)
	assert.Empty(t, qms)
}

func buildOrQuery(t *testing.T, qms []qm.QueryMod) (string, []interface{}) {
	t.Helper()

	o := db.CombineWithOr(qms)
	require.Greater(t, len(o), 0)

	o = append(o, qm.Select("*"), qm.From("users"))
	q := models.NewQuery(o...)

	return queries.BuildQuery(q)
}
