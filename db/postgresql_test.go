package db

import (
	"context"
	"testing"

	"github.com/dokidokikoi/go-common/db/base"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/lib/pq"
)

type Game struct {
	ID   uint           `grom:"primaryKey"`
	Name string         `gorm:"type:varchar(1024)"`
	Tags pq.StringArray `gorm:"type:text[]"`
}

var gameStore base.BasicCURD[Game]

func init() {
	db, err := NewPostgresql("postgres", "postgres", WithHost("localhost"), WithPassword("password"))
	if err != nil {
		panic(err)
	}
	gameStore = &base.PgModel[Game]{
		DB: db.Debug(),
	}

	db.AutoMigrate(&Game{})
}

func Test_Create(t *testing.T) {
	err := gameStore.Create(context.TODO(), &Game{
		Name: "gta 6", Tags: []string{"犯罪", "都市"},
	}, &meta.CreateOption{})
	if err != nil {
		t.Error(err)
	}

	err = gameStore.Creates(
		context.TODO(), []*Game{
			{Name: "赛博朋克2077", Tags: []string{"赛博朋克", "都市"}},
			{Name: "塞尔达: 荒野之息", Tags: []string{"奇幻", "冒险"}},
		},
		nil)
	if err != nil {
		t.Error(err)
	}
}

func Test_Delete(t *testing.T) {
	err := gameStore.Delete(context.TODO(), &Game{Name: "gta 6"}, nil)
	if err != nil {
		t.Error(err)
	}
}

func Test_Select(t *testing.T) {
	games, err := gameStore.List(context.TODO(), &Game{Name: "gta 6"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 1 {
		t.Fail()
	}

	games, err = gameStore.ListComplex(context.TODO(), &Game{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "Tags",
				Operator: "@>",
				Value:    pq.Array([]string{"都市", "犯罪"}),
			},
		},
	}, nil)
	if err != nil {
		t.Error(err)
	}
	if len(games) != 1 {
		t.Fail()
	}
	t.Logf("%+v", games[0])

	games, err = gameStore.ListComplex(context.TODO(), &Game{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "Tags",
				Operator: "&&",
				Value:    pq.Array([]string{"都市", "犯罪"}),
			},
		},
	}, nil)
	if err != nil {
		t.Error(err)
	}
	if len(games) != 2 {
		t.Fail()
	}
}
