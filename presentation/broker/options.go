package broker

import (
	"github.com/kingstonduy/mcs-money-transfer/domain/account"
	"github.com/lengocson131002/go-clean-core/transport/broker"
	"golang.org/x/sync/errgroup"
)

func (s *BrokerServer) GetStartOptions() []BrokerServerStartOption {
	return []BrokerServerStartOption{
		WithHandlerSubscriptionRoutes(),
	}
}

type BrokerServerStartOption func(*BrokerServer) error

func WithHandlerSubscriptionRoutes() BrokerServerStartOption {
	return func(b *BrokerServer) error {
		eg := new(errgroup.Group)
		eg.Go(func() error {
			_, e := b.broker.Subscribe(TopicRequestCheckBalance, func(e broker.Event) error {
				return HandleBrokerEvent[*account.CheckBalanceRequest, *account.CheckBalanceResponse](b.broker, e, WithReplyTopic(TopicReplyCheckBalance))
			})
			return e
		})

		return eg.Wait()
	}
}
