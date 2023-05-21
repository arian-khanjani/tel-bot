package mongodb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tel-bot/model"
	"time"
)

type Repo struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

type ConnProps struct {
	URI  string
	DB   string
	Coll string
}

func New(props ConnProps) (*Repo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(props.URI)
	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := c.Database(props.DB)

	return &Repo{
		client: c,
		db:     db,
		coll:   db.Collection(props.Coll),
	}, nil
}

func (r *Repo) ListClients(ctx context.Context, providerID int64) (*[]model.Client, error) {
	cur, err := r.coll.Find(ctx, bson.D{{"provider_id", providerID}})
	if err != nil {
		return nil, err
	}

	var res []model.Client
	err = cur.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *Repo) GetClient(ctx context.Context, id int64) (*model.Client, error) {
	var res model.Client
	err := r.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *Repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	var res model.User
	err := r.coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}

func (r *Repo) CreateClient(ctx context.Context, u *model.Client) (*model.Client, error) {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repo) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	_, err := r.coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	one, err := r.coll.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	if one.DeletedCount == 0 {
		return errors.New("record was not deleted. maybe ID was incorrect")
	}

	return nil
}

func (r *Repo) Disconnect(ctx context.Context) error {
	err := r.client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CreateIndexes(ctx context.Context, index bson.D) ([]string, error) {
	res, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: index}})
	if err != nil {
		return []string{}, err
	}

	return res, nil
}
