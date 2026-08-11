package iotdb

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/iotdb-client-go/client"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
)

const deviceTemplateValueType = "DOUBLE"

// deviceTemplate 设备模板
type deviceTemplate struct {
	iotDb *iotdb.IotDb
}

func newDeviceTemplate(iotDb *iotdb.IotDb) *deviceTemplate {
	return &deviceTemplate{
		iotDb: iotDb,
	}
}

// Delete 删除模板
func (d *deviceTemplate) Delete(template entity.DeviceTemplate) error {
	return nil
}

// Update 更新模板
func (d *deviceTemplate) Update(template []entity.DeviceTemplate) error {
	if len(template) == 0 {
		return nil
	}
	if d == nil || d.iotDb == nil {
		return fmt.Errorf("iotdb device template storage is nil")
	}

	templateData, err := groupTemplateData(template)
	if err != nil {
		return err
	}
	if len(templateData) == 0 {
		return nil
	}

	session, err := d.iotDb.GetSession()
	if err != nil {
		zap.L().Error("get iotdb session error", zap.Error(err))
		return err
	}
	defer d.iotDb.PutSession(session)

	templateIds := make([]int64, 0, len(templateData))
	for templateId := range templateData {
		templateIds = append(templateIds, templateId)
	}
	sort.Slice(templateIds, func(i, j int) bool { return templateIds[i] < templateIds[j] })

	for _, templateId := range templateIds {
		if err := d.upsertDeviceTemplate(session, templateId, templateData[templateId]); err != nil {
			return err
		}
	}
	return nil
}

// upsertDeviceTemplate 更新设备模板
func (d *deviceTemplate) upsertDeviceTemplate(session client.Session, templateId int64, dataIds []int64) error {
	name := iotdb2.MakeDeviceTemplateName(templateId)
	exists, err := d.deviceTemplateExists(session, name)
	if err != nil {
		return fmt.Errorf("check iotdb device template %s: %w", name, err)
	}

	if !exists {
		if len(dataIds) == 0 {
			return nil
		}
		return d.createDeviceTemplate(session, name, dataIds)
	}

	existingNodes, err := d.deviceTemplateNodes(session, name)
	if err != nil {
		return fmt.Errorf("read iotdb device template nodes %s: %w", name, err)
	}

	desiredNodes := make(map[string]bool, len(dataIds))
	missingDataIds := make([]int64, 0)
	for _, dataId := range dataIds {
		node := deviceTemplateNodeName(dataId)
		desiredNodes[node] = true
		if !existingNodes[node] {
			missingDataIds = append(missingDataIds, dataId)
		}
	}

	for node := range existingNodes {
		if !desiredNodes[node] {
			return d.recreateDeviceTemplate(session, name, dataIds)
		}
	}
	if len(missingDataIds) == 0 {
		return nil
	}

	sql := fmt.Sprintf("alter device template %s add (%s)", name, deviceTemplateDefinitions(missingDataIds))
	if err := d.executeStatement(session, sql); err != nil {
		zap.L().Error("alter iotdb device template error", zap.String("template", name), zap.Error(err))
		return fmt.Errorf("alter iotdb device template %s: %w", name, err)
	}
	return nil
}

func (d *deviceTemplate) createDeviceTemplate(session client.Session, name string, dataIds []int64) error {
	sql := fmt.Sprintf("create device template %s aligned (%s)", name, deviceTemplateDefinitions(dataIds))
	if err := d.executeStatement(session, sql); err != nil {
		zap.L().Error("create iotdb device template error", zap.String("template", name), zap.Error(err))
		return fmt.Errorf("create iotdb device template %s: %w", name, err)
	}
	return nil
}

func (d *deviceTemplate) recreateDeviceTemplate(session client.Session, name string, dataIds []int64) error {
	if err := d.dropDeviceTemplate(session, name); err != nil {
		return err
	}
	if len(dataIds) == 0 {
		return nil
	}
	return d.createDeviceTemplate(session, name, dataIds)
}

func (d *deviceTemplate) dropDeviceTemplate(session client.Session, name string) error {
	sql := fmt.Sprintf("drop device template %s", name)
	if err := d.executeStatement(session, sql); err != nil {
		zap.L().Error("drop iotdb device template error", zap.String("template", name), zap.Error(err))
		return fmt.Errorf("drop iotdb device template %s: %w", name, err)
	}
	return nil
}

// groupTemplateData 分组创建
func groupTemplateData(template []entity.DeviceTemplate) (map[int64][]int64, error) {
	data := make(map[int64][]int64)
	seen := make(map[int64]map[int64]struct{})
	for _, item := range template {
		if item.TemplateId <= 0 {
			return nil, fmt.Errorf("invalid device template id: %d", item.TemplateId)
		}
		if _, ok := data[item.TemplateId]; !ok {
			data[item.TemplateId] = make([]int64, 0)
		}
		if item.DataId <= 0 {
			continue
		}
		if seen[item.TemplateId] == nil {
			seen[item.TemplateId] = make(map[int64]struct{})
		}
		if _, ok := seen[item.TemplateId][item.DataId]; ok {
			continue
		}
		seen[item.TemplateId][item.DataId] = struct{}{}
		data[item.TemplateId] = append(data[item.TemplateId], item.DataId)
	}

	for templateId := range data {
		sort.Slice(data[templateId], func(i, j int) bool { return data[templateId][i] < data[templateId][j] })
	}
	return data, nil
}

func (d *deviceTemplate) deviceTemplateExists(session client.Session, name string) (bool, error) {
	statement, err := session.ExecuteStatement("show device templates")
	if err != nil {
		return false, err
	}
	defer statement.Close()

	column := firstStatementColumn(statement, "template name", "template", "name")
	for next, err := statement.Next(); ; next, err = statement.Next() {
		if err != nil {
			return false, err
		}
		if !next {
			break
		}
		if strings.TrimSpace(statement.GetText(column)) == name {
			return true, nil
		}
	}
	return false, nil
}

func (d *deviceTemplate) deviceTemplateNodes(session client.Session, name string) (map[string]bool, error) {
	statement, err := session.ExecuteStatement(fmt.Sprintf("show nodes in device template %s", name))
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	nodes := make(map[string]bool)
	column := firstStatementColumn(statement, "child nodes", "child node", "node")
	for next, err := statement.Next(); ; next, err = statement.Next() {
		if err != nil {
			return nil, err
		}
		if !next {
			break
		}
		node := strings.TrimSpace(statement.GetText(column))
		if node != "" {
			nodes[node] = true
		}
	}
	return nodes, nil
}

func (d *deviceTemplate) executeStatement(session client.Session, sql string) error {
	statement, err := session.ExecuteStatement(sql)
	if statement != nil {
		defer statement.Close()
	}
	return err
}

func firstStatementColumn(statement *client.SessionDataSet, names ...string) string {
	columns := statement.GetColumnNames()
	for _, name := range names {
		for _, column := range columns {
			if strings.EqualFold(column, name) {
				return column
			}
		}
	}
	if len(columns) == 0 {
		return ""
	}
	return columns[0]
}

func deviceTemplateDefinitions(dataIds []int64) string {
	definitions := make([]string, 0, len(dataIds))
	for _, dataId := range dataIds {
		definitions = append(definitions, fmt.Sprintf("%s %s", deviceTemplateNodeName(dataId), deviceTemplateValueType))
	}
	return strings.Join(definitions, ", ")
}

func deviceTemplateNodeName(dataId int64) string {
	return "d" + strconv.FormatInt(dataId, 10)
}
