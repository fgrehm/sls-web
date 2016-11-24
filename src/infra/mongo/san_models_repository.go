package mongo

import (
	"github.com/fgrehm/sls-web/src/core/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SAN_MODELS_COLLECTION = "sanModels"

type sanModelRepository struct {
	db      string
	session *mgo.Session
}

func (r *sanModelRepository) Store(m *models.SanModel) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SAN_MODELS_COLLECTION)

	_, err := c.Upsert(bson.M{"id": m.ID}, bson.M{"$set": m})

	return err
}

func (r *sanModelRepository) Find(id models.SanModelID) (*models.SanModel, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SAN_MODELS_COLLECTION)

	var result models.SanModel
	if err := c.Find(bson.M{"id": id}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, models.ErrUnknown
		}
		return nil, err
	}

	return &result, nil
}

func (r *sanModelRepository) All() ([]*models.SanModel, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SAN_MODELS_COLLECTION)

	var result []*models.SanModel
	if err := c.Find(bson.M{}).All(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// NewSanModelsRepository returns a new instance of the models repository.
func NewSanModelsRepository(session *mgo.Session, db string) models.Repository {
	r := &sanModelRepository{db, session}

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SAN_MODELS_COLLECTION)
	if err := c.EnsureIndex(index); err != nil {
		panic(err)
	}

	return r
}
