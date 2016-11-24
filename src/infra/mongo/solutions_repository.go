package mongo

import (
	"github.com/fgrehm/sls-web/src/core/solver"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SOLUTIONS_COLLECTION = "solutions"

type solutionsRepository struct {
	db      string
	session *mgo.Session
}

func (r *solutionsRepository) Store(m *solver.Solution) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SOLUTIONS_COLLECTION)

	_, err := c.Upsert(bson.M{"id": m.ID}, bson.M{"$set": m})

	return err
}

func (r *solutionsRepository) Find(id solver.SolutionID, externalKey *string) (*solver.Solution, error) {
	sess := r.session.Copy()
	defer sess.Close()

	query := bson.M{"id": id}
	if externalKey != nil {
		query["externalkey"] = *externalKey
	}

	c := sess.DB(r.db).C(SOLUTIONS_COLLECTION)

	var solution solver.Solution
	if err := c.Find(query).One(&solution); err != nil {
		if err == mgo.ErrNotFound {
			return nil, solver.ErrUnknown
		}
		return nil, err
	}

	return &solution, nil
}

func (r *solutionsRepository) FindForReuse(externalKey, modelHash string, maxIterations uint, tolerance float64) (*solver.Solution, error) {
	sess := r.session.Copy()
	defer sess.Close()

	query := bson.M{
		"externalkey": externalKey,
		"modelhash":   modelHash,
		"tolerance":   tolerance,
		"errored":     false,
		"$or": []bson.M{
			bson.M{
				"steps": bson.M{"$lte": maxIterations},
				"found": true,
			},
			bson.M{
				"maxiterations": maxIterations,
				"found":         false,
			},
		},
	}

	c := sess.DB(r.db).C(SOLUTIONS_COLLECTION)

	var solution solver.Solution
	if err := c.Find(query).One(&solution); err != nil {
		if err == mgo.ErrNotFound {
			return nil, solver.ErrUnknown
		}
		return nil, err
	}

	return &solution, nil
}

// NewSolutionsRepository returns a new instance of a solutions repository.
func NewSolutionsRepository(session *mgo.Session, db string) solver.SolutionRepository {
	r := &solutionsRepository{db, session}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(SOLUTIONS_COLLECTION)

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	if err := c.EnsureIndex(index); err != nil {
		panic(err)
	}

	index = mgo.Index{
		Key:        []string{"externalkey", "modelhash"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	if err := c.EnsureIndex(index); err != nil {
		panic(err)
	}

	return r
}
