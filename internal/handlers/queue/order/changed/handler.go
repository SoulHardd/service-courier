package changed

import (
	"encoding/json"
	"github.com/SoulHardd/service-courier/internal/domain"
	"github.com/SoulHardd/service-courier/internal/handlers/queue/order/changed/dto"
	"log"

	"github.com/IBM/sarama"
)

type OrderController struct {
	useCase OrderUseCase
}

func New(uc OrderUseCase) *OrderController {
	return &OrderController{useCase: uc}
}

func (c *OrderController) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *OrderController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *OrderController) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for dtoMsg := range claim.Messages() {
		ctx := session.Context()
		log.Printf("order.changed handler: received message: key=%s, value=%s, partition=%d, offset=%d\n",
			string(dtoMsg.Key),
			string(dtoMsg.Value),
			dtoMsg.Partition,
			dtoMsg.Offset,
		)

		var msg dto.Message
		err := json.Unmarshal(dtoMsg.Value, &msg)
		if err != nil {
			log.Printf("order.changed handler: received bad message: %v", err)
			session.MarkMessage(dtoMsg, "")
			continue
		}
		order := domain.Order{
			ID:        msg.OrderId,
			Status:    domain.OrderStatus(msg.Status),
			CreatedAt: msg.CreatedAt,
		}

		if err := c.useCase.Process(ctx, order); err != nil {
			log.Printf("order.changed handler: failed process order %s: %v", order.ID, err)
		}

		session.MarkMessage(dtoMsg, "")
	}
	return nil
}
