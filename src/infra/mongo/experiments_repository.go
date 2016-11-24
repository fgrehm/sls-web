package mongo

import (
	"github.com/fgrehm/sls-web/src/core/solver"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const EXPERIMENTS_COLLECTION = "experiments"

type experimentsRepository struct {
	db        string
	session   *mgo.Session
	solutions solver.SolutionRepository
}

func (r *experimentsRepository) Store(e *solver.Experiment) error {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(EXPERIMENTS_COLLECTION)

	_, err := c.Upsert(bson.M{"id": e.ID}, bson.M{"$set": e})

	return err
}

func (r *experimentsRepository) Find(id solver.ExperimentID, externalKey *string) (*solver.Experiment, error) {
	sess := r.session.Copy()
	defer sess.Close()

	query := bson.M{"id": id}
	if externalKey != nil {
		query["externalkey"] = *externalKey
	}

	c := sess.DB(r.db).C(EXPERIMENTS_COLLECTION)
	var experiment solver.Experiment
	if err := c.Find(query).One(&experiment); err != nil {
		if err == mgo.ErrNotFound {
			return nil, solver.ErrUnknown
		}
		return nil, err
	}

	// TODO: This is kinda slow, improve it
	if experiment.Solutions != nil {
		for _, es := range *experiment.Solutions {
			// TODO: Handle unknown
			es.Solution, _ = r.solutions.Find(es.ID, nil)
		}
	}

	return &experiment, nil
}

// NewExperimentsRepository returns a new instance of an experiments repository.
func NewExperimentsRepository(session *mgo.Session, db string, solutions solver.SolutionRepository) solver.ExperimentRepository {
	r := &experimentsRepository{db, session, solutions}

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C(EXPERIMENTS_COLLECTION)
	if err := c.EnsureIndex(index); err != nil {
		panic(err)
	}

	return r
}
