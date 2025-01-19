import { orderBy } from "lodash-es";
import type { SQLResultSetV1 } from "@/types";
import { NullValue } from "@/types/proto/google/protobuf/struct";
import { Timestamp } from "@/types/proto/google/protobuf/timestamp";
import { QueryRow, RowValue } from "@/types/proto/v1/sql_service";

type NoSQLRowData = {
  key: string;
  value: any;
};

const flattenNoSQLColumn = (value: any): any => {
  if (typeof value !== "object") {
    return value;
  }
  if (value === null) {
    return value;
  }
  if (Array.isArray(value)) {
    return value.map(flattenNoSQLColumn);
  }

  const dict = value as { [key: string]: any };
  if (Object.keys(dict).length === 1 && Object.keys(dict)[0].startsWith("$")) {
    // Used by the MongoDB response.
    const key = Object.keys(dict)[0];
    switch (key) {
      case "$oid":
        return dict[key];
      case "$date":
        if (typeof dict[key] !== "object") {
          return dict[key];
        }
        if (!dict[key]["$numberLong"]) {
          return dict[key];
        }
        return new Date(parseInt(dict[key]["$numberLong"]));
      case "$numberLong":
        return parseInt(dict[key]);
      case "$numberDouble":
        return parseFloat(dict[key]);
      case "$numberInt":
        return parseInt(dict[key]);
      case "$timestamp":
        return (dict[key] as { t: number; i: number }).t;
      default:
        return dict[key];
    }
  }

  return Object.keys(dict).reduce(
    (d, key) => {
      d[key] = flattenNoSQLColumn(dict[key]);
      return d;
    },
    {} as { [key: string]: any }
  );
};

const convertAnyToRowValue = (
  value: any,
  nested: boolean
): { value: RowValue; type: string } => {
  switch (typeof value) {
    case "number": {
      if (Math.floor(value) === value) {
        return {
          value: RowValue.fromPartial({
            int32Value: value,
          }),
          type: "INTEGER",
        };
      }
      return {
        value: RowValue.fromPartial({
          floatValue: value,
        }),
        type: "FLOAT",
      };
    }
    case "string":
      return {
        value: RowValue.fromPartial({
          stringValue: value,
        }),
        type: "TEXT",
      };
    case "undefined":
      return {
        value: RowValue.fromPartial({
          nullValue: NullValue.NULL_VALUE,
        }),
        type: "NULL",
      };
    case "boolean":
      return {
        value: RowValue.fromPartial({
          boolValue: value,
        }),
        type: "BOOLEAN",
      };
    case "bigint":
      return {
        value: RowValue.fromPartial({
          stringValue: value.toString(),
        }),
        type: "TEXT",
      };
    case "object": {
      if (value === null) {
        return {
          value: RowValue.fromPartial({
            nullValue: NullValue.NULL_VALUE,
          }),
          type: "NULL",
        };
      }
      if (Array.isArray(value)) {
        return {
          value: RowValue.fromPartial({
            stringValue: JSON.stringify(value.map(flattenNoSQLColumn)),
          }),
          type: "OBJECT",
        };
      }
      if (value instanceof Date) {
        return {
          value: RowValue.fromPartial({
            timestampValue: Timestamp.fromPartial({
              seconds: Math.floor(value.getTime() / 1000),
              nanos: (value.getTime() % 1000) * 1e6,
            }),
          }),
          type: "DATETIME",
        };
      }

      if (nested) {
        const formatted = flattenNoSQLColumn(value);
        return convertAnyToRowValue(formatted, !nested);
      } else {
        return {
          value: RowValue.fromPartial({
            stringValue: JSON.stringify(value),
          }),
          type: "TEXT",
        };
      }
    }
    default:
      return {
        value: RowValue.fromPartial({
          stringValue: JSON.stringify(value),
        }),
        type: "TEXT",
      };
  }
};

export const flattenNoSQLResult = (resultSet: SQLResultSetV1) => {
  for (const result of resultSet.results) {
    const { columns, columnIndexMap } = getNoSQLColumns(result.rows);

    const rows: QueryRow[] = [];
    const columnTypeNames: string[] = Array.from({
      length: columns.length,
    }).map((_) => "TEXT");

    for (const row of result.rows) {
      if (row.values.length !== 1 || !row.values[0].stringValue) {
        continue;
      }
      const data = JSON.parse(row.values[0].stringValue);
      const values: RowValue[] = Array.from({ length: columns.length }).map(
        (_) =>
          RowValue.fromPartial({
            nullValue: NullValue.NULL_VALUE,
          })
      );

      for (const [key, value] of Object.entries(data)) {
        const index = columnIndexMap.get(key) ?? 0;

        const { value: formatted, type } = convertAnyToRowValue(value, true);
        values[index] = formatted;
        columnTypeNames[index] = type;
      }

      rows.push(
        QueryRow.fromPartial({
          values: values,
        })
      );
    }

    result.rows = rows;
    result.columnNames = columns;
    result.columnTypeNames = columnTypeNames;
  }
};

const getNoSQLColumns = (rows: QueryRow[]) => {
  const columnSet = new Set<string>();
  const columnIndexMap = new Map<string, number>();

  for (const row of rows) {
    const parsedRows = getNoSQLRows(row);
    if (!parsedRows) {
      continue;
    }
    for (const item of parsedRows) {
      columnSet.add(item.key);
    }
  }

  const builtInColumns = new Map<string, number>([
    ["_rid", 0],
    ["_self", 1],
    ["_etag", 2],
    ["_attachments", 3],
    ["_ts", 4],
  ]);

  const sortedColumns = orderBy(
    [...columnSet],
    [
      (column) => (column === "id" || column === "-id" ? -1 : 1),
      (column) => (builtInColumns.has(column) ? 1 : 0),
      (column) => builtInColumns.get(column) ?? 0,
    ],
    ["asc", "asc"]
  );

  for (let i = 0; i < sortedColumns.length; i++) {
    columnIndexMap.set(sortedColumns[i], i);
  }

  return {
    columns: sortedColumns,
    columnIndexMap,
  };
};

const getNoSQLRows = (row: QueryRow): NoSQLRowData[] | undefined => {
  if (row.values.length !== 1 || !row.values[0].stringValue) {
    return;
  }
  const parsedRow = JSON.parse(row.values[0].stringValue) as {
    [key: string]: any;
  };
  const results: NoSQLRowData[] = [];

  for (const [key, value] of Object.entries(parsedRow)) {
    results.push({
      key: key,
      value: value,
    });
  }
  return results;
};
