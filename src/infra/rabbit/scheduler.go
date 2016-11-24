package rabbit

import (
	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/streadway/amqp"
)

type amqpScheduler struct {
	channel          *amqp.Channel
	experimentsQueue amqp.Queue
	solutionsQueue   amqp.Queue
}

func (s *amqpScheduler) ScheduleSolution(solution *solver.Solution) error {
	return s.channel.Publish(
		"", // exchange
		s.solutionsQueue.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(solution.ID),
		})
}

func (s *amqpScheduler) ScheduleExperiment(experiment *solver.Experiment) error {
	return s.channel.Publish(
		"", // exchange
		s.experimentsQueue.Name, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(experiment.ID),
		})
}

func NewScheduler(ch *amqp.Channel, eq, sq amqp.Queue) solver.Scheduler {
	return &amqpScheduler{ch, eq, sq}
}
