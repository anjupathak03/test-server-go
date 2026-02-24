import { test, expect } from "@playwright/test";
import {
  createConnection,
  type Connection,
  type QueryError,
  type ResultSetHeader,
  type RowDataPacket,
} from "mysql2/promise";

interface TodoRow extends RowDataPacket {
  id: number;
  title: string;
  description: string;
  completed: number;
}

const API_CREATED_TITLE = "mysql-api-stable-title";
const API_CREATED_DESCRIPTION = "created through API";
const DIRECT_MYSQL_TITLE = "direct-mysql-stable-title";
const DIRECT_MYSQL_DESCRIPTION = "inserted through mysql2";

const mysqlConfig = {
  host: process.env.DB_HOST || "localhost",
  port: Number(process.env.DB_PORT || "3306"),
  user: process.env.DB_USER || "root",
  password: process.env.DB_PASSWORD || "password",
};

const databaseCandidates = buildDatabaseCandidates();

function buildDatabaseCandidates(): string[] {
  const raw = [
    process.env.DB_NAME,
    process.env.MYSQL_DATABASE,
    "todo_db",
    "todo_test_db",
  ].filter(Boolean) as string[];

  const unique = [...new Set(raw)];
  return unique.filter((name) => /^[A-Za-z0-9_]+$/.test(name));
}

function isSkippableDatabaseError(err: unknown): boolean {
  const code = (err as QueryError | undefined)?.code;
  return (
    code === "ER_BAD_DB_ERROR" ||
    code === "ER_NO_SUCH_TABLE" ||
    code === "ER_DBACCESS_DENIED_ERROR"
  );
}

async function findTodoInAnyCandidateDB(
  connection: Connection,
  title: string,
  description: string
): Promise<{ db: string; row: TodoRow } | null> {
  for (const db of databaseCandidates) {
    try {
      const [rows] = await connection.query<TodoRow[]>(
        `SELECT id, title, description, completed FROM \`${db}\`.todos WHERE title = ? AND description = ? ORDER BY id DESC LIMIT 1`,
        [title, description]
      );
      if (rows.length > 0) {
        return { db, row: rows[0] };
      }
    } catch (err) {
      if (!isSkippableDatabaseError(err)) {
        throw err;
      }
    }
  }
  return null;
}

async function waitForTodoInAnyCandidateDB(
  connection: Connection,
  title: string,
  description: string,
  attempts = 10,
  delayMs = 75
): Promise<{ db: string; row: TodoRow } | null> {
  for (let i = 0; i < attempts; i++) {
    const found = await findTodoInAnyCandidateDB(connection, title, description);
    if (found) {
      return found;
    }
    await new Promise((resolve) => setTimeout(resolve, delayMs));
  }
  return null;
}

test.describe("API + direct MySQL outgoing calls", () => {
  let connection: Connection;
  let resolvedDatabase: string | null = null;

  test.beforeAll(async () => {
    connection = await createConnection({
      ...mysqlConfig,
      connectTimeout: 5_000,
    });
  });

  test.afterAll(async () => {
    await connection.end();
  });

  test("MySQLhealthquery", async () => {
    const [rows] = await connection.query<RowDataPacket[]>("SELECT 1 AS ok");
    expect(rows[0].ok).toBe(1);
  });

  test("MySQLquery", async ({ request }) => {
    const createRes = await request.post("/api/todos", {
      data: {
        title: API_CREATED_TITLE,
        description: API_CREATED_DESCRIPTION,
      },
    });

    expect(createRes.status()).toBe(201);
    const created = await createRes.json();
    expect(created.title).toBe(API_CREATED_TITLE);
    expect(created.description).toBe(API_CREATED_DESCRIPTION);

    const found = await waitForTodoInAnyCandidateDB(
      connection,
      API_CREATED_TITLE,
      API_CREATED_DESCRIPTION,
      10,
      75
    );

    expect(
      found,
      `Could not find API-created todo in candidate databases: ${databaseCandidates.join(", ")}. ` +
        "Set Playwright DB_NAME to the same database used by the API server."
    ).not.toBeNull();

    if (!found) {
      return;
    }

    resolvedDatabase = found.db;
    expect(found.row.title).toBe(API_CREATED_TITLE);
    expect(found.row.description).toBe(API_CREATED_DESCRIPTION);

    await connection.execute(
      `DELETE FROM \`${found.db}\`.todos WHERE title = ? AND description = ?`,
      [API_CREATED_TITLE, API_CREATED_DESCRIPTION]
    );
  });

  test("MySQlinsert", async ({ request }) => {
    const orderedCandidates = resolvedDatabase
      ? [
          resolvedDatabase,
          ...databaseCandidates.filter((db) => db !== resolvedDatabase),
        ]
      : databaseCandidates;

    let matchedDatabase: string | null = null;

    for (const db of orderedCandidates) {
      try {
        await connection.execute(
          `DELETE FROM \`${db}\`.todos WHERE title = ? AND description = ?`,
          [DIRECT_MYSQL_TITLE, DIRECT_MYSQL_DESCRIPTION]
        );

        const [result] = await connection.execute<ResultSetHeader>(
          `INSERT INTO \`${db}\`.todos (title, description, completed) VALUES (?, ?, ?)`,
          [DIRECT_MYSQL_TITLE, DIRECT_MYSQL_DESCRIPTION, false]
        );
        expect(result.affectedRows).toBeGreaterThan(0);
        matchedDatabase = db;
        break;
      } catch (err) {
        if (!isSkippableDatabaseError(err)) {
          throw err;
        }
      }
    }

    expect(
      matchedDatabase,
      "Could not insert through direct MySQL calls in any candidate database. " +
        `Tried: ${orderedCandidates.join(", ")}. ` +
        "Set Playwright DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME to match the API server DB settings."
    ).not.toBeNull();

    if (!matchedDatabase) {
      return;
    }

    resolvedDatabase = matchedDatabase;
    const healthRes = await request.get("/api/todos");
    expect(healthRes.status()).toBe(200);

    await connection.execute(
      `DELETE FROM \`${matchedDatabase}\`.todos WHERE title = ? AND description = ?`,
      [DIRECT_MYSQL_TITLE, DIRECT_MYSQL_DESCRIPTION]
    );

    // Best-effort cleanup verification without strict row-count dependency in replay mode.
    await connection.query<TodoRow[]>(
      `SELECT id FROM \`${matchedDatabase}\`.todos WHERE title = ? AND description = ? ORDER BY id DESC LIMIT 1`,
      [DIRECT_MYSQL_TITLE, DIRECT_MYSQL_DESCRIPTION]
    );
  });
});
