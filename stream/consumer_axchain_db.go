// (c) 2021, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/axiacoin/axia-network-v2/utils/hashing"
	"github.com/axiacoin/axia-network-v2-coreth/core/types"
	"github.com/axiacoin/axia-network-v2-magellan/cfg"
	"github.com/axiacoin/axia-network-v2-magellan/db"
	"github.com/axiacoin/axia-network-v2-magellan/modelsc"
	"github.com/axiacoin/axia-network-v2-magellan/services"
	"github.com/axiacoin/axia-network-v2-magellan/services/indexes/cvm"
	"github.com/axiacoin/axia-network-v2-magellan/servicesctrl"
	"github.com/axiacoin/axia-network-v2-magellan/utils"
)

type consumerAXChainDB struct {
	id string
	sc *servicesctrl.Control

	// metrics
	metricProcessedCountKey       string
	metricProcessMillisCounterKey string
	metricSuccessCountKey         string
	metricFailureCountKey         string

	conf cfg.Config

	// Concurrency control
	quitCh   chan struct{}
	consumer *cvm.Writer

	topicName     string
	topicTrcName  string
	topicLogsName string
}

func NewConsumerAXChainDB() ProcessorFactoryInstDB {
	return func(sc *servicesctrl.Control, conf cfg.Config) (ProcessorDB, error) {
		c := &consumerAXChainDB{
			conf:                          conf,
			sc:                            sc,
			metricProcessedCountKey:       fmt.Sprintf("consume_records_processed_%s_axchain", conf.AXchainID),
			metricProcessMillisCounterKey: fmt.Sprintf("consume_records_process_millis_%s_axchain", conf.AXchainID),
			metricSuccessCountKey:         fmt.Sprintf("consume_records_success_%s_axchain", conf.AXchainID),
			metricFailureCountKey:         fmt.Sprintf("consume_records_failure_%s_axchain", conf.AXchainID),
			id:                            fmt.Sprintf("consumer %d %s axchain", conf.NetworkID, conf.AXchainID),

			quitCh: make(chan struct{}),
		}
		utils.Prometheus.CounterInit(c.metricProcessedCountKey, "records processed")
		utils.Prometheus.CounterInit(c.metricProcessMillisCounterKey, "records processed millis")
		utils.Prometheus.CounterInit(c.metricSuccessCountKey, "records success")
		utils.Prometheus.CounterInit(c.metricFailureCountKey, "records failure")
		sc.InitConsumeMetrics()

		var err error
		c.consumer, err = cvm.NewWriter(c.conf.NetworkID, c.conf.AXchainID)
		if err != nil {
			_ = c.Close()
			return nil, err
		}

		c.topicName = fmt.Sprintf("%d-%s-axchain", c.conf.NetworkID, c.conf.AXchainID)
		c.topicTrcName = fmt.Sprintf("%d-%s-axchain-trc", c.conf.NetworkID, c.conf.AXchainID)
		c.topicLogsName = fmt.Sprintf("%d-%s-axchain-logs", conf.NetworkID, conf.AXchainID)

		return c, nil
	}
}

// Close shuts down the producer
func (c *consumerAXChainDB) Close() error {
	return nil
}

func (c *consumerAXChainDB) ID() string {
	return c.id
}

func (c *consumerAXChainDB) Topic() []string {
	return []string{c.topicName, c.topicTrcName, c.topicLogsName}
}

func (c *consumerAXChainDB) Process(conns *utils.Connections, row *db.TxPool) error {
	switch row.Topic {
	case c.topicName:
		msg := &Message{
			id:         row.MsgKey,
			chainID:    c.conf.AXchainID,
			body:       row.Serialization,
			timestamp:  row.CreatedAt.UTC().Unix(),
			nanosecond: int64(row.CreatedAt.UTC().Nanosecond()),
		}
		return c.Consume(conns, msg)
	case c.topicTrcName:
		msg := &Message{
			id:         row.MsgKey,
			chainID:    c.conf.AXchainID,
			body:       row.Serialization,
			timestamp:  row.CreatedAt.UTC().Unix(),
			nanosecond: int64(row.CreatedAt.UTC().Nanosecond()),
		}
		return c.ConsumeTrace(conns, msg)
	case c.topicLogsName:
		msg := &Message{
			id:         row.MsgKey,
			chainID:    c.conf.AXchainID,
			body:       row.Serialization,
			timestamp:  row.CreatedAt.UTC().Unix(),
			nanosecond: int64(row.CreatedAt.UTC().Nanosecond()),
		}
		return c.ConsumeLogs(conns, msg)
	}

	return nil
}

func (c *consumerAXChainDB) ConsumeLogs(conns *utils.Connections, msg services.Consumable) error {
	txLogs := &types.Log{}
	err := json.Unmarshal(msg.Body(), txLogs)
	if err != nil {
		return err
	}
	collectors := utils.NewCollectors(
		utils.NewCounterIncCollect(c.metricProcessedCountKey),
		utils.NewCounterObserveMillisCollect(c.metricProcessMillisCounterKey),
		utils.NewCounterIncCollect(servicesctrl.MetricConsumeProcessedCountKey),
		utils.NewCounterObserveMillisCollect(servicesctrl.MetricConsumeProcessMillisCounterKey),
	)
	defer func() {
		err := collectors.Collect()
		if err != nil {
			c.sc.Log.Error("collectors.Collect: %s", err)
		}
	}()

	id := hashing.ComputeHash256(msg.Body())

	nmsg := NewMessage(string(id), msg.ChainID(), msg.Body(), msg.Timestamp(), msg.Nanosecond())

	rsleep := utils.NewRetrySleeper(1, 100*time.Millisecond, time.Second)
	for {
		err = c.persistConsumeLogs(conns, nmsg, txLogs)
		if !utils.ErrIsLockError(err) {
			break
		}
		rsleep.Inc()
	}

	if err != nil {
		c.Failure()
		collectors.Error()
		c.sc.Log.Error("consumer.Consume: %s", err)
		return err
	}
	c.Success()

	return nil
}

func (c *consumerAXChainDB) ConsumeTrace(conns *utils.Connections, msg services.Consumable) error {
	transactionTrace := &modelsc.TransactionTrace{}
	err := json.Unmarshal(msg.Body(), transactionTrace)
	if err != nil {
		return err
	}
	collectors := utils.NewCollectors(
		utils.NewCounterIncCollect(c.metricProcessedCountKey),
		utils.NewCounterObserveMillisCollect(c.metricProcessMillisCounterKey),
		utils.NewCounterIncCollect(servicesctrl.MetricConsumeProcessedCountKey),
		utils.NewCounterObserveMillisCollect(servicesctrl.MetricConsumeProcessMillisCounterKey),
	)
	defer func() {
		err := collectors.Collect()
		if err != nil {
			c.sc.Log.Error("collectors.Collect: %s", err)
		}
	}()

	id := hashing.ComputeHash256(transactionTrace.Trace)

	nmsg := NewMessage(string(id), msg.ChainID(), transactionTrace.Trace, msg.Timestamp(), msg.Nanosecond())

	rsleep := utils.NewRetrySleeper(1, 100*time.Millisecond, time.Second)
	for {
		err = c.persistConsumeTrace(conns, nmsg, transactionTrace)
		if !utils.ErrIsLockError(err) {
			break
		}
		rsleep.Inc()
	}

	if err != nil {
		c.Failure()
		collectors.Error()
		c.sc.Log.Error("consumer.Consume: %s", err)
		return err
	}
	c.Success()

	return nil
}

func (c *consumerAXChainDB) Consume(conns *utils.Connections, msg services.Consumable) error {
	block, err := modelsc.Unmarshal(msg.Body())
	if err != nil {
		return err
	}

	collectors := utils.NewCollectors(
		utils.NewCounterIncCollect(c.metricProcessedCountKey),
		utils.NewCounterObserveMillisCollect(c.metricProcessMillisCounterKey),
		utils.NewCounterIncCollect(servicesctrl.MetricConsumeProcessedCountKey),
		utils.NewCounterObserveMillisCollect(servicesctrl.MetricConsumeProcessMillisCounterKey),
	)
	defer func() {
		err := collectors.Collect()
		if err != nil {
			c.sc.Log.Error("collectors.Collect: %s", err)
		}
	}()

	if block.BlockExtraData == nil {
		block.BlockExtraData = []byte("")
	}

	id := hashing.ComputeHash256(block.BlockExtraData)
	nmsg := NewMessage(string(id), msg.ChainID(), block.BlockExtraData, msg.Timestamp(), msg.Nanosecond())

	rsleep := utils.NewRetrySleeper(1, 100*time.Millisecond, time.Second)
	for {
		err = c.persistConsume(conns, nmsg, block)
		if !utils.ErrIsLockError(err) {
			break
		}
		rsleep.Inc()
	}

	if err != nil {
		c.Failure()
		collectors.Error()
		c.sc.Log.Error("consumer.Consume: %s", err)
		return err
	}
	c.Success()

	c.sc.BalanceManager.Exec()

	return nil
}

func (c *consumerAXChainDB) persistConsumeLogs(conns *utils.Connections, msg services.Consumable, txLogs *types.Log) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), cfg.DefaultConsumeProcessWriteTimeout)
	defer cancelFn()
	return c.consumer.ConsumeLogs(ctx, conns, msg, txLogs, c.sc.Persist)
}

func (c *consumerAXChainDB) persistConsumeTrace(conns *utils.Connections, msg services.Consumable, transactionTrace *modelsc.TransactionTrace) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), cfg.DefaultConsumeProcessWriteTimeout)
	defer cancelFn()
	return c.consumer.ConsumeTrace(ctx, conns, msg, transactionTrace, c.sc.Persist)
}

func (c *consumerAXChainDB) persistConsume(conns *utils.Connections, msg services.Consumable, block *modelsc.Block) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), cfg.DefaultConsumeProcessWriteTimeout)
	defer cancelFn()
	return c.consumer.Consume(ctx, conns, msg, block, c.sc.Persist)
}

func (c *consumerAXChainDB) Failure() {
	_ = utils.Prometheus.CounterInc(c.metricFailureCountKey)
	_ = utils.Prometheus.CounterInc(servicesctrl.MetricConsumeFailureCountKey)
}

func (c *consumerAXChainDB) Success() {
	_ = utils.Prometheus.CounterInc(c.metricSuccessCountKey)
	_ = utils.Prometheus.CounterInc(servicesctrl.MetricConsumeSuccessCountKey)
}
