package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v32 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/parse"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/rancher/pkg/clustermanager"
	v3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	pv3 "github.com/rancher/rancher/pkg/generated/norman/project.cattle.io/v3"
	monitorutil "github.com/rancher/rancher/pkg/monitoring"
	"github.com/rancher/rancher/pkg/types/config/dialer"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewProjectGraphHandler(dialerFactory dialer.Factory, clustermanager *clustermanager.Manager) *ProjectGraphHandler {
	return &ProjectGraphHandler{
		dialerFactory:  dialerFactory,
		clustermanager: clustermanager,
		projectLister:  clustermanager.ScaledContext.Management.Projects(metav1.NamespaceAll).Controller().Lister(),
		appLister:      clustermanager.ScaledContext.Project.Apps(metav1.NamespaceAll).Controller().Lister(),
	}
}

type ProjectGraphHandler struct {
	dialerFactory  dialer.Factory
	clustermanager *clustermanager.Manager
	projectLister  v3.ProjectLister
	appLister      pv3.AppLister
}

func (h *ProjectGraphHandler) QuerySeriesAction(actionName string, action *types.Action, apiContext *types.APIContext) error {
	var queryGraphInput v32.QueryGraphInput
	actionInput, err := parse.ReadBody(apiContext.Request)
	if err != nil {
		return err
	}

	if err = convert.ToObj(actionInput, &queryGraphInput); err != nil {
		return err
	}

	inputParser := newProjectGraphInputParser(queryGraphInput)
	if err = inputParser.parse(); err != nil {
		return err
	}

	clusterName := inputParser.ClusterName
	userContext, err := h.clustermanager.UserContextNoControllers(clusterName)
	if err != nil {
		return fmt.Errorf("get usercontext failed, %v", err)
	}

	check := newAuthChecker(apiContext.Request.Context(), userContext, inputParser.Input, inputParser.ProjectID)
	if err = check.check(); err != nil {
		return err
	}

	reqContext, cancel := context.WithTimeout(context.Background(), prometheusReqTimeout)
	defer cancel()

	var svcName, svcNamespace, svcPort, token string
	var queries []*PrometheusQuery
	prometheusName, prometheusNamespace := monitorutil.ClusterMonitoringInfo()
	token, err = getAuthToken(userContext, prometheusName, prometheusNamespace)
	if err != nil {
		return err
	}

	prometheusQuery, err := NewPrometheusQuery(reqContext, clusterName, token, svcNamespace, svcName, svcPort, h.dialerFactory, userContext)
	if err != nil {
		return err
	}
	seriesSlice, err := prometheusQuery.Do(queries)
	if err != nil {
		logrus.WithError(err).Warn("query series failed")
		return httperror.NewAPIError(httperror.ServerError, "Failed to obtain metrics. The metrics service may not be available.")
	}

	if seriesSlice == nil {
		apiContext.WriteResponse(http.StatusNoContent, nil)
		return nil
	}

	collection := v32.QueryProjectGraphOutput{Type: "collection"}
	for k, v := range seriesSlice {
		graphName, _, _ := parseID(k)
		queryGraph := v32.QueryProjectGraph{
			GraphName: graphName,
			Series:    parseResponse(v),
		}
		collection.Data = append(collection.Data, queryGraph)
	}

	res, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("marshal query series result failed, %v", err)
	}
	apiContext.Response.Write(res)
	return nil
}

func parseResponse(seriesSlice []*TimeSeries) []*v32.TimeSeries {
	var series []*v32.TimeSeries
	for _, v := range seriesSlice {
		series = append(series, &v32.TimeSeries{
			Name:   v.Name,
			Points: v.Points,
		})
	}
	return series
}
