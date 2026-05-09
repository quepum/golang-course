package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	platformKafka "repo-stat/platform/kafka"

	"github.com/segmentio/kafka-go"
)

type Adapter struct {
	log          *slog.Logger
	reader       *kafka.Reader
	writer       *kafka.Writer
	topicRequest string
	topicResult  string
}

func NewAdapter(log *slog.Logger, broker, topicReq, topicRes, groupId string) *Adapter {
	return &Adapter{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topicReq,
			GroupID: groupId,
		}),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
		topicRequest: topicReq,
		topicResult:  topicRes,
	}
}

func (a *Adapter) StartListening(ctx context.Context, handler func(req platformKafka.FetchRequest)) error {
	for {
		msg, err := a.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			a.log.Error("failed to read message from kafka", "error", err)
			continue
		}

		var req platformKafka.FetchRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			a.log.Error("failed to unmarshal request", "error", err)
			continue
		}
		handler(req)
	}
}

func (a *Adapter) SendFetchResult(ctx context.Context, res platformKafka.FetchResult) error {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return a.writer.WriteMessages(ctx, kafka.Message{
		Topic: a.topicResult,
		Key:   []byte(res.Owner + "/" + res.Repo),
		Value: data,
	})
}

func (a *Adapter) SendFetchRequest(ctx context.Context, req platformKafka.FetchRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return a.writer.WriteMessages(ctx, kafka.Message{
		Topic: a.topicRequest,
		Value: data,
	})
}

func (a *Adapter) Close() error {
	_ = a.reader.Close()
	return a.writer.Close()
}
