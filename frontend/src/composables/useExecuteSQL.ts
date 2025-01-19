import Emittery from "emittery";
import { head, isEmpty } from "lodash-es";
import { Status } from "nice-grpc-common";
import { markRaw, reactive } from "vue";
import { parseSQL } from "@/components/MonacoEditor/sqlParser";
import { sqlServiceClient } from "@/grpcweb";
import { t } from "@/plugins/i18n";
import {
  pushNotification,
  useDatabaseV1Store,
  useSQLEditorStore,
  useSQLEditorTabStore,
  useSQLStore,
  useSQLEditorQueryHistoryStore,
  useAppFeature,
} from "@/store";
import type {
  ComposedDatabase,
  SQLResultSetV1,
  BBNotificationStyle,
  SQLEditorQueryParams,
  SQLEditorTab,
} from "@/types";
import { isValidDatabaseName } from "@/types";
import { Engine } from "@/types/proto/v1/common";
import {
  Advice,
  Advice_Status,
  advice_StatusToJSON,
  CheckRequest_ChangeType,
} from "@/types/proto/v1/sql_service";
import {
  emptySQLEditorTabQueryContext,
  ensureDataSourceSelection,
  hasPermissionToCreateChangeDatabaseIssue,
} from "@/utils";
import { extractGrpcErrorMessage } from "@/utils/grpcweb";
import { flattenNoSQLResult } from "./utils";

// SKIP_CHECK_THRESHOLD is the MaxSheetCheckSize in the backend.
const SKIP_CHECK_THRESHOLD = 2 * 1024 * 1024;
// QUERY_INTERVAL_LIMIT is the minimal gap between two queries
const QUERY_INTERVAL_LIMIT = 1000;

const events = new Emittery<{
  "update:advices": {
    tab: SQLEditorTab;
    params: SQLEditorQueryParams;
    advices: Advice[];
  };
}>();

type SQLCheckResult = { passed: boolean; advices?: Advice[] };

const useExecuteSQL = () => {
  const state = reactive<{
    lastQueryTime?: number;
  }>({});
  const databaseStore = useDatabaseV1Store();
  const tabStore = useSQLEditorTabStore();
  const sqlEditorStore = useSQLEditorStore();
  const sqlCheckStyle = useAppFeature("bb.feature.sql-editor.sql-check-style");

  const notify = (
    type: BBNotificationStyle,
    title: string,
    description?: string
  ) => {
    pushNotification({
      module: "bytebase",
      style: type,
      title,
      description,
    });
  };

  const preflight = async (params: SQLEditorQueryParams) => {
    state.lastQueryTime = Date.now();

    const tab = tabStore.currentTab;
    if (!tab) {
      return false;
    }

    if (tab.queryContext?.status === "EXECUTING") {
      notify("INFO", t("common.tips"), t("sql-editor.can-not-execute-query"));
      return false;
    }

    if (tabStore.isDisconnected) {
      notify("CRITICAL", t("sql-editor.select-connection"));
      return false;
    }

    if (isEmpty(params.statement)) {
      notify("CRITICAL", t("sql-editor.notify-empty-statement"));
      return false;
    }

    tab.queryContext = {
      ...emptySQLEditorTabQueryContext(),
      params,
      status: "EXECUTING",
    };
    return true;
  };

  const check = async (): Promise<SQLCheckResult> => {
    const tab = tabStore.currentTab;
    if (!tab) {
      return { passed: false };
    }

    const params = tab.queryContext?.params;
    const abortController = tab.queryContext?.abortController;
    if (!params) {
      return { passed: false };
    }
    if (new Blob([params.statement]).size > SKIP_CHECK_THRESHOLD) {
      return { passed: true };
    }
    const response = await sqlServiceClient.check(
      {
        name: params.connection.database,
        statement: params.statement,
        changeType: CheckRequest_ChangeType.SQL_EDITOR,
      },
      {
        ignoredCodes: [Status.PERMISSION_DENIED],
        signal: abortController?.signal,
      }
    );
    const { advices } = response;
    events.emit("update:advices", { tab, params, advices });
    return { passed: advices.length === 0, advices };
  };

  const notifyAdvices = (advices: Advice[]) => {
    let adviceStatus: "SUCCESS" | "ERROR" | "WARNING" = "SUCCESS";
    let adviceNotifyMessage = "";
    for (const advice of advices) {
      if (advice.status === Advice_Status.SUCCESS) {
        continue;
      }

      if (advice.status === Advice_Status.ERROR) {
        adviceStatus = "ERROR";
      } else if (adviceStatus !== "ERROR") {
        adviceStatus = "WARNING";
      }

      adviceNotifyMessage += `${advice_StatusToJSON(advice.status)}: ${
        advice.title
      }\n`;
      if (advice.content) {
        adviceNotifyMessage += `${advice.content}\n`;
      }
    }

    if (adviceStatus !== "SUCCESS") {
      const notifyStyle = adviceStatus === "ERROR" ? "CRITICAL" : "WARN";
      notify(
        notifyStyle,
        t("sql-editor.sql-review-result"),
        adviceNotifyMessage
      );
    }
  };

  const cleanup = () => {
    const tab = tabStore.currentTab;
    if (!tab) return;
    if (!tab.queryContext) return;
    tab.queryContext.status = "IDLE";
  };

  const execute = async (params: SQLEditorQueryParams) => {
    const now = Date.now();
    if (
      state.lastQueryTime &&
      now - state.lastQueryTime < QUERY_INTERVAL_LIMIT
    ) {
      return;
    }

    if (!preflight(params)) {
      return cleanup();
    }

    const tab = tabStore.currentTab;
    if (!tab) {
      return cleanup();
    }
    const { mode } = tab;
    if (mode === "ADMIN") {
      return cleanup();
    }

    const queryContext = tab.queryContext!;
    const batchQueryContext = tab.batchQueryContext;
    const { data } = await parseSQL(params.statement);

    if (data === undefined) {
      notify("CRITICAL", t("sql-editor.notify-invalid-sql-statement"));
      return cleanup();
    }

    const selectedDatabase = useDatabaseV1Store().getDatabaseByName(
      params.connection.database
    );
    const databaseName = isValidDatabaseName(selectedDatabase.name)
      ? selectedDatabase.databaseName
      : "";
    const batchQueryDatabases: ComposedDatabase[] = [selectedDatabase];

    // Check if the user selects multiple databases to query.
    if (
      databaseName &&
      batchQueryContext &&
      batchQueryContext.databases.length > 0
    ) {
      for (const databaseResourceName of batchQueryContext.databases) {
        const database = databaseStore.getDatabaseByName(databaseResourceName);
        if (database.name === selectedDatabase.name) {
          continue;
        }
        batchQueryDatabases.push(database);
      }
    }

    const queryResultMap = new Map<string, SQLResultSetV1>();
    for (const database of batchQueryDatabases) {
      queryResultMap.set(database.name, {
        error: "",
        results: [],
        advices: [],
        allowExport: false,
      });
    }
    queryContext.results = queryResultMap;

    const fail = (database: ComposedDatabase, result: SQLResultSetV1) => {
      queryResultMap.set(database.name, {
        error: result.error,
        results: [],
        advices: result.advices,
        status: result.status,
        allowExport: false,
      });
    };
    const abort = (error: string, advices: Advice[] = []) => {
      fail(batchQueryDatabases[0], {
        error,
        results: [],
        advices,
        status: Status.ABORTED,
        allowExport: false,
      });
      return cleanup();
    };

    const { abortController } = queryContext;

    const sqlStore = useSQLStore();
    queryContext.beginTimestampMS = Date.now();

    const checkBehavior = params.skipCheck ? "SKIP" : sqlCheckStyle.value;
    let checkResult: SQLCheckResult = { passed: true };
    if (checkBehavior !== "SKIP") {
      try {
        checkResult = await check();
      } catch (error) {
        return abort(extractGrpcErrorMessage(error));
      }
    }
    if (checkBehavior === "PREFLIGHT" && !checkResult.passed) {
      return abort(
        head(checkResult.advices)?.content ?? "",
        checkResult.advices
      );
    }
    if (
      checkBehavior === "NOTIFICATION" &&
      !checkResult.passed &&
      checkResult.advices
    ) {
      const { advices } = checkResult;
      const errorAdvice = advices.find(
        (advice) => advice.status === Advice_Status.ERROR
      );
      if (errorAdvice) {
        notifyAdvices(advices);
        return abort(errorAdvice.content, advices);
      }
    }

    for (const database of batchQueryDatabases) {
      if (abortController.signal.aborted) {
        // Once any one of the batch queries is aborted, don't go further
        // and mock an "Aborted" result for the rest queries.
        fail(database, {
          advices: [],
          allowExport: false,
          error: "AbortError: The user aborted a request.",
          results: [],
          status: Status.ABORTED,
        });
        continue;
      }

      try {
        const instance = isValidDatabaseName(database.name)
          ? database.instance
          : params.connection.instance;
        const dataSourceId =
          ensureDataSourceSelection(
            instance === params.connection.instance
              ? params.connection.dataSourceId
              : undefined,
            database
          ) ?? "";

        if (!dataSourceId) {
          fail(database, {
            advices: [],
            allowExport: false,
            error: "No queriable data source.",
            results: [],
            status: Status.NOT_FOUND,
          });
          continue;
        }

        const resultSet = await sqlStore.query(
          {
            name: database.name,
            dataSourceId: dataSourceId,
            statement: params.statement,
            limit: sqlEditorStore.resultRowsLimit,
            explain: params.explain,
            schema: params.connection.schema,
            container: params.connection.table,
            timeout: undefined, // TODO: make this param configurable
            queryOption: {
              redisRunCommandsOn: sqlEditorStore.redisCommandOption,
            },
          },
          abortController.signal
        );

        if (
          database.instanceResource.engine === Engine.MONGODB ||
          database.instanceResource.engine === Engine.COSMOSDB
        ) {
          flattenNoSQLResult(resultSet);
        }

        if (checkBehavior === "NOTIFICATION") {
          notifyAdvices(checkResult.advices ?? []);
        }

        if (resultSet.error) {
          // The error message should be consistent with the one from the backend.
          if (isOnlySelectError(resultSet)) {
            const database = databaseStore.getDatabaseByName(
              params.connection.database
            );

            // Show a tips to navigate to issue creation
            // if the user is allowed to create issue in the project.
            if (hasPermissionToCreateChangeDatabaseIssue(database)) {
              sqlEditorStore.isShowExecutingHint = true;
              sqlEditorStore.executingHintDatabase = database;
              cleanup();
            }
            fail(database, resultSet);
          } else {
            fail(database, resultSet);
          }
        } else {
          queryResultMap.set(database.name, markRaw(resultSet));
        }
      } catch (error: any) {
        fail(database, error.response?.data?.message ?? String(error));
      }
    }

    // After all the queries are executed, we update the tab with the latest query result map.
    // Refresh the query history list when the query executed successfully
    // (with or without warnings).
    useSQLEditorQueryHistoryStore().fetchQueryHistoryList();
    cleanup();
  };

  return {
    events,
    execute,
  };
};

const isOnlySelectError = (resultSet: SQLResultSetV1) => {
  if (
    resultSet.error === "Support SELECT sql statement only" &&
    resultSet.status === Status.INVALID_ARGUMENT
  ) {
    return true;
  }
  if (
    resultSet.error.match(/disallow execute (DML|DDL) statement/) &&
    resultSet.status === Status.PERMISSION_DENIED
  ) {
    return true;
  }
  return false;
};

export { useExecuteSQL };
