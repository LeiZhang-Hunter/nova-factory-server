package catalogwire

import (
	"strings"
	"testing"
	"time"
)

func TestSignedCatalogIdentityFreshnessAndError(t *testing.T) {
	now := time.Unix(1800000000, 0)
	req := Request{
		RequestID: strings.Repeat("a1", 24),
		GatewayID: "42",
		Username:  "agent",
		Timestamp: now.Unix(),
	}
	req.Signature = req.Sign("secret")
	if !req.Verify("secret", now) || req.Verify("wrong", now) || req.Verify("secret", now.Add(2*time.Minute)) {
		t.Fatal("signature/freshness failed")
	}
	changed := req
	changed.GatewayID = "43"
	if changed.Verify("secret", now) {
		t.Fatal("gateway substitution accepted")
	}
	reply := Response{
		GatewayID: req.GatewayID,
		RequestID: req.RequestID,
		Error:     "catalog unavailable",
	}
	reply.Signature = reply.Sign("secret")
	if !reply.Authentic("secret", req) || reply.Verify("secret", req) == nil {
		t.Fatal("signed errors must authenticate but report failure")
	}
	reply.GatewayID = "43"
	if reply.Authentic("secret", req) {
		t.Fatal("foreign response accepted")
	}
	if TopicGateway(RequestTopic("42")) != "42" || TopicGateway(RequestTopic("42/43")) != "" {
		t.Fatal("topic identity invalid")
	}
}
