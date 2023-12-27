package worker

const (
	QueuePriorityCritical = "critical"
	QueuePriorityDefault  = "default"
	QueuePriorityLow      = "low"

	// task queue
	DeliveryEmailQueue = "task:email_delivery"
	DeliverySayHello   = "task:say_hello"
)
