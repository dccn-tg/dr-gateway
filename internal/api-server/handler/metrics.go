package handler

import (
	"net/http"

	"github.com/dccn-tg/dr-gateway/pkg/swagger/server/restapi/operations"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

// This trick of integrating promhttp handler with swagger server is taken from
// the blog: https://www.kaznacheev.me/posts/en/go_swagger_tricks/
type CustomResponder func(http.ResponseWriter, runtime.Producer)

func (c CustomResponder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	c(w, p)
}

func NewCustomResponder(r *http.Request, h http.Handler) middleware.Responder {
	return CustomResponder(func(w http.ResponseWriter, _ runtime.Producer) {
		h.ServeHTTP(w, r)
	})
}

// GetMetrics handles the metrics request with the Prometheus handler
func GetMetrics(ucache *UsersCache, ccache *CollectionsCache) func(p operations.GetMetricsParams) middleware.Responder {

	promRegistry := prometheus.NewRegistry()
	promRegistry.MustRegister(
		collectionCount,
		storageQuota,
		storageUsage,
		userCount,
		userCountOu,
	)

	log.Debugf("GetMetrics called %p", promRegistry)

	return func(p operations.GetMetricsParams) middleware.Responder {
		collectMetrics(ucache, ccache)
		return NewCustomResponder(
			p.HTTPRequest,
			promhttp.HandlerFor(
				promRegistry,
				promhttp.HandlerOpts{
					EnableOpenMetrics: false,
				},
			),
		)
	}
}

// metrics definition
var (
	collectionCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "repository_coll_count",
			Help: "number of collections",
		},
		[]string{
			"projectId",
			"type",
			"state",
			"ou",
		},
	)

	storageQuota = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_storage_quota",
			Help: "storage quota of collections",
		},
		[]string{
			"projectId",
			"type",
			"state",
			"ou",
		},
	)

	storageUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_storage_usage",
			Help: "used storage space of collections",
		},
		[]string{
			"projectId",
			"type",
			"state",
			"ou",
		},
	)

	userCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "repository_user_count",
			Help: "number of users",
		},
		[]string{
			"homeOrg",
		},
	)

	userCountOu = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "repository_ou_user_count",
			Help: "number of OU users",
		},
		[]string{
			"ou",
		},
	)
)

// metrics recording function in an infinite loop
func collectMetrics(ucache *UsersCache, ccache *CollectionsCache) {
	// collections
	collectionCount.Reset()
	storageQuota.Reset()
	storageUsage.Reset()
	for _, coll := range ccache.GetCollections() {

		collectionCount.WithLabelValues(
			coll.ProjectID,
			coll.Type.String(),
			coll.State.String(),
			coll.OrganisationalUnit,
		).Add(1)

		storageQuota.WithLabelValues(
			coll.ProjectID,
			coll.Type.String(),
			coll.State.String(),
			coll.OrganisationalUnit,
		).Add(float64(coll.QuotaInBytes))

		storageUsage.WithLabelValues(
			coll.ProjectID,
			coll.Type.String(),
			coll.State.String(),
			coll.OrganisationalUnit,
		).Add(float64(coll.SizeInBytes))
	}

	// users
	userCount.Reset()
	userCountOu.Reset()
	for _, user := range ucache.GetUsers() {
		userCount.WithLabelValues(
			user.IdentityProvider,
		).Add(1)

		for _, ou := range user.OrganizationalUnits {
			userCountOu.WithLabelValues(
				ou,
			).Add(1)
		}
	}

}
