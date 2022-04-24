/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dapr

import (
	"strings"

	"github.com/valyala/fasthttp"

	// Included components in compiled daprd.

	// Secret stores.

	"github.com/dapr/kit/logger"

	"github.com/dapr/dapr/pkg/runtime"

	"github.com/dapr/components-contrib/secretstores"
	secretstore_file "github.com/dapr/components-contrib/secretstores/local/file"
	secretstores_loader "github.com/dapr/dapr/pkg/components/secretstores"

	// State Stores.

	"github.com/dapr/components-contrib/state"
	"github.com/dapr/components-contrib/state/postgresql"
	state_redis "github.com/dapr/components-contrib/state/redis"

	state_loader "github.com/dapr/dapr/pkg/components/state"

	// Pub/Sub.
	pubs "github.com/dapr/components-contrib/pubsub"
	pubsub_kafka "github.com/dapr/components-contrib/pubsub/kafka"
	pubsub_redis "github.com/dapr/components-contrib/pubsub/redis"

	pubsub_loader "github.com/dapr/dapr/pkg/components/pubsub"

	// Name resolutions.
	nr "github.com/dapr/components-contrib/nameresolution"
	nr_consul "github.com/dapr/components-contrib/nameresolution/consul"
	nr_kubernetes "github.com/dapr/components-contrib/nameresolution/kubernetes"
	nr_mdns "github.com/dapr/components-contrib/nameresolution/mdns"

	nr_loader "github.com/dapr/dapr/pkg/components/nameresolution"

	// Bindings.
	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/bindings/postgres"
	"github.com/dapr/components-contrib/bindings/twitter"

	bindings_loader "github.com/dapr/dapr/pkg/components/bindings"

	// HTTP Middleware.

	middleware "github.com/dapr/components-contrib/middleware"
	"github.com/dapr/components-contrib/middleware/http/bearer"
	"github.com/dapr/components-contrib/middleware/http/oauth2"
	"github.com/dapr/components-contrib/middleware/http/oauth2clientcredentials"
	"github.com/dapr/components-contrib/middleware/http/opa"
	"github.com/dapr/components-contrib/middleware/http/ratelimit"
	"github.com/dapr/components-contrib/middleware/http/sentinel"

	http_middleware_loader "github.com/dapr/dapr/pkg/components/middleware/http"
	http_middleware "github.com/dapr/dapr/pkg/middleware/http"
)

var (
	logContrib = logger.NewLogger("dapr.contrib")
)

func RegisterComponents(rt *Runtime) {
	rt.AddOptions(
		runtime.WithSecretStores(
			// 	secretstores_loader.New("kubernetes", func() secretstores.SecretStore {
			// 		return sercetstores_kubernetes.NewKubernetesSecretStore(logContrib)
			// 	}),
			// 	secretstores_loader.New("azure.keyvault", func() secretstores.SecretStore {
			// 		return keyvault.NewAzureKeyvaultSecretStore(logContrib)
			// 	}),
			// 	secretstores_loader.New("hashicorp.vault", func() secretstores.SecretStore {
			// 		return vault.NewHashiCorpVaultSecretStore(logContrib)
			// 	}),
			// 	secretstores_loader.New("aws.secretmanager", func() secretstores.SecretStore {
			// 		return secretmanager.NewSecretManager(logContrib)
			// 	}),
			// 	secretstores_loader.New("aws.parameterstore", func() secretstores.SecretStore {
			// 		return parameterstore.NewParameterStore(logContrib)
			// 	}),
			// 	secretstores_loader.New("gcp.secretmanager", func() secretstores.SecretStore {
			// 		return gcp_secretmanager.NewSecreteManager(logContrib)
			// 	}),
			secretstores_loader.New("local.file", func() secretstores.SecretStore {
				return secretstore_file.NewLocalSecretStore(logContrib)
			}),
			// 	secretstores_loader.New("local.env", func() secretstores.SecretStore {
			// 		return secretstore_env.NewEnvSecretStore(logContrib)
			// 	}),
		),
		runtime.WithStates(
			state_loader.New("redis", func() state.Store {
				return state_redis.NewRedisStateStore(logContrib)
			}),
			// state_loader.New("consul", func() state.Store {
			// 	return consul.NewConsulStateStore(logContrib)
			// }),
			// state_loader.New("azure.blobstorage", func() state.Store {
			// 	return state_azure_blobstorage.NewAzureBlobStorageStore(logContrib)
			// }),
			// state_loader.New("azure.cosmosdb", func() state.Store {
			// 	return state_cosmosdb.NewCosmosDBStateStore(logContrib)
			// }),
			// state_loader.New("azure.tablestorage", func() state.Store {
			// 	return state_azure_tablestorage.NewAzureTablesStateStore(logContrib)
			// }),
			// state_loader.New("cassandra", func() state.Store {
			// 	return cassandra.NewCassandraStateStore(logContrib)
			// }),
			// state_loader.New("memcached", func() state.Store {
			// 	return memcached.NewMemCacheStateStore(logContrib)
			// }),
			// state_loader.New("mongodb", func() state.Store {
			// 	return mongodb.NewMongoDB(logContrib)
			// }),
			// state_loader.New("zookeeper", func() state.Store {
			// 	return zookeeper.NewZookeeperStateStore(logContrib)
			// }),
			// state_loader.New("gcp.firestore", func() state.Store {
			// 	return firestore.NewFirestoreStateStore(logContrib)
			// }),
			state_loader.New("postgresql", func() state.Store {
				return postgresql.NewPostgreSQLStateStore(logContrib)
			}),
			// state_loader.New("sqlserver", func() state.Store {
			// 	return sqlserver.NewSQLServerStateStore(logContrib)
			// }),
			// state_loader.New("hazelcast", func() state.Store {
			// 	return hazelcast.NewHazelcastStore(logContrib)
			// }),
			// state_loader.New("cloudstate.crdt", func() state.Store {
			// 	return cloudstate.NewCRDT(logContrib)
			// }),
			// state_loader.New("couchbase", func() state.Store {
			// 	return couchbase.NewCouchbaseStateStore(logContrib)
			// }),
			// state_loader.New("aerospike", func() state.Store {
			// 	return aerospike.NewAerospikeStateStore(logContrib)
			// }),
			// state_loader.New("rethinkdb", func() state.Store {
			// 	return rethinkdb.NewRethinkDBStateStore(logContrib)
			// }),
			// state_loader.New("aws.dynamodb", state_dynamodb.NewDynamoDBStateStore),
			// state_loader.New("mysql", func() state.Store {
			// 	return state_mysql.NewMySQLStateStore(logContrib)
			// }),
		),
		runtime.WithPubSubs(
			// pubsub_loader.New("azure.eventhubs", func() pubs.PubSub {
			// 	return pubsub_eventhubs.NewAzureEventHubs(logContrib)
			// }),
			// pubsub_loader.New("azure.servicebus", func() pubs.PubSub {
			// 	return servicebus.NewAzureServiceBus(logContrib)
			// }),
			// pubsub_loader.New("gcp.pubsub", func() pubs.PubSub {
			// 	return pubsub_gcp.NewGCPPubSub(logContrib)
			// }),
			// pubsub_loader.New("hazelcast", func() pubs.PubSub {
			// 	return pubsub_hazelcast.NewHazelcastPubSub(logContrib)
			// }),
			pubsub_loader.New("kafka", func() pubs.PubSub {
				return pubsub_kafka.NewKafka(logContrib)
			}),
			// pubsub_loader.New("mqtt", func() pubs.PubSub {
			// 	return pubsub_mqtt.NewMQTTPubSub(logContrib)
			// }),
			// pubsub_loader.New("natsstreaming", func() pubs.PubSub {
			// 	return natsstreaming.NewNATSStreamingPubSub(logContrib)
			// }),
			// pubsub_loader.New("pulsar", func() pubs.PubSub {
			// 	return pubsub_pulsar.NewPulsar(logContrib)
			// }),
			// pubsub_loader.New("rabbitmq", func() pubs.PubSub {
			// 	return rabbitmq.NewRabbitMQ(logContrib)
			// }),
			pubsub_loader.New("redis", func() pubs.PubSub {
				return pubsub_redis.NewRedisStreams(logContrib)
			}),
			// pubsub_loader.New("snssqs", func() pubs.PubSub {
			// 	return pubsub_snssqs.NewSnsSqs(logContrib)
			// }),
		),
		runtime.WithNameResolutions(
			nr_loader.New("mdns", func() nr.Resolver {
				return nr_mdns.NewResolver(logContrib)
			}),
			nr_loader.New("kubernetes", func() nr.Resolver {
				return nr_kubernetes.NewResolver(logContrib)
			}),
			nr_loader.New("consul", func() nr.Resolver {
				return nr_consul.NewResolver(logContrib)
			}),
		),
		runtime.WithInputBindings(
			// 	bindings_loader.NewInput("aws.sqs", func() bindings.InputBinding {
			// 		return sqs.NewAWSSQS(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("aws.kinesis", func() bindings.InputBinding {
			// 		return kinesis.NewAWSKinesis(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("azure.eventgrid", func() bindings.InputBinding {
			// 		return eventgrid.NewAzureEventGrid(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("azure.eventhubs", func() bindings.InputBinding {
			// 		return eventhubs.NewAzureEventHubs(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("azure.servicebusqueues", func() bindings.InputBinding {
			// 		return servicebusqueues.NewAzureServiceBusQueues(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("azure.storagequeues", func() bindings.InputBinding {
			// 		return storagequeues.NewAzureStorageQueues(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("cron", func() bindings.InputBinding {
			// 		return cron.NewCron(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("dingtalk.webhook", func() bindings.InputBinding {
			// 		return dingtalk_webhook.NewDingTalkWebhook(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("gcp.pubsub", func() bindings.InputBinding {
			// 		return pubsub.NewGCPPubSub(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("kafka", func() bindings.InputBinding {
			// 		return kafka.NewKafka(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("kubernetes", func() bindings.InputBinding {
			// 		return kubernetes.NewKubernetes(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("mqtt", func() bindings.InputBinding {
			// 		return mqtt.NewMQTT(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("rabbitmq", func() bindings.InputBinding {
			// 		return bindings_rabbitmq.NewRabbitMQ(logContrib)
			// 	}),
			// 	bindings_loader.NewInput("rethinkdb.statechange", func() bindings.InputBinding {
			// 		return statechange.NewRethinkDBStateChangeBinding(logContrib)
			// 	}),
			bindings_loader.NewInput("twitter", func() bindings.InputBinding {
				return twitter.NewTwitter(logContrib)
			}),
			// 	bindings_loader.NewInput("zeebe.jobworker", func() bindings.InputBinding {
			// 		return bindings_zeebe_jobworker.NewZeebeJobWorker(logContrib)
			// 	}),
		),
		runtime.WithOutputBindings(
			// bindings_loader.NewOutput("alicloud.oss", func() bindings.OutputBinding {
			// 	return oss.NewAliCloudOSS(logContrib)
			// }),
			// bindings_loader.NewOutput("apns", func() bindings.OutputBinding {
			// 	return apns.NewAPNS(logContrib)
			// }),
			// bindings_loader.NewOutput("aws.s3", func() bindings.OutputBinding {
			// 	return s3.NewAWSS3(logContrib)
			// }),
			// bindings_loader.NewOutput("aws.sqs", func() bindings.OutputBinding {
			// 	return sqs.NewAWSSQS(logContrib)
			// }),
			// bindings_loader.NewOutput("aws.sns", func() bindings.OutputBinding {
			// 	return sns.NewAWSSNS(logContrib)
			// }),
			// bindings_loader.NewOutput("aws.kinesis", func() bindings.OutputBinding {
			// 	return kinesis.NewAWSKinesis(logContrib)
			// }),
			// bindings_loader.NewOutput("aws.dynamodb", func() bindings.OutputBinding {
			// 	return dynamodb.NewDynamoDB(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.blobstorage", func() bindings.OutputBinding {
			// 	return blobstorage.NewAzureBlobStorage(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.cosmosdb", func() bindings.OutputBinding {
			// 	return bindings_cosmosdb.NewCosmosDB(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.eventgrid", func() bindings.OutputBinding {
			// 	return eventgrid.NewAzureEventGrid(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.eventhubs", func() bindings.OutputBinding {
			// 	return eventhubs.NewAzureEventHubs(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.servicebusqueues", func() bindings.OutputBinding {
			// 	return servicebusqueues.NewAzureServiceBusQueues(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.signalr", func() bindings.OutputBinding {
			// 	return signalr.NewSignalR(logContrib)
			// }),
			// bindings_loader.NewOutput("azure.storagequeues", func() bindings.OutputBinding {
			// 	return storagequeues.NewAzureStorageQueues(logContrib)
			// }),
			// bindings_loader.NewOutput("cron", func() bindings.OutputBinding {
			// 	return cron.NewCron(logContrib)
			// }),
			// bindings_loader.NewOutput("dingtalk.webhook", func() bindings.OutputBinding {
			// 	return dingtalk_webhook.NewDingTalkWebhook(logContrib)
			// }),
			// bindings_loader.NewOutput("gcp.bucket", func() bindings.OutputBinding {
			// 	return bucket.NewGCPStorage(logContrib)
			// }),
			// bindings_loader.NewOutput("gcp.pubsub", func() bindings.OutputBinding {
			// 	return pubsub.NewGCPPubSub(logContrib)
			// }),
			// bindings_loader.NewOutput("http", func() bindings.OutputBinding {
			// 	return http.NewHTTP(logContrib)
			// }),
			// bindings_loader.NewOutput("influx", func() bindings.OutputBinding {
			// 	return influx.NewInflux(logContrib)
			// }),
			// bindings_loader.NewOutput("kafka", func() bindings.OutputBinding {
			// 	return kafka.NewKafka(logContrib)
			// }),
			// bindings_loader.NewOutput("localstorage", func() bindings.OutputBinding {
			// 	return localstorage.NewLocalStorage(logContrib)
			// }),
			// bindings_loader.NewOutput("mqtt", func() bindings.OutputBinding {
			// 	return mqtt.NewMQTT(logContrib)
			// }),
			// bindings_loader.NewOutput("mysql", func() bindings.OutputBinding {
			// 	return mysql.NewMysql(logContrib)
			// }),
			bindings_loader.NewOutput("postgres", func() bindings.OutputBinding {
				return postgres.NewPostgres(logContrib)
			}),
			// bindings_loader.NewOutput("postmark", func() bindings.OutputBinding {
			// 	return postmark.NewPostmark(logContrib)
			// }),
			// bindings_loader.NewOutput("rabbitmq", func() bindings.OutputBinding {
			// 	return bindings_rabbitmq.NewRabbitMQ(logContrib)
			// }),
			// bindings_loader.NewOutput("redis", func() bindings.OutputBinding {
			// 	return redis.NewRedis(logContrib)
			// }),
			// bindings_loader.NewOutput("smtp", func() bindings.OutputBinding {
			// 	return smtp.NewSMTP(logContrib)
			// }),
			// bindings_loader.NewOutput("twilio.sms", func() bindings.OutputBinding {
			// 	return sms.NewSMS(logContrib)
			// }),
			// bindings_loader.NewOutput("twilio.sendgrid", func() bindings.OutputBinding {
			// 	return sendgrid.NewSendGrid(logContrib)
			// }),
			// bindings_loader.NewOutput("twitter", func() bindings.OutputBinding {
			// 	return twitter.NewTwitter(logContrib)
			// }),
			// bindings_loader.NewOutput("zeebe.command", func() bindings.OutputBinding {
			// 	return bindings_zeebe_command.NewZeebeCommand(logContrib)
			// }),
			// bindings_loader.NewOutput("graphql", func() bindings.OutputBinding {
			// 	return graphql.NewGraphQL(logContrib)
			// }),
		),
		runtime.WithHTTPMiddleware(
			http_middleware_loader.New("uppercase", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
					return func(ctx *fasthttp.RequestCtx) {
						body := string(ctx.PostBody())
						ctx.Request.SetBody([]byte(strings.ToUpper(body)))
						h(ctx)
					}
				}, nil
			}),
			http_middleware_loader.New("oauth2", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return oauth2.NewOAuth2Middleware().GetHandler(metadata)
			}),
			http_middleware_loader.New("oauth2clientcredentials", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return oauth2clientcredentials.NewOAuth2ClientCredentialsMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("ratelimit", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return ratelimit.NewRateLimitMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("bearer", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return bearer.NewBearerMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("opa", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return opa.NewMiddleware(log).GetHandler(metadata)
			}),
			http_middleware_loader.New("sentinel", func(metadata middleware.Metadata) (http_middleware.Middleware, error) {
				return sentinel.NewMiddleware(log).GetHandler(metadata)
			}),
		),
	)
}
