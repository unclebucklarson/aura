# Aura Examples

Complete working examples demonstrating Aura's language features. These examples correspond to the test files in the `testdata/` directory.

---

## Table of Contents

1. [Minimal Program](#minimal-program)
2. [Data Models](#data-models)
3. [Specifications](#specifications)
4. [Service Functions](#service-functions)
5. [Control Flow](#control-flow)
6. [Expressions](#expressions)
7. [Complete Application](#complete-application)

---

## Minimal Program

The simplest valid Aura program. Demonstrates basic types, structs, enums, and functions.

```aura
module simple

type Name = String

pub struct Point:
    pub x: Float
    pub y: Float

pub enum Color:
    Red
    Green
    Blue

pub fn add(a: Int, b: Int) -> Int:
    return a + b

let result = add(1, 2)
```

**Key features shown:**
- Module declaration
- Type alias
- Struct with public fields
- Simple enum (no data)
- Function with return type
- Let binding with function call

---

## Data Models

Defining data types for a task management application with refinement types, default values, and enums with data.

```aura
module auratask.models

import std.time as time

## Core type aliases with refinement constraints.

type TaskId = String where len >= 1 and len <= 64
type Priority = Int where self >= 1 and self <= 5
type TaskStatus = "pending" | "in_progress" | "done" | "cancelled"

## The primary Task data model.
pub struct Task:
    pub id: TaskId
    pub title: String where len >= 1 and len <= 200
    pub description: String = ""
    pub status: TaskStatus = "pending"
    pub priority: Priority = 3
    pub created_at: time.Instant
    pub completed_at: time.Instant? = none
    pub tags: [String] = []

## Errors that can occur during task operations.
pub enum TaskError:
    NotFound(TaskId)
    InvalidTitle(String)
    AlreadyCompleted(TaskId)
    Unauthorized(String)

pub enum TaskEvent:
    Created(Task)
    StatusChanged(TaskId, TaskStatus, TaskStatus)
    Deleted(TaskId)
```

**Key features shown:**
- Imports with aliases
- Doc comments (`##`)
- Refinement types (`where len >= 1 and len <= 64`)
- Union string literal types (`"pending" | "done" | ...`)
- Struct fields with defaults and optional types (`time.Instant? = none`)
- Enums with data variants

---

## Specifications

Writing specs that capture function intent before implementation.

```aura
module auratask.specs

from auratask.models import Task, TaskId, TaskError, Priority, TaskStatus

spec CreateNewTask:
    doc: "Creates a new task with the given title and optional priority."

    inputs:
        title: String where len >= 1 and len <= 200 - "The task title"
        priority: Priority = 3 - "Task priority, 1 (low) to 5 (critical)"

    guarantees:
        - "Returns a Task with status 'pending'"
        - "The returned Task has a unique, non-empty id"
        - "The returned Task's created_at is the current time"
        - "The task is persisted in the database"

    effects: db, time

    errors:
        InvalidTitle(String) - "When title is empty or exceeds 200 characters"

spec CompleteTask:
    doc: "Marks an existing task as done and records the completion time."

    inputs:
        task_id: TaskId - "The ID of the task to complete"

    guarantees:
        - "The task's status is changed to 'done'"
        - "The task's completed_at is set to the current time"
        - "The updated task is persisted in the database"

    effects: db, time

    errors:
        NotFound(TaskId) - "When no task exists with the given ID"
        AlreadyCompleted(TaskId) - "When the task is already in 'done' status"

spec ListTasksByStatus:
    doc: "Retrieves all tasks matching the given status, ordered by priority descending."

    inputs:
        status: TaskStatus - "The status to filter by"

    guarantees:
        - "All returned tasks have the requested status"
        - "Results are ordered by priority descending (5 first, 1 last)"

    effects: db
```

**Key features shown:**
- Multiple spec blocks with all sections (`doc`, `inputs`, `guarantees`, `effects`, `errors`)
- Typed inputs with descriptions
- Default values in spec inputs
- Effect declarations
- Error variant descriptions

---

## Service Functions

Implementing functions that satisfy specifications, with effects and error handling.

```aura
module auratask.service

from auratask.models import Task, TaskId, TaskError, Priority, TaskStatus
import std.time as time

pub fn create_task(title: String, priority: Priority = 3) -> Result[Task, TaskError] with db, time satisfies CreateNewTask:
    if title.len == 0 or title.len > 200:
        return Err(TaskError.InvalidTitle("Title must be between 1 and 200 characters, got {title.len}"))

    let now = time.now()
    let id = generate_id()

    let task = Task(
        id: id,
        title: title,
        priority: priority,
        status: "pending",
        created_at: now,
    )

    db.insert("tasks", task)
    return Ok(task)

pub fn complete_task(task_id: TaskId) -> Result[Task, TaskError] with db, time satisfies CompleteTask:
    let maybe_task = db.query_one("tasks", task_id)
    if maybe_task is none:
        return Err(TaskError.NotFound(task_id))

    let task = maybe_task!

    if task.status == "done":
        return Err(TaskError.AlreadyCompleted(task_id))

    let now = time.now()
    let updated = Task(
        id: task.id,
        title: task.title,
        description: task.description,
        status: "done",
        priority: task.priority,
        created_at: task.created_at,
        completed_at: now,
        tags: task.tags,
    )

    db.update("tasks", task_id, updated)
    return Ok(updated)

pub fn list_tasks_by_status(status: TaskStatus) -> Result[[Task], TaskError] with db satisfies ListTasksByStatus:
    let tasks = db.query("tasks", status: status, order_by: "priority", descending: true)
    return Ok(tasks)
```

**Key features shown:**
- `satisfies` clause binding functions to specs
- Effect annotations (`with db, time`)
- `Result` return types with `Ok()` and `Err()`
- Struct construction with named fields
- String interpolation in error messages
- Force unwrap (`!`) operator
- `is` operator for none-checking

---

## Control Flow

Demonstrating all control flow constructs.

```aura
module control

# If / elif / else
pub fn classify_priority(p: Int) -> String:
    if p >= 4:
        return "critical"
    elif p >= 2:
        return "normal"
    else:
        return "low"

# Match with patterns
pub fn describe_color(c: Color) -> String:
    match c:
        case Color.Red:
            return "The color of fire"
        case Color.Green:
            return "The color of nature"
        case Color.Blue:
            return "The color of sky"
        case _:
            return "Unknown color"

# Match with destructuring
pub fn handle_result(result: Result[Task, TaskError]) -> String:
    match result:
        case Ok(task):
            return "Got task: {task.title}"
        case Err(TaskError.NotFound(id)):
            return "Not found: {id}"
        case Err(TaskError.InvalidTitle(msg)):
            return "Invalid: {msg}"
        case Err(e):
            return "Error: {e}"

# For loop
pub fn count_pending(tasks: [Task]) -> Int:
    let mut count = 0
    for task in tasks:
        if task.status == "pending":
            count = count + 1
    return count

# While loop
pub fn countdown(n: Int) -> [Int]:
    let mut result: [Int] = []
    let mut i = n
    while i > 0:
        result.push(i)
        i = i - 1
    return result

# Match with guards
pub fn categorize(value: Int) -> String:
    match value:
        case v if v < 0:
            return "negative"
        case 0:
            return "zero"
        case v if v <= 100:
            return "small"
        case _:
            return "large"
```

**Key features shown:**
- `if` / `elif` / `else` chains
- `match` with constructor patterns
- `match` with destructuring
- `match` with guard clauses
- `for ... in` loops
- `while` loops
- Mutable variables (`let mut`)

---

## Expressions

Showcasing Aura's expression features.

```aura
module expressions

import std.time as time

# Pipeline operator
pub fn process_tasks(tasks: [Task]) -> [String]:
    return tasks
        |> filter_active
        |> sort_by_priority
        |> extract_titles

# List comprehension
pub fn get_urgent_titles(tasks: [Task]) -> [String]:
    return [t.title for t in tasks if t.priority >= 4]

# List comprehension with transform
pub fn double_positives(numbers: [Int]) -> [Int]:
    return [x * 2 for x in numbers if x > 0]

# Lambda expressions
pub fn sort_tasks(tasks: [Task]) -> [Task]:
    return tasks.sort_by(|a, b| -> a.priority > b.priority)

# If expressions (ternary)
pub fn priority_label(p: Int) -> String:
    return if p >= 4 then "urgent" else "normal"

# String interpolation
pub fn format_task(task: Task) -> String:
    return "[{task.priority}] {task.title} ({task.status})"

# Option chaining
pub fn get_display_name(user: User?) -> String:
    let name = user?.profile?.display_name
    return name or "Anonymous"

# Struct construction
pub fn new_task(title: String) -> Task with time:
    return Task(
        id: generate_id(),
        title: title,
        priority: 3,
        status: "pending",
        created_at: time.now(),
    )
```

**Key features shown:**
- Pipeline operator (`|>`)
- List comprehensions with filtering
- Lambda expressions
- If expressions (ternary style)
- String interpolation
- Option chaining (`?.`)
- Struct construction with named fields

---

## Complete Application

Putting it all together — a complete task management module with models, specs, implementation, and tests.

```aura
module taskapp

import std.time as time
import std.uuid as uuid

# ─── Types ───

type TaskId = String where len >= 1
type Priority = Int where self >= 1 and self <= 5

pub struct Task:
    pub id: TaskId
    pub title: String where len >= 1
    pub status: String = "pending"
    pub priority: Priority = 3
    pub created_at: time.Instant
    pub tags: [String] = []

pub enum TaskError:
    NotFound(TaskId)
    InvalidTitle(String)

# ─── Spec ───

spec CreateTask:
    doc: "Creates and persists a new task."

    inputs:
        title: String where len >= 1 - "Task title"
        priority: Priority = 3 - "Priority level"

    guarantees:
        - "Returns a task with status 'pending'"
        - "Task is persisted to the database"

    effects: db, time

    errors:
        InvalidTitle(String) - "When title is empty"

# ─── Implementation ───

pub fn create_task(title: String, priority: Priority = 3) -> Result[Task, TaskError] with db, time satisfies CreateTask:
    if title.len == 0:
        return Err(TaskError.InvalidTitle("Title cannot be empty"))

    let task = Task(
        id: uuid.v4(),
        title: title,
        priority: priority,
        status: "pending",
        created_at: time.now(),
    )

    db.insert("tasks", task)
    return Ok(task)

pub fn find_task(id: TaskId) -> Result[Task, TaskError] with db:
    let maybe = db.query_one("tasks", id)
    if maybe is none:
        return Err(TaskError.NotFound(id))
    return Ok(maybe!)

pub fn list_by_priority(min_priority: Priority) -> [Task] with db:
    let all_tasks = db.query("tasks")
    return [t for t in all_tasks if t.priority >= min_priority]

# ─── Helpers ───

fn format_task(task: Task) -> String:
    return "[P{task.priority}] {task.title} ({task.status})"

fn is_urgent(task: Task) -> Bool:
    return task.priority >= 4

# ─── Tests ───

test "create_task returns pending task":
    with mock_db(), mock_time():
        let result = create_task("Buy groceries")
        assert result.is_ok()
        let task = result.unwrap()
        assert task.title == "Buy groceries"
        assert task.status == "pending"
        assert task.priority == 3

test "create_task rejects empty title":
    with mock_db(), mock_time():
        let result = create_task("")
        assert result.is_err()
        match result:
            case Err(TaskError.InvalidTitle(msg)):
                assert "empty" in msg
            case _:
                assert false, "Expected InvalidTitle"

test "find_task returns not found for missing id":
    with mock_db():
        let result = find_task("nonexistent")
        assert result.is_err()

test "list_by_priority filters correctly":
    with mock_db(), mock_time():
        let _ = create_task("Low", priority: 1)
        let _ = create_task("High", priority: 5)
        let _ = create_task("Medium", priority: 3)

        let urgent = list_by_priority(4)
        assert urgent.len == 1
        assert urgent[0].priority == 5

test "format_task produces expected output":
    with mock_time():
        let task = Task(
            id: "t-1",
            title: "Test",
            priority: 4,
            status: "pending",
            created_at: time.now(),
        )
        assert format_task(task) == "[P4] Test (pending)"

test "is_urgent returns true for priority 4 and above":
    with mock_time():
        let task = Task(id: "t-1", title: "Test", priority: 4, created_at: time.now())
        assert is_urgent(task) == true

        let low = Task(id: "t-2", title: "Test", priority: 2, created_at: time.now())
        assert is_urgent(low) == false
```

**This example demonstrates the full Aura workflow:**
1. Define types with constraints
2. Write specs capturing intent
3. Implement functions that satisfy specs
4. Write tests to verify behavior
5. Use effects to track side effects
6. Use pattern matching for error handling

---

## Running Examples

You can format and parse any of the example files:

```bash
# Format an example
./aura format testdata/models.aura

# Parse and inspect the AST
./aura parse testdata/specs.aura

# Format all testdata files
for f in testdata/*.aura; do ./aura format "$f"; done
```

---

*See the [Language Guide](language_guide.md) for a tutorial-style introduction, or the [Language Reference](language_reference.md) for formal specifications.*
