package asapp

// RdsService defines the protocol to use redis
type RdsService interface {
	Enqueue(queue string, env PublicEnvelope) CompoundError
	Dequeue(queue string) (PublicEnvelope, CompoundError)

	// QM: queue manager
	AddToQM(key string, queue string) CompoundError
	QMMembers(key string) ([]string, CompoundError)
	RemoveFromQM(key string, queue string) CompoundError
}
