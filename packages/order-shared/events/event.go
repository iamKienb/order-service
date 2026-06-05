package events

const (
	TopicOrderCreated   = "order-service.order.created"
	TopicOrderConfirmed = "order-service.order.confirmed"
	TopicOrderCancelled = "order-service.order.cancelled"
)

var Topics = []string{
	TopicOrderCreated,
	TopicOrderConfirmed,
	TopicOrderCancelled,
}
