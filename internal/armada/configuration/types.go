package configuration

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/G-Research/armada/internal/armada/authorization/permissions"
	"github.com/G-Research/armada/internal/common"
)

type UserInfo struct {
	Password string
	Groups   []string
}

type ArmadaConfig struct {
	AnonymousAuth bool

	GrpcPort               uint16
	HttpPort               uint16
	MetricsPort            uint16
	PriorityHalfTime       time.Duration
	Redis                  redis.UniversalOptions
	EventsKafka            KafkaConfig
	EventsRedis            redis.UniversalOptions
	BasicAuth              BasicAuthenticationConfig
	OpenIdAuth             OpenIdAuthenticationConfig
	Kerberos               KerberosAuthenticationConfig
	PermissionGroupMapping map[permissions.Permission][]string
	PermissionScopeMapping map[permissions.Permission][]string

	Scheduling     SchedulingConfig
	EventRetention EventRetentionPolicy
}

type OpenIdAuthenticationConfig struct {
	ProviderUrl string
	GroupsClaim string
}

type BasicAuthenticationConfig struct {
	Users map[string]UserInfo
}

type KerberosAuthenticationConfig struct {
	KeytabLocation string
	PrincipalName  string
	UserNameSuffix string
}

type SchedulingConfig struct {
	UseProbabilisticSchedulingForAllResources bool
	QueueLeaseBatchSize                       uint
	MinimumResourceToSchedule                 common.ComputeResourcesFloat
	MaximalClusterFractionToSchedule          map[string]float64
	Lease                                     LeaseSettings
}

type EventRetentionPolicy struct {
	ExpiryEnabled     bool
	RetentionDuration time.Duration
}

type LeaseSettings struct {
	ExpireAfter        time.Duration
	ExpiryLoopInterval time.Duration
}

type KafkaConfig struct {
	Brokers         []string
	Topic           string
	ConsumerGroupID string
}
