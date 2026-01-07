package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

const (
	EnvEVMRPCUrl       = "EVM_RPC_URL"
	EnvMaxBlocksPerRun = "MAX_BLOCKS_PER_RUN"
)

const (
	EnvPGDSN = "PG_DSN"

	EnvJSHost                   = "JS_HOST"
	EnvJSUser                   = "JS_USER"
	EnvJSPassword               = "JS_PASSWORD"
	EnvJSClientName             = "JS_CLIENT_NAME"
	EnvJSVerbose                = "JS_VERBOSE"
	EnvJSAllowReconnect         = "JS_ALLOW_RECONNECT"
	EnvJSRetryOnFailedConnect   = "JS_RETRY_ON_FAILED_CONNECT"
	EnvJSPublishAsyncMaxPending = "JS_PUBLISH_ASYNC_MAX_PENDING"
)

const (
	EnvAPIHTTPAddr = "API_HTTP_ADDR"
)

const (
	EnvWorkerSyncTick   = "WORKER_SYNC_TICK"
	EnvWorkerClaimLimit = "WORKER_CLAIM_LIMIT"
	EnvWorkerLockTTL    = "WORKER_LOCK_TTL"

	EnvWorkerPoolSize      = "WORKER_POOL_SIZE"
	EnvWorkerPoolQueueSize = "WORKER_POOL_QUEUE_SIZE"

	EnvWorkerStreamSubject    = "WORKER_JS_STREAM_SUBJECT"
	EnvWorkerStreamName       = "WORKER_JS_STREAM_NAME"
	EnvWorkerStreamRetention  = "WORKER_JS_STREAM_RETENTION" // limits|interest|workqueue
	EnvWorkerStreamReplicas   = "WORKER_JS_STREAM_REPLICAS"
	EnvWorkerStreamDuplicates = "WORKER_JS_STREAM_DUPLICATES"

	EnvWorkerConsumerSubject       = "WORKER_JS_CONSUMER_SUBJECT"
	EnvWorkerConsumerName          = "WORKER_JS_CONSUMER_NAME"
	EnvWorkerConsumerAckPolicy     = "WORKER_JS_CONSUMER_ACK_POLICY" // explicit|all|none
	EnvWorkerConsumerAckWait       = "WORKER_JS_CONSUMER_ACK_WAIT"
	EnvWorkerConsumerMaxAckPending = "WORKER_JS_CONSUMER_MAX_ACK_PENDING"
)

type Config struct {
	Common CommonConfig
	API    APIConfig
	Worker WorkerConfig
}

type CommonConfig struct {
	Postgres  PostgresParams
	JetStream JetStreamConfig
}

type APIConfig struct {
	HTTP HTTPConfig
}

type WorkerConfig struct {
	Pool     PoolConfig
	Stream   StreamConfig
	Consumer ConsumerConfig
	Chain    ChainConfig

	SyncTick   time.Duration
	ClaimLimit int
	LockTTL    time.Duration
}

type JetStreamConfig struct {
	Host       string
	User       string
	Password   string
	ClientName string

	Verbose                bool
	AllowReconnect         bool
	RetryOnFailedConnect   bool
	PublishAsyncMaxPending int
}

type PoolConfig struct {
	Workers   int
	QueueSize int
}

type ChainConfig struct {
	EVM EVMConfig
}

type EVMConfig struct {
	RPCURL          string
	MaxBlocksPerRun int
}

type StreamConfig struct {
	Subject         string
	StreamName      string
	RetentionPolicy jetstream.RetentionPolicy
	Replicas        int
	Duplicates      time.Duration
}

type ConsumerConfig struct {
	Subject       string
	Name          string
	AckPolicy     jetstream.AckPolicy
	AckWait       time.Duration
	MaxAckPending int
}

type PostgresParams struct {
	DSN string
}

type HTTPConfig struct {
	Addr string
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		Common: CommonConfig{
			Postgres: PostgresParams{
				DSN: getEnv(EnvPGDSN, ""),
			},
			JetStream: JetStreamConfig{
				Host:                   getEnv(EnvJSHost, "nats://localhost:4222"),
				User:                   getEnv(EnvJSUser, ""),
				Password:               getEnv(EnvJSPassword, ""),
				ClientName:             getEnv(EnvJSClientName, "app"),
				Verbose:                getEnvBool(EnvJSVerbose, false),
				AllowReconnect:         getEnvBool(EnvJSAllowReconnect, true),
				RetryOnFailedConnect:   getEnvBool(EnvJSRetryOnFailedConnect, true),
				PublishAsyncMaxPending: getEnvInt(EnvJSPublishAsyncMaxPending, 1024),
			},
		},

		API: APIConfig{
			HTTP: HTTPConfig{
				Addr: getEnv(EnvAPIHTTPAddr, ":8080"),
			},
		},

		Worker: WorkerConfig{
			SyncTick:   getEnvDuration(EnvWorkerSyncTick, 10*time.Second),
			ClaimLimit: getEnvInt(EnvWorkerClaimLimit, 100),
			LockTTL:    getEnvDuration(EnvWorkerLockTTL, 2*time.Minute),

			Pool: PoolConfig{
				Workers:   getEnvInt(EnvWorkerPoolSize, runtime.NumCPU()),
				QueueSize: getEnvInt(EnvWorkerPoolQueueSize, runtime.NumCPU()),
			},

			Chain: ChainConfig{
				EVM: EVMConfig{
					RPCURL:          getEnv(EnvEVMRPCUrl, ""),
					MaxBlocksPerRun: getEnvInt(EnvMaxBlocksPerRun, 1000),
				},
			},

			Stream: StreamConfig{
				Subject:    getEnv(EnvWorkerStreamSubject, ""),
				StreamName: getEnv(EnvWorkerStreamName, ""),
				Replicas:   getEnvInt(EnvWorkerStreamReplicas, 1),
				Duplicates: getEnvDuration(EnvWorkerStreamDuplicates, 2*time.Minute),
			},
			Consumer: ConsumerConfig{
				Subject:       getEnv(EnvWorkerConsumerSubject, ""),
				Name:          getEnv(EnvWorkerConsumerName, ""),
				AckWait:       getEnvDuration(EnvWorkerConsumerAckWait, 30*time.Second),
				MaxAckPending: getEnvInt(EnvWorkerConsumerMaxAckPending, 1024),
			},
		},
	}

	var err error
	cfg.Worker.Stream.RetentionPolicy, err = parseRetentionPolicy(getEnv(EnvWorkerStreamRetention, "workqueue"))
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", EnvWorkerStreamRetention, err)
	}

	cfg.Worker.Consumer.AckPolicy, err = parseAckPolicy(getEnv(EnvWorkerConsumerAckPolicy, "explicit"))
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", EnvWorkerConsumerAckPolicy, err)
	}

	// Default consumer subject to stream subject if not set
	if strings.TrimSpace(cfg.Worker.Consumer.Subject) == "" && strings.TrimSpace(cfg.Worker.Stream.Subject) != "" {
		cfg.Worker.Consumer.Subject = cfg.Worker.Stream.Subject
	}

	return cfg, nil
}

func (c Config) ValidateAPI() error {
	if strings.TrimSpace(c.Common.Postgres.DSN) == "" {
		return fmt.Errorf("%s is required", EnvPGDSN)
	}
	if strings.TrimSpace(c.API.HTTP.Addr) == "" {
		return fmt.Errorf("%s is required", EnvAPIHTTPAddr)
	}
	return nil
}

func (c Config) ValidateWorker() error {
	if strings.TrimSpace(c.Common.Postgres.DSN) == "" {
		return fmt.Errorf("%s is required", EnvPGDSN)
	}
	if strings.TrimSpace(c.Common.JetStream.Host) == "" {
		return fmt.Errorf("%s is required", EnvJSHost)
	}

	if strings.TrimSpace(c.Worker.Stream.StreamName) == "" {
		return fmt.Errorf("%s is required", EnvWorkerStreamName)
	}
	if strings.TrimSpace(c.Worker.Stream.Subject) == "" {
		return fmt.Errorf("%s is required", EnvWorkerStreamSubject)
	}
	if strings.TrimSpace(c.Worker.Consumer.Name) == "" {
		return fmt.Errorf("%s is required", EnvWorkerConsumerName)
	}

	if c.Worker.SyncTick <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerSyncTick)
	}
	if c.Worker.ClaimLimit <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerClaimLimit)
	}
	if c.Worker.LockTTL <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerLockTTL)
	}

	if c.Worker.Consumer.MaxAckPending <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerConsumerMaxAckPending)
	}
	if c.Worker.Consumer.AckWait <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerConsumerAckWait)
	}
	if c.Worker.Stream.Replicas <= 0 {
		return fmt.Errorf("%s must be > 0", EnvWorkerStreamReplicas)
	}

	if c.Worker.Chain.EVM.RPCURL == "" {
		return fmt.Errorf("%s is required", EnvEVMRPCUrl)
	}
	if c.Worker.Chain.EVM.MaxBlocksPerRun < 1 {
		return fmt.Errorf("%s must be > 0", EnvMaxBlocksPerRun)
	}

	return nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(v)
	}
	return def
}

func getEnvInt(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getEnvBool(key string, def bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func parseRetentionPolicy(s string) (jetstream.RetentionPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "limits", "limit":
		return jetstream.LimitsPolicy, nil
	case "interest":
		return jetstream.InterestPolicy, nil
	case "workqueue", "work_queue", "work-queue":
		return jetstream.WorkQueuePolicy, nil
	default:
		return 0, fmt.Errorf("unknown retention policy %q (use: limits|interest|workqueue)", s)
	}
}

func parseAckPolicy(s string) (jetstream.AckPolicy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "explicit":
		return jetstream.AckExplicitPolicy, nil
	case "all":
		return jetstream.AckAllPolicy, nil
	case "none":
		return jetstream.AckNonePolicy, nil
	default:
		return 0, fmt.Errorf("unknown ack policy %q (use: explicit|all|none)", s)
	}
}
