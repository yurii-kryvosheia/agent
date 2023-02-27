package kafkatarget

// This code is copied from Promtail. The kafkatarget package is used to
// configure and run the targets that can read kafka entries and forward them
// to other loki components.

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/relabel"

	"github.com/Shopify/sarama"
	"github.com/prometheus/common/model"

	"github.com/grafana/agent/component/common/loki"
	"github.com/grafana/loki/clients/pkg/promtail/targets/target"
	"github.com/grafana/loki/pkg/logproto"
)

type runnableDroppedTarget struct {
	target.Target
	runFn func()
}

func (d *runnableDroppedTarget) run() {
	d.runFn()
}

type KafkaTarget struct {
	discoveredLabels     model.LabelSet
	lbs                  model.LabelSet
	details              ConsumerDetails
	claim                sarama.ConsumerGroupClaim
	session              sarama.ConsumerGroupSession
	client               loki.EntryHandler
	relabelConfig        []*relabel.Config
	useIncomingTimestamp bool
}

func NewKafkaTarget(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
	discoveredLabels, lbs model.LabelSet,
	relabelConfig []*relabel.Config,
	client loki.EntryHandler,
	useIncomingTimestamp bool,
) *KafkaTarget {

	return &KafkaTarget{
		discoveredLabels:     discoveredLabels,
		lbs:                  lbs,
		details:              newDetails(session, claim),
		claim:                claim,
		session:              session,
		client:               client,
		relabelConfig:        relabelConfig,
		useIncomingTimestamp: useIncomingTimestamp,
	}
}

const (
	defaultKafkaMessageKey  = "none"
	labelKeyKafkaMessageKey = "__meta_kafka_message_key"
)

func (t *KafkaTarget) run() {
	defer t.client.Stop()
	for message := range t.claim.Messages() {
		mk := string(message.Key)
		if len(mk) == 0 {
			mk = defaultKafkaMessageKey
		}

		// TODO: Possibly need to format after merging with discovered labels because we can specify multiple labels in source labels
		// https://github.com/grafana/loki/pull/4745#discussion_r750022234
		lbs := format([]labels.Label{{
			Name:  labelKeyKafkaMessageKey,
			Value: mk,
		}}, t.relabelConfig)

		out := t.lbs.Clone()
		if len(lbs) > 0 {
			out = out.Merge(lbs)
		}
		t.client.Chan() <- loki.Entry{
			Entry: logproto.Entry{
				Line:      string(message.Value),
				Timestamp: timestamp(t.useIncomingTimestamp, message.Timestamp),
			},
			Labels: out,
		}
		t.session.MarkMessage(message, "")
	}
}

func timestamp(useIncoming bool, incoming time.Time) time.Time {
	if useIncoming {
		return incoming
	}
	return time.Now()
}

func (t *KafkaTarget) Type() target.TargetType {
	return target.KafkaTargetType
}

func (t *KafkaTarget) Ready() bool {
	return true
}

func (t *KafkaTarget) DiscoveredLabels() model.LabelSet {
	return t.discoveredLabels
}

func (t *KafkaTarget) Labels() model.LabelSet {
	return t.lbs
}

// Details returns target-specific details.
func (t *KafkaTarget) Details() interface{} {
	return t.details
}

type ConsumerDetails struct {

	// MemberID returns the cluster member ID.
	MemberID string

	// GenerationID returns the current generation ID.
	GenerationID int32

	Topic         string
	Partition     int32
	InitialOffset int64
}

func (c ConsumerDetails) String() string {
	return fmt.Sprintf("member_id=%s generation_id=%d topic=%s partition=%d initial_offset=%d", c.MemberID, c.GenerationID, c.Topic, c.Partition, c.InitialOffset)
}

func newDetails(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) ConsumerDetails {
	return ConsumerDetails{
		MemberID:      session.MemberID(),
		GenerationID:  session.GenerationID(),
		Topic:         claim.Topic(),
		Partition:     claim.Partition(),
		InitialOffset: claim.InitialOffset(),
	}
}
