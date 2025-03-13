package db_test

import (
	"testing"

	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"test-project/internal/models"
	"test-project/internal/test"
	"test-project/internal/util/db"
)

func TestWhereIn(t *testing.T) {
	q := models.NewQuery(
		qm.Select("*"),
		qm.From("users"),
		db.InnerJoin("users", "id", "app_user_profiles", "user_id"),
		db.WhereIn("app_user_profiles", "username", []string{"max", "muster", "peter"}),
	)

	sql, args := queries.BuildQuery(q)

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}

func TestNIN(t *testing.T) {
	q := models.NewQuery(
		qm.Select("*"),
		qm.From("users"),
		db.InnerJoin("users", "id", "app_user_profiles", "user_id"),
		db.NIN("app_user_profiles.username", []string{"max", "muster", "peter"}),
	)

	sql, args := queries.BuildQuery(q)

	test.Snapshoter.Label("SQL").Save(t, sql)
	test.Snapshoter.Label("Args").Save(t, args)
}
