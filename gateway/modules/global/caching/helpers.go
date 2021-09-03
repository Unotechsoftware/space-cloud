package caching

import (
	"encoding/json"
	"fmt"
	"github.com/spaceuptech/space-cloud/gateway/config"
)

const (
	databaseJoinTypeResult = "result"
	databaseJoinTypeJoin   = "join"
	databaseJoinTypeAlways = "always"

	keyTypeTTL        = "ttl"
	keyTypeInvalidate = "invalidate"
)

// Ingress keys
func (c *Cache) generateIngressRoutingKey(routeID string, cacheOptions []interface{}) string {
	data, _ := json.Marshal(cacheOptions)
	return fmt.Sprintf("%s::%s", c.generateIngressRoutingPrefixWithRouteID(routeID), string(data))
}

func (c *Cache) generateIngressRoutingPrefixWithRouteID(routeID string) string {
	return fmt.Sprintf("%s::%s", c.generateIngressRoutingResourcePrefixKey(), routeID)
}

func (c *Cache) generateIngressRoutingResourcePrefixKey() string {
	return fmt.Sprintf("%s::%s", c.clusterID, config.ResourceIngressRoute)
}