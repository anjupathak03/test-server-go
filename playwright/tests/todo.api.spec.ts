import { test, expect, APIRequestContext } from "@playwright/test";

// ---------------------------------------------------------------------------
// Types mirroring the Go models
// ---------------------------------------------------------------------------
interface Todo {
  id: number;
  title: string;
  description: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

// ---------------------------------------------------------------------------
// Helper to clean up todos created during tests
// ---------------------------------------------------------------------------
async function deleteTodo(request: APIRequestContext, id: number) {
  await request.delete(`/api/todos/${id}`);
}

// ---------------------------------------------------------------------------
// POST /api/todos – Create
// ---------------------------------------------------------------------------
test.describe("POST /api/todos", () => {
  let createdId: number;

  test.afterEach(async ({ request }) => {
    if (createdId) await deleteTodo(request, createdId);
  });

  test("createwithdesc", async ({ request }) => {
    const res = await request.post("/api/todos", {
      data: { title: "Buy groceries", description: "Milk, eggs, bread" },
    });

    expect(res.status()).toBe(201);
    const body: Todo = await res.json();
    expect(body.id).toBeGreaterThan(0);
    expect(body.title).toBe("Buy groceries");
    expect(body.description).toBe("Milk, eggs, bread");
    expect(body.completed).toBe(false);
    createdId = body.id;
  });

  test("titleonly", async ({ request }) => {
    const res = await request.post("/api/todos", {
      data: { title: "Minimal todo" },
    });

    expect(res.status()).toBe(201);
    const body: Todo = await res.json();
    expect(body.title).toBe("Minimal todo");
    expect(body.description).toBe("");
    createdId = body.id;
  });

  test("missingtitle", async ({ request }) => {
    const res = await request.post("/api/todos", {
      data: { description: "No title provided" },
    });

    expect(res.status()).toBe(400);
    const body = await res.json();
    expect(body).toHaveProperty("error");
  });

  test("emptybody", async ({ request }) => {
    const res = await request.post("/api/todos", { data: {} });
    expect(res.status()).toBe(400);
  });
});

// ---------------------------------------------------------------------------
// GET /api/todos – List all
// ---------------------------------------------------------------------------
test.describe("GET /api/todos", () => {
  const ids: number[] = [];

  test.beforeAll(async ({ request }) => {
    for (const title of ["Alpha", "Beta", "Gamma"]) {
      const res = await request.post("/api/todos", { data: { title } });
      const body: Todo = await res.json();
      ids.push(body.id);
    }
  });

  test.afterAll(async ({ request }) => {
    for (const id of ids) await deleteTodo(request, id);
  });

  test("listtodos", async ({ request }) => {
    const res = await request.get("/api/todos");
    expect(res.status()).toBe(200);

    const body: Todo[] = await res.json();
    expect(Array.isArray(body)).toBe(true);
    expect(body.length).toBeGreaterThanOrEqual(3);
  });

  test("todoshape", async ({ request }) => {
    const res = await request.get("/api/todos");
    const todos: Todo[] = await res.json();

    for (const todo of todos) {
      expect(todo).toHaveProperty("id");
      expect(todo).toHaveProperty("title");
      expect(todo).toHaveProperty("description");
      expect(todo).toHaveProperty("completed");
      expect(todo).toHaveProperty("created_at");
      expect(todo).toHaveProperty("updated_at");
    }
  });

  test("seededtodos", async ({ request }) => {
    const res = await request.get("/api/todos");
    const todos: Todo[] = await res.json();
    const titles = todos.map((t) => t.title);

    expect(titles).toContain("Alpha");
    expect(titles).toContain("Beta");
    expect(titles).toContain("Gamma");
  });
});

// ---------------------------------------------------------------------------
// GET /api/todos/:id – Get by ID
// ---------------------------------------------------------------------------
test.describe("GET /api/todos/:id", () => {
  let todoId: number;

  test.beforeAll(async ({ request }) => {
    const res = await request.post("/api/todos", {
      data: { title: "Find me", description: "Specific todo" },
    });
    const body: Todo = await res.json();
    todoId = body.id;
  });

  test.afterAll(async ({ request }) => {
    await deleteTodo(request, todoId);
  });

  test("getbyid", async ({ request }) => {
    const res = await request.get(`/api/todos/${todoId}`);
    expect(res.status()).toBe(200);

    const body: Todo = await res.json();
    expect(body.id).toBe(todoId);
    expect(body.title).toBe("Find me");
    expect(body.description).toBe("Specific todo");
  });

  test("notfound", async ({ request }) => {
    const res = await request.get("/api/todos/99999999");
    expect(res.status()).toBe(404);
    const body = await res.json();
    expect(body).toHaveProperty("error");
  });

  test("invalidid", async ({ request }) => {
    const res = await request.get("/api/todos/not-a-number");
    expect(res.status()).toBe(400);
  });
});

// ---------------------------------------------------------------------------
// PUT /api/todos/:id – Update
// ---------------------------------------------------------------------------
test.describe("PUT /api/todos/:id", () => {
  let todoId: number;

  test.beforeEach(async ({ request }) => {
    const res = await request.post("/api/todos", {
      data: { title: "Original title", description: "Original description" },
    });
    const body: Todo = await res.json();
    todoId = body.id;
  });

  test.afterEach(async ({ request }) => {
    await deleteTodo(request, todoId);
  });

  test("updatetitle", async ({ request }) => {
    const res = await request.put(`/api/todos/${todoId}`, {
      data: { title: "Updated title" },
    });

    expect(res.status()).toBe(200);
    const body: Todo = await res.json();
    expect(body.title).toBe("Updated title");
    expect(body.description).toBe("Original description"); // unchanged
  });

  test("markcompleted", async ({ request }) => {
    const res = await request.put(`/api/todos/${todoId}`, {
      data: { completed: true },
    });

    expect(res.status()).toBe(200);
    const body: Todo = await res.json();
    expect(body.completed).toBe(true);
  });

  test("updateall", async ({ request }) => {
    const res = await request.put(`/api/todos/${todoId}`, {
      data: {
        title: "New title",
        description: "New description",
        completed: true,
      },
    });

    expect(res.status()).toBe(200);
    const body: Todo = await res.json();
    expect(body.title).toBe("New title");
    expect(body.description).toBe("New description");
    expect(body.completed).toBe(true);
  });

  test("updatenotfound", async ({ request }) => {
    const res = await request.put("/api/todos/99999999", {
      data: { title: "Ghost" },
    });
    expect(res.status()).toBe(404);
  });

  test("invalidid", async ({ request }) => {
    const res = await request.put("/api/todos/abc", {
      data: { title: "Bad ID" },
    });
    expect(res.status()).toBe(400);
  });
});

// ---------------------------------------------------------------------------
// DELETE /api/todos/:id – Delete
// ---------------------------------------------------------------------------
test.describe("DELETE /api/todos/:id", () => {
  test("deletestodo", async ({ request }) => {
    const createRes = await request.post("/api/todos", {
      data: { title: "To be deleted" },
    });
    const { id } = await createRes.json();

    const deleteRes = await request.delete(`/api/todos/${id}`);
    expect(deleteRes.status()).toBe(200);

    // Confirm it is gone
    const getRes = await request.get(`/api/todos/${id}`);
    expect(getRes.status()).toBe(404);
  });

  test("deletenotfound", async ({ request }) => {
    const res = await request.delete("/api/todos/99999999");
    expect(res.status()).toBe(404);
  });

  test("invalidid", async ({ request }) => {
    const res = await request.delete("/api/todos/not-an-id");
    expect(res.status()).toBe(400);
  });
});

// ---------------------------------------------------------------------------
// End-to-end CRUD lifecycle
// ---------------------------------------------------------------------------
test.describe("Full CRUD lifecycle", () => {
  test("singletodo", async ({ request }) => {
    // Create
    const createRes = await request.post("/api/todos", {
      data: { title: "Lifecycle test", description: "E2E test" },
    });
    expect(createRes.status()).toBe(201);
    const created: Todo = await createRes.json();
    const id = created.id;

    // Read
    const readRes = await request.get(`/api/todos/${id}`);
    expect(readRes.status()).toBe(200);
    const read: Todo = await readRes.json();
    expect(read.title).toBe("Lifecycle test");

    // Update
    const updateRes = await request.put(`/api/todos/${id}`, {
      data: { title: "Updated lifecycle", completed: true },
    });
    expect(updateRes.status()).toBe(200);
    const updated: Todo = await updateRes.json();
    expect(updated.title).toBe("Updated lifecycle");
    expect(updated.completed).toBe(true);

    // Delete
    const deleteRes = await request.delete(`/api/todos/${id}`);
    expect(deleteRes.status()).toBe(200);

    // Verify deletion
    const afterDelete = await request.get(`/api/todos/${id}`);
    expect(afterDelete.status()).toBe(404);
  });
});
