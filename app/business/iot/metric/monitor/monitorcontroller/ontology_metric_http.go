// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"
	"nova-factory-server/app/business/iot/metric/monitor/monitorservice"
)

type OntologyMetricHTTPHandler struct {
	service monitorservice.IOntologyMetricReportService
	logger  *zap.Logger
}

func NewOntologyMetricHTTPHandler(service monitorservice.IOntologyMetricReportService, logger *zap.Logger) *OntologyMetricHTTPHandler {
	return &OntologyMetricHTTPHandler{
		service: service,
		logger:  logger,
	}
}

func (h *OntologyMetricHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"code": http.StatusMethodNotAllowed, "msg": "method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "msg": fmt.Sprintf("read body: %v", err)})
		return
	}
	body = bytes.TrimSpace(body)
	if len(body) == 0 {
		h.writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "msg": "empty body"})
		return
	}

	accepted, err := h.receive(r, body)
	if err != nil {
		h.logger.Error("receive ontology metrics failed", zap.Error(err))
		h.writeJSON(w, http.StatusBadRequest, map[string]any{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	h.writeJSON(w, http.StatusAccepted, map[string]any{"code": 0, "msg": "ok", "data": map[string]any{"accepted": accepted}})
}

func (h *OntologyMetricHTTPHandler) receive(r *http.Request, body []byte) (int, error) {
	if body[0] == '[' {
		var reports []modelrequest.OntologyMetricReport
		if err := json.Unmarshal(body, &reports); err != nil {
			return 0, fmt.Errorf("decode ontology metric reports: %w", err)
		}
		return h.service.ReceiveOntologyReports(r.Context(), reports)
	}

	var payload modelrequest.ExportDML
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0, fmt.Errorf("decode monitor export payload: %w", err)
	}
	if payload.Signal == "" && len(payload.Body) == 0 {
		return 0, fmt.Errorf("unsupported monitor payload")
	}
	if err := h.service.ReceiveExportPayload(r.Context(), payload); err != nil {
		return 0, err
	}
	return 1, nil
}

func (h *OntologyMetricHTTPHandler) writeJSON(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Debug("write ontology metrics response failed", zap.Error(err))
	}
}
