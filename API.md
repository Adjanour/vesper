# Vesper API Documentation

This document provides detailed information about the Vesper HTTP API endpoints, including request/response formats and examples.

## Base URL

```
http://localhost:8080/api
```

## Table of Contents

- [Health Check](#health-check)
- [Tasks](#tasks)
  - [List All Tasks](#list-all-tasks)
  - [Create Task](#create-task)
  - [Get Task](#get-task)
  - [Update Task](#update-task)
  - [Delete Task](#delete-task)
- [Error Responses](#error-responses)
- [Data Models](#data-models)

---

## Health Check

Check if the API server is running.

### Endpoint

```
GET /api/health
```

### Response

**Status Code:** `200 OK`

```json
{
  "status": "ok"
}
```

### Example (cURL)

```bash
curl http://localhost:8080/api/health
```

---

## Tasks

### List All Tasks

Retrieve all tasks for the current user (currently hardcoded to user_id = "1").

#### Endpoint

```
GET /api/tasks/
```

#### Response

**Status Code:** `200 OK`

```json
{
  "tasks": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Morning Review",
      "start": "2026-02-08T09:00:00Z",
      "end": "2026-02-08T10:00:00Z",
      "user_id": "1",
      "status": "scheduled"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "title": "Team Meeting",
      "start": "2026-02-08T14:00:00Z",
      "end": "2026-02-08T15:00:00Z",
      "user_id": "1",
      "status": "scheduled"
    }
  ]
}
```

#### Example (cURL)

```bash
curl http://localhost:8080/api/tasks/
```

---

### Create Task

Create a new time block task. The API will check for overlapping tasks and return a conflict error if one exists.

#### Endpoint

```
POST /api/tasks/
```

#### Request Headers

```
Content-Type: application/json
```

#### Request Body

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Morning Review",
  "start": "2026-02-08T09:00:00Z",
  "end": "2026-02-08T10:00:00Z",
  "user_id": "1",
  "status": "scheduled"
}
```

#### Response

**Status Code:** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Morning Review",
  "start": "2026-02-08T09:00:00Z",
  "end": "2026-02-08T10:00:00Z",
  "user_id": "1",
  "status": "scheduled"
}
```

#### Example (cURL)

```bash
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Morning Review",
    "start": "2026-02-08T09:00:00Z",
    "end": "2026-02-08T10:00:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'
```

#### Error Responses

- `400 Bad Request` - Invalid JSON format
- `409 Conflict` - Task overlaps with an existing task

---

### Get Task

Retrieve a specific task by ID.

#### Endpoint

```
GET /api/tasks/{id}
```

#### Path Parameters

| Parameter | Type   | Description              |
|-----------|--------|--------------------------|
| id        | string | The unique task ID (UUID)|

#### Response

**Status Code:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Morning Review",
  "start": "2026-02-08T09:00:00Z",
  "end": "2026-02-08T10:00:00Z",
  "user_id": "1",
  "status": "scheduled"
}
```

#### Example (cURL)

```bash
curl http://localhost:8080/api/tasks/550e8400-e29b-41d4-a716-446655440000
```

#### Error Responses

- `404 Not Found` - Task with the specified ID does not exist

---

### Update Task

Update an existing task. The API will check for overlaps with other tasks (excluding the task being updated).

#### Endpoint

```
PUT /api/tasks/{id}
```

#### Path Parameters

| Parameter | Type   | Description              |
|-----------|--------|--------------------------|
| id        | string | The unique task ID (UUID)|

#### Request Headers

```
Content-Type: application/json
```

#### Request Body

```json
{
  "title": "Morning Review (Extended)",
  "start": "2026-02-08T09:00:00Z",
  "end": "2026-02-08T10:30:00Z",
  "user_id": "1",
  "status": "scheduled"
}
```

**Note:** The `id` field in the request body is ignored; the ID from the URL path is used.

#### Response

**Status Code:** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Morning Review (Extended)",
  "start": "2026-02-08T09:00:00Z",
  "end": "2026-02-08T10:30:00Z",
  "user_id": "1",
  "status": "scheduled"
}
```

#### Example (cURL)

```bash
curl -X PUT http://localhost:8080/api/tasks/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Morning Review (Extended)",
    "start": "2026-02-08T09:00:00Z",
    "end": "2026-02-08T10:30:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'
```

#### Error Responses

- `400 Bad Request` - Invalid request format or validation error
- `404 Not Found` - Task with the specified ID does not exist
- `409 Conflict` - Updated task would overlap with another task

---

### Delete Task

Delete a specific task by ID.

#### Endpoint

```
DELETE /api/tasks/{id}
```

#### Path Parameters

| Parameter | Type   | Description              |
|-----------|--------|--------------------------|
| id        | string | The unique task ID (UUID)|

#### Response

**Status Code:** `204 No Content`

No response body.

#### Example (cURL)

```bash
curl -X DELETE http://localhost:8080/api/tasks/550e8400-e29b-41d4-a716-446655440000
```

#### Error Responses

- `404 Not Found` - Task with the specified ID does not exist

---

## Error Responses

All error responses follow a consistent format:

**Status Codes:**
- `400 Bad Request` - Invalid request format or parameters
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., overlapping tasks)
- `500 Internal Server Error` - Server-side error

**Error Response Body:**

Plain text error message.

Example:
```
task not found
```

---

## Data Models

### Task

Represents a time block in the schedule.

| Field   | Type      | Description                                    | Required |
|---------|-----------|------------------------------------------------|----------|
| id      | string    | Unique identifier (UUID format recommended)    | Yes      |
| title   | string    | Task/event title                               | Yes      |
| start   | datetime  | Start time (RFC3339 format)                    | Yes      |
| end     | datetime  | End time (RFC3339 format)                      | Yes      |
| user_id | string    | User identifier                                | Yes      |
| status  | string    | Task status: "scheduled", "deleted", "replaced"| Yes      |

**Time Format:** ISO 8601 / RFC3339  
Example: `2026-02-08T09:00:00Z`

**Valid Status Values:**
- `scheduled` - Active task
- `deleted` - Soft-deleted task
- `replaced` - Task replaced by another

---

## Complete Example Workflow

Here's a complete example of creating and managing tasks:

```bash
# 1. Check API health
curl http://localhost:8080/api/health

# 2. Create a morning task
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "task-001",
    "title": "Team Standup",
    "start": "2026-02-08T09:00:00Z",
    "end": "2026-02-08T09:30:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'

# 3. Create an afternoon task
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "task-002",
    "title": "Code Review",
    "start": "2026-02-08T14:00:00Z",
    "end": "2026-02-08T15:00:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'

# 4. Get a specific task
curl http://localhost:8080/api/tasks/task-001

# 5. Try to create overlapping task (will fail with 409)
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "task-003",
    "title": "Overlapping Meeting",
    "start": "2026-02-08T09:15:00Z",
    "end": "2026-02-08T09:45:00Z",
    "user_id": "1",
    "status": "scheduled"
  }'

# 6. Delete a task
curl -X DELETE http://localhost:8080/api/tasks/task-001
```

---

## Notes

- All timestamps should be in UTC and follow RFC3339 format
- The API currently uses a hardcoded `user_id = "1"` for all operations
- Multi-user authentication is planned for future releases
- CORS is enabled for all origins (`*`) to support browser-based clients
