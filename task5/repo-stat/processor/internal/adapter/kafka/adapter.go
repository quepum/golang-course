package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	platformKafka "repo-stat/platform/kafka"
	"repo-stat/processor/internal/adapter/storage"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"
)

type Adapter struct {
	log    *slog.Logger
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewAdapter(log *slog.Logger, broker, topicReq, topicRes string) *Adapter {
	return &Adapter{
		log: log,
		writer: &kafka.Writer{
			Addr:        kafka.TCP(broker),
			Topic:       topicReq,
			Balancer:    &kafka.LeastBytes{},
			MaxAttempts: 3,
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{broker},
			Topic:          topicRes,
			GroupID:        "processor-group",
			ReadBackoffMin: 100 * time.Millisecond,
		}),
	}
}

func (a *Adapter) SendFetchRequest(ctx context.Context, req platformKafka.FetchRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal fetch request: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%s/%s", req.Owner, req.Repo)),
		Value: data,
	}

	if err := a.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}
	return nil
}

func (a *Adapter) StartConsumingResults(ctx context.Context, repoStorage *storage.Repo) error {
	a.log.Info("starting kafka result consumer...")

	for {
		msg, err := a.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				a.log.Info("kafka consumer stopping due to context cancellation")
				return nil
			}
			a.log.Error("failed to read message from kafka", "error", err)
			continue
		}

		var res platformKafka.FetchResult
		if err := json.Unmarshal(msg.Value, &res); err != nil {
			a.log.Error("failed to unmarshal result", "error", err)
			continue
		}

		if res.Error != "" {
			a.log.Warn("received error from collector", "owner", res.Owner, "repo", res.Repo, "error", res.Error)
			continue
		}

		var createdAtPg pgtype.Timestamptz
		if res.CreatedAt != "" {
			t, err := time.Parse(time.RFC3339, res.CreatedAt)
			if err != nil {
				a.log.Error("failed to parse CreatedAt", "value", res.CreatedAt, "error", err)
			} else {
				createdAtPg = pgtype.Timestamptz{Time: t, Valid: true}
			}
		}

		err = repoStorage.UpsertInfo(ctx, storage.UpsertInfoParams{
			Owner:       res.Owner,
			Repo:        res.Repo,
			FullName:    pgtype.Text{String: res.FullName, Valid: res.FullName != ""},
			Description: pgtype.Text{String: res.Description, Valid: res.Description != ""},
			Stars:       pgtype.Int4{Int32: int32(res.Stars), Valid: true},
			Forks:       pgtype.Int4{Int32: int32(res.Forks), Valid: true},
			CreatedAt:   createdAtPg,
		})

		if err != nil {
			a.log.Error("failed to upsert repo to db", "owner", res.Owner, "repo", res.Repo, "error", err)
		} else {
			a.log.Debug("repo data saved to cache", "owner", res.Owner, "repo", res.Repo)
		}
	}
}

func (a *Adapter) Close() error {
	var errs []error
	if err := a.writer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close writer: %w", err))
	}
	if err := a.reader.Close(); err != nil {
		errs = append(errs, fmt.Errorf("close reader: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
