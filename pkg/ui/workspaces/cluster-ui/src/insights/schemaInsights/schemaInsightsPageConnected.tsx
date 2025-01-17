// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

import { connect } from "react-redux";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { Dispatch } from "redux";

import {
  actions,
  selectSchemaInsights,
  selectSchemaInsightsDatabases,
  selectSchemaInsightsError,
  selectSchemaInsightsMaxApiSizeReached,
  selectSchemaInsightsTypes,
  selectFilters,
  selectSortSetting,
} from "src/store/schemaInsights";
import { AppState, uiConfigActions } from "src/store";
import { SortSetting } from "src/sortedtable";
import { selectDropUnusedIndexDuration } from "src/store/clusterSettings/clusterSettings.selectors";

import { SchemaInsightEventFilters } from "../types";
import { actions as localStorageActions } from "../../store/localStorage";
import { selectHasAdminRole } from "../../store/uiConfig";
import { actions as analyticsActions } from "../../store/analytics";

import {
  SchemaInsightsView,
  SchemaInsightsViewDispatchProps,
  SchemaInsightsViewStateProps,
} from "./schemaInsightsView";

const mapStateToProps = (
  state: AppState,
  _props: RouteComponentProps,
): SchemaInsightsViewStateProps => ({
  schemaInsights: selectSchemaInsights(state),
  schemaInsightsDatabases: selectSchemaInsightsDatabases(state),
  schemaInsightsTypes: selectSchemaInsightsTypes(state),
  schemaInsightsError: selectSchemaInsightsError(state),
  filters: selectFilters(state),
  sortSetting: selectSortSetting(state),
  hasAdminRole: selectHasAdminRole(state),
  maxSizeApiReached: selectSchemaInsightsMaxApiSizeReached(state),
  csIndexUnusedDuration: selectDropUnusedIndexDuration(state),
});

const mapDispatchToProps = (
  dispatch: Dispatch,
): SchemaInsightsViewDispatchProps => ({
  onFiltersChange: (filters: SchemaInsightEventFilters) => {
    dispatch(
      localStorageActions.update({
        key: "filters/SchemaInsightsPage",
        value: filters,
      }),
    );
    dispatch(
      analyticsActions.track({
        name: "Filter Clicked",
        page: "Schema Insights",
        filterName: "filters",
        value: filters.toString(),
      }),
    );
  },
  onSortChange: (ss: SortSetting) => {
    dispatch(
      localStorageActions.update({
        key: "sortSetting/SchemaInsightsPage",
        value: ss,
      }),
    );
    dispatch(
      analyticsActions.track({
        name: "Column Sorted",
        page: "Schema Insights",
        tableName: "Schema Insights Table",
        columnName: ss.columnTitle,
      }),
    );
  },
  refreshSchemaInsights: (csIndexUnusedDuration: string) => {
    dispatch(actions.refresh({ csIndexUnusedDuration }));
  },
  refreshUserSQLRoles: () => dispatch(uiConfigActions.refreshUserSQLRoles()),
});

export const SchemaInsightsPageConnected = withRouter(
  connect<
    SchemaInsightsViewStateProps,
    SchemaInsightsViewDispatchProps,
    RouteComponentProps
  >(
    mapStateToProps,
    mapDispatchToProps,
  )(SchemaInsightsView),
);
