# Introduction to Programming with Aura

> **Your first steps into the Aura programming language.**
> This tutorial will take you from zero to writing real programs in Aura — no prior Aura experience required!

**Aura version:** v0.9.0
**Difficulty:** Beginner
**Time to complete:** ~45 minutes

---

## Table of Contents

1. [What is Aura?](#what-is-aura)
2. [Getting Started](#getting-started)
3. [Tutorial 1: Hello, World!](#tutorial-1-hello-world)
4. [Tutorial 2: Calculator Program](#tutorial-2-calculator-program)
5. [Tutorial 3: Number Guessing Game](#tutorial-3-number-guessing-game)
6. [Aura Syntax Quick Reference](#aura-syntax-quick-reference)
7. [Next Steps](#next-steps)

---

## What is Aura?

Aura is a **modern, Python-inspired programming language** designed for clarity, safety, and AI-human collaboration. If you've ever written Python, Aura will feel familiar — but with powerful extras like static types, pattern matching, and an effect system that makes your code safer and easier to reason about.

### Why Learn Aura?

- **Clean syntax** — Indentation-based blocks, just like Python. No curly braces or semicolons.
- **Safe by default** — `Option` and `Result` types prevent null pointer errors and unhandled exceptions.
- **Pattern matching** — Elegant `match` expressions replace messy chains of `if`/`elif`.
- **Batteries included** — 17 standard library modules with 117+ functions ready to use.
- **AI-first design** — Built for seamless collaboration between humans and AI agents.

### What You'll Learn

In these three tutorials, you'll progressively learn:

| Tutorial | What You'll Build | Key Concepts |
|---|---|---|
| **Hello World** | A greeting program | Functions, strings, `std.io`, program structure |
| **Calculator** | A math calculator | Pattern matching, arithmetic, functions with parameters |
| **Guessing Game** | A number guessing game | Loops, conditionals, random numbers, mutable variables |

Let's dive in! 🚀

---

## Getting Started

### Installing Aura

Aura is built in Go. To get started:

```bash
# 1. Make sure you have Go 1.22+ installed
go version

# 2. Clone the Aura repository
git clone https://github.com/unclebucklarson/aura.git
cd aura

# 3. Build the Aura CLI
go build -o aura ./cmd/aura

# 4. (Optional) Move it to your PATH so you can use it from anywhere
sudo mv aura /usr/local/bin/
```

### Verify It Works

```bash
./aura repl
```

You should see:

```
Aura REPL v0.3 (Phase 3)
Type expressions or statements. Press Ctrl+D to exit.
```

Type `42 + 8` and press Enter — you should see `50`. Press `Ctrl+D` to exit.

### Creating Your First Aura File

Aura source files use the `.aura` extension. You can create them with any text editor:

```bash
# Create a new file
touch my_program.aura
```

### Running Aura Programs

To run an Aura program, use the `aura run` command:

```bash
aura run my_program.aura
```

This will:
1. Parse your `.aura` file
2. Execute all top-level code
3. Call the `main()` function (if one exists)

### The Aura Toolchain

Aura comes with several useful commands:

| Command | What It Does |
|---|---|
| `aura run file.aura` | Execute a program |
| `aura test file.aura` | Run test blocks |
| `aura format file.aura` | Pretty-print your code |
| `aura check file.aura` | Type-check without running |
| `aura repl` | Interactive playground |

Now you're ready to write your first program!

---

## Tutorial 1: Hello, World!

Every programming journey starts with "Hello, World!" — a simple program that prints a message to the screen. In Aura, it's delightfully simple.

### The Complete Program

Create a file called `hello_world.aura`:

```aura
module hello_world

import std.io

fn main():
    std.io.println("Hello, World!")
    std.io.println("Welcome to Aura! 🎉")
```

### Running It

```bash
aura run hello_world.aura
```

### Expected Output

```
Hello, World!
Welcome to Aura! 🎉
```

🎉 **Congratulations!** You just ran your first Aura program!

### Step-by-Step Explanation

Let's break down every line:

#### Line 1: `module hello_world`

```aura
module hello_world
```

Every Aura file **must** start with a `module` declaration. This gives your file a name — think of it as a label for your code. The module name is typically the same as the filename (without `.aura`).

> 💡 **Key concept:** Modules are how Aura organizes code. Each file is a module.

#### Line 3: `import std.io`

```aura
import std.io
```

This imports the **standard I/O module** — `std.io` — which gives us functions for printing to the screen. In Aura, you need to explicitly import the tools you want to use.

Some commonly used standard library modules:
- `std.io` — Input/output (printing)
- `std.math` — Math functions (sqrt, abs, etc.)
- `std.string` — String utilities
- `std.random` — Random number generation

#### Lines 5–7: The `main()` Function

```aura
fn main():
    std.io.println("Hello, World!")
    std.io.println("Welcome to Aura! 🎉")
```

- `fn` — The keyword that declares a function.
- `main` — The function name. When you use `aura run`, Aura automatically looks for and calls a function named `main()`.
- `():` — Empty parentheses mean this function takes no parameters. The colon `:` starts the function body.
- The indented lines are the **function body**. Aura uses **4-space indentation** to define blocks (just like Python).

#### `std.io.println("Hello, World!")`

- `std.io.println` — Calls the `println` function from the `std.io` module. It prints text followed by a newline.
- `"Hello, World!"` — A **string literal** enclosed in double quotes. Strings in Aura always use double quotes `"`.

### Key Concepts Learned

| Concept | Example | Description |
|---|---|---|
| Module declaration | `module hello_world` | Every file needs one |
| Imports | `import std.io` | Bring in standard library modules |
| Functions | `fn main():` | `fn` keyword, name, params, colon |
| Indentation | 4 spaces | Defines code blocks |
| Strings | `"Hello, World!"` | Double-quoted text |
| Printing | `std.io.println(...)` | Output text to screen |

### 🧪 Exercises

Try these modifications to deepen your understanding:

1. **Change the message:** Print your own name instead of "World".
2. **Add more lines:** Print three lines — your name, your favorite color, and your favorite number.
3. **String interpolation:** Try using Aura's string interpolation:
   ```aura
   let name = "Aura"
   std.io.println("Hello, {name}!")
   ```
   This will print: `Hello, Aura!`

### ⚠️ Common Mistakes

| Mistake | What Happens | Fix |
|---|---|---|
| Missing `module` line | Parse error | Always start with `module name` |
| Using tabs instead of spaces | Parse error | Use exactly 4 spaces for indentation |
| Forgetting to import `std.io` | Runtime error: undefined | Add `import std.io` at the top |
| Using single quotes `'hello'` | Parse error | Always use double quotes `"hello"` |
| Wrong indentation (2 or 3 spaces) | Unexpected behavior | Always use exactly **4 spaces** |

---

## Tutorial 2: Calculator Program

Now let's build something more useful — a calculator! This program demonstrates functions with parameters, arithmetic operations, and Aura's powerful **pattern matching**.

### The Complete Program

Create a file called `calculator.aura`:

```aura
module calculator

import std.io

fn add(a, b):
    return a + b

fn subtract(a, b):
    return a - b

fn multiply(a, b):
    return a * b

fn divide(a, b):
    if b == 0:
        std.io.println("  ⚠️  Error: Cannot divide by zero!")
        return none
    return a / b

fn calculate(a, op, b):
    std.io.println("Calculating: {a} {op} {b}")
    let result = match op:
        "add"      -> add(a, b)
        "subtract" -> subtract(a, b)
        "multiply" -> multiply(a, b)
        "divide"   -> divide(a, b)
        _          -> none
    if result == none:
        if op != "add" and op != "subtract" and op != "multiply" and op != "divide":
            std.io.println("  ⚠️  Unknown operation: {op}")
    else:
        std.io.println("  ✅ Result: {result}")
    return result

fn main():
    std.io.println("=================================")
    std.io.println("   🧮 Aura Calculator")
    std.io.println("=================================")
    std.io.println("")

    # Basic arithmetic
    calculate(10, "add", 5)
    calculate(10, "subtract", 3)
    calculate(10, "multiply", 4)
    calculate(10, "divide", 3)

    std.io.println("")
    std.io.println("----- Edge Cases -----")

    # Edge cases
    calculate(10, "divide", 0)
    calculate(7, "modulo", 3)

    std.io.println("")
    std.io.println("----- Chained Calculations -----")

    # Chained calculations
    let step1 = add(100, 50)
    std.io.println("Step 1: 100 + 50 = {step1}")

    let step2 = multiply(step1, 2)
    std.io.println("Step 2: {step1} * 2 = {step2}")

    let step3 = subtract(step2, 25)
    std.io.println("Step 3: {step2} - 25 = {step3}")

    std.io.println("")
    std.io.println("Final answer: {step3}")
    std.io.println("=================================")
```

### Running It

```bash
aura run calculator.aura
```

### Expected Output

```
=================================
   🧮 Aura Calculator
=================================

Calculating: 10 add 5
  ✅ Result: 15
Calculating: 10 subtract 3
  ✅ Result: 7
Calculating: 10 multiply 4
  ✅ Result: 40
Calculating: 10 divide 3
  ✅ Result: 3

----- Edge Cases -----
Calculating: 10 divide 0
  ⚠️  Error: Cannot divide by zero!
Calculating: 7 modulo 3
  ⚠️  Unknown operation: modulo

----- Chained Calculations -----
Step 1: 100 + 50 = 150
Step 2: 150 * 2 = 300
Step 3: 300 - 25 = 275

Final answer: 275
=================================
```

### Step-by-Step Explanation

#### Arithmetic Functions

```aura
fn add(a, b):
    return a + b

fn subtract(a, b):
    return a - b
```

- Each function takes **two parameters** (`a` and `b`).
- The `return` keyword sends a value back to the caller.
- Aura supports all standard arithmetic operators: `+`, `-`, `*`, `/`, `%` (modulo).

> 💡 **Key concept:** Functions are defined with `fn`, take parameters in parentheses, and use `return` to send back results.

#### Safe Division with Error Handling

```aura
fn divide(a, b):
    if b == 0:
        std.io.println("  ⚠️  Error: Cannot divide by zero!")
        return none
    return a / b
```

- We check `if b == 0` to prevent division by zero.
- `none` is Aura's "no value" — similar to `null` in other languages, but safer because Aura tracks it through the type system.
- The `if` block uses **4-space indentation** for its body.

> 💡 **Key concept:** In Aura, `none` represents the absence of a value. It's part of the `Option` type system that prevents null pointer errors.

#### Pattern Matching with `match`

```aura
let result = match op:
    "add"      -> add(a, b)
    "subtract" -> subtract(a, b)
    "multiply" -> multiply(a, b)
    "divide"   -> divide(a, b)
    _          -> none
```

This is one of Aura's most powerful features — **pattern matching**! Here's how it works:

- `match op:` — Look at the value of `op` and compare it against patterns.
- `"add" -> add(a, b)` — If `op` equals `"add"`, call `add(a, b)` and return the result.
- `_ -> none` — The underscore `_` is a **wildcard** that matches anything. It's like the `default` case in a switch statement.
- The `match` expression returns a value, so we can assign it to `let result`.

> 💡 **Key concept:** `match` is an **expression** — it produces a value. This is cleaner than long `if`/`elif` chains.

#### String Interpolation

```aura
std.io.println("Calculating: {a} {op} {b}")
```

Aura supports **string interpolation** using curly braces `{}` inside double-quoted strings. Any expression inside `{}` is evaluated and converted to a string.

#### Variables with `let`

```aura
let step1 = add(100, 50)
```

- `let` declares an **immutable variable** — once set, it cannot be changed.
- Aura infers the type automatically, so you don't need to write `let step1: Int = ...`.

### Key Concepts Learned

| Concept | Example | Description |
|---|---|---|
| Function parameters | `fn add(a, b):` | Functions can take inputs |
| Return values | `return a + b` | Send results back to the caller |
| `if` conditionals | `if b == 0:` | Execute code based on conditions |
| `none` value | `return none` | Represents "no value" |
| Pattern matching | `match op: "add" -> ...` | Elegant branching based on values |
| Wildcard pattern | `_ -> none` | Matches anything not matched above |
| String interpolation | `"Result: {result}"` | Embed expressions in strings |
| Immutable variables | `let x = 42` | Variables that can't change |

### 🧪 Exercises

1. **Add modulo:** Add a `modulo(a, b)` function and handle `"modulo"` in the `match` expression.
2. **Add power:** Add a `power(a, b)` function. Hint: you can import `std.math` and use `std.math.pow(a, b)`.
3. **Floating point:** Try `calculate(10, "divide", 3)` — does it return `3` or `3.333...`? Try with `10.0` and `3.0` instead.
4. **Chain more:** Create a sequence of 5 calculations where each step uses the previous result.

### ⚠️ Common Mistakes

| Mistake | What Happens | Fix |
|---|---|---|
| Forgetting `return` | Function returns `none` | Always `return` your result |
| Missing wildcard `_` in match | Runtime error if no pattern matches | Always add a `_` catch-all |
| Wrong comparison `=` vs `==` | Parse error | Use `==` for comparison, `=` for assignment |
| Indentation inside match arms | Parse error | Match arms use `->` syntax |

---

## Tutorial 3: Number Guessing Game

Now for something fun — a number guessing game! This tutorial brings together everything you've learned and introduces **loops**, **mutable variables**, **random numbers**, and **comparisons**.

### The Complete Program

Create a file called `guessing_game.aura`:

```aura
module guessing_game

import std.io
import std.random

fn give_hint(guess, secret):
    if guess < secret:
        std.io.println("  📈 Too low! Try a higher number.")
        return false
    elif guess > secret:
        std.io.println("  📉 Too high! Try a lower number.")
        return false
    else:
        return true

fn play_round(secret, guesses):
    std.io.println("")
    std.io.println("--- Round with {guesses.len()} guesses ---")
    let mut attempts = 0
    let mut found = false
    for guess in guesses:
        attempts = attempts + 1
        std.io.println("")
        std.io.println("Guess #{attempts}: {guess}")
        if give_hint(guess, secret):
            found = true
            std.io.println("  🎉 CORRECT! You got it in {attempts} attempt(s)!")
            break
    if found == false:
        std.io.println("")
        std.io.println("  😢 Out of guesses! The number was {secret}.")
    return attempts

fn simulate_game(secret, min_val, max_val):
    std.io.println("🎯 Secret number is between {min_val} and {max_val}")
    std.io.println("   (Psst... the secret is {secret})")

    # Simulate a player making guesses
    let guesses = [50, 25, 75, 37, 62, 42, 55, 48, 52, secret]
    play_round(secret, guesses)

fn play_binary_search(secret, min_val, max_val):
    std.io.println("")
    std.io.println("🤖 Now let's watch an AI play with binary search!")
    std.io.println("   Target: a number between {min_val} and {max_val}")
    std.io.println("")

    let mut low = min_val
    let mut high = max_val
    let mut attempts = 0
    let mut found = false

    while found == false:
        let guess = (low + high) / 2
        attempts = attempts + 1
        std.io.println("  Attempt {attempts}: AI guesses {guess}")

        if guess == secret:
            std.io.println("  🎉 AI found it in {attempts} attempts!")
            found = true
        elif guess < secret:
            std.io.println("    Too low! Searching higher...")
            low = guess + 1
        else:
            std.io.println("    Too high! Searching lower...")
            high = guess - 1

        if attempts > 20:
            std.io.println("  ⚠️ Safety limit reached!")
            found = true

    return attempts

fn rate_performance(attempts, max_val):
    std.io.println("")
    if attempts <= 3:
        std.io.println("⭐⭐⭐ Amazing! Only {attempts} attempts!")
    elif attempts <= 7:
        std.io.println("⭐⭐ Good job! {attempts} attempts.")
    elif attempts <= 10:
        std.io.println("⭐ Not bad! {attempts} attempts.")
    else:
        std.io.println("Keep practicing! {attempts} attempts.")

fn main():
    std.io.println("╔═══════════════════════════════════╗")
    std.io.println("║   🎲 Aura Number Guessing Game    ║")
    std.io.println("╚═══════════════════════════════════╝")
    std.io.println("")

    # Generate a random secret number between 1 and 100
    let secret = std.random.int(1, 100)
    let min_val = 1
    let max_val = 100

    # Play a simulated round
    simulate_game(secret, min_val, max_val)

    # Watch the binary search AI play
    let ai_attempts = play_binary_search(secret, min_val, max_val)
    rate_performance(ai_attempts, max_val)

    std.io.println("")
    std.io.println("Thanks for playing! 🎮")
```

### Running It

```bash
aura run guessing_game.aura
```

### Example Output

(The secret number is random, so your output will differ!)

```
╔═══════════════════════════════════╗
║   🎲 Aura Number Guessing Game    ║
╚═══════════════════════════════════╝

🎯 Secret number is between 1 and 100
   (Psst... the secret is 73)

--- Round with 10 guesses ---

Guess #1: 50
  📈 Too low! Try a higher number.

Guess #2: 25
  📈 Too low! Try a higher number.

Guess #3: 75
  📉 Too high! Try a lower number.

Guess #4: 37
  📈 Too low! Try a higher number.

Guess #5: 62
  📈 Too low! Try a higher number.

Guess #6: 42
  📈 Too low! Try a higher number.

Guess #7: 55
  📈 Too low! Try a higher number.

Guess #8: 48
  📈 Too low! Try a higher number.

Guess #9: 52
  📈 Too low! Try a higher number.

Guess #10: 73
  🎉 CORRECT! You got it in 10 attempt(s)!

🤖 Now let's watch an AI play with binary search!
   Target: a number between 1 and 100

  Attempt 1: AI guesses 50
    Too low! Searching higher...
  Attempt 2: AI guesses 75
    Too high! Searching lower...
  Attempt 3: AI guesses 62
    Too low! Searching higher...
  Attempt 4: AI guesses 68
    Too low! Searching higher...
  Attempt 5: AI guesses 71
    Too low! Searching higher...
  Attempt 6: AI guesses 73
  🎉 AI found it in 6 attempts!

⭐⭐ Good job! 6 attempts.

Thanks for playing! 🎮
```

### Step-by-Step Explanation

#### Importing Multiple Modules

```aura
import std.io
import std.random
```

You can import as many modules as you need, one per line. Here we use:
- `std.io` for printing to the screen
- `std.random` for generating random numbers

#### Generating Random Numbers

```aura
let secret = std.random.int(1, 100)
```

`std.random.int(min, max)` returns a random integer between `min` and `max` (inclusive). Every time you run the program, you'll get a different secret number!

Other useful `std.random` functions:
- `std.random.float()` — Random decimal between 0.0 and 1.0
- `std.random.choice(list)` — Pick a random element from a list
- `std.random.shuffle(list)` — Shuffle a list randomly

#### Mutable Variables with `let mut`

```aura
let mut attempts = 0
let mut found = false
```

- `let` creates an **immutable** variable (cannot change).
- `let mut` creates a **mutable** variable (can be reassigned).
- Use `let` by default; only use `let mut` when you need the value to change.

```aura
attempts = attempts + 1    # ✅ OK — attempts is mutable
```

> 💡 **Key concept:** Aura encourages immutability. Using `let mut` signals "this value will change" — making your code's intent clearer.

#### Loops with `for` and `while`

**`for` loop** — Iterate over a collection:

```aura
for guess in guesses:
    std.io.println("Guess: {guess}")
```

This loops through each element in the `guesses` list. On each iteration, `guess` holds the current element.

**`while` loop** — Repeat while a condition is true:

```aura
while found == false:
    let guess = (low + high) / 2
    # ... check the guess ...
```

The `while` loop keeps running as long as `found == false`. Once `found` becomes `true`, the loop stops.

#### `break` — Exit a Loop Early

```aura
if give_hint(guess, secret):
    found = true
    std.io.println("  🎉 CORRECT!")
    break
```

`break` immediately exits the innermost loop. Use it when you've found what you're looking for.

#### `if` / `elif` / `else` Chains

```aura
if guess < secret:
    std.io.println("  📈 Too low!")
    return false
elif guess > secret:
    std.io.println("  📉 Too high!")
    return false
else:
    return true
```

- `if` — Check the first condition.
- `elif` — Check additional conditions (short for "else if"). You can have as many `elif` as you need.
- `else` — Runs if no previous condition was true.

#### Lists

```aura
let guesses = [50, 25, 75, 37, 62, 42, 55, 48, 52, secret]
```

Lists in Aura use square brackets `[]` and can contain any values. You can even include variables like `secret` — its value will be used.

Useful list methods:
- `list.len()` — Get the length
- `list.first()` — Get the first element (as Option)
- `list.last()` — Get the last element (as Option)
- `list.contains(x)` — Check if `x` is in the list

#### The Binary Search Algorithm

The `play_binary_search` function demonstrates a classic algorithm:

```aura
let mut low = min_val
let mut high = max_val
let guess = (low + high) / 2
```

1. Start with the full range (1–100).
2. Guess the middle value.
3. If too low, search the upper half. If too high, search the lower half.
4. Repeat until found.

This always finds the number in at most 7 attempts for a range of 1–100! (Because log₂(100) ≈ 7)

### Key Concepts Learned

| Concept | Example | Description |
|---|---|---|
| Random numbers | `std.random.int(1, 100)` | Generate random values |
| Mutable variables | `let mut count = 0` | Variables that can change |
| `for` loops | `for x in list:` | Iterate over collections |
| `while` loops | `while cond:` | Repeat while condition is true |
| `break` | `break` | Exit a loop early |
| `if`/`elif`/`else` | `if x < y: ... elif: ... else:` | Conditional branching |
| Lists | `[1, 2, 3]` | Ordered collections |
| Comparison operators | `<`, `>`, `==`, `!=`, `<=`, `>=` | Compare values |
| Boolean values | `true`, `false` | Logical values |

### 🧪 Exercises

1. **Change the range:** Modify the game to use a range of 1–1000. How many attempts does the binary search need now?
2. **Add difficulty levels:** Create functions for `easy_game()` (1–10), `medium_game()` (1–100), and `hard_game()` (1–1000).
3. **Random guesses:** Instead of the fixed `guesses` list, use `std.random.int()` to generate random guesses.
4. **Count comparisons:** Add a counter that tracks how many comparisons (`<`, `>`, `==`) are made during binary search.
5. **Multiple rounds:** Use a `for` loop to play 5 rounds and track the total attempts across all games.

### ⚠️ Common Mistakes

| Mistake | What Happens | Fix |
|---|---|---|
| Forgetting `mut` | Cannot reassign: "variable is immutable" | Use `let mut` for variables you'll change |
| Infinite `while` loop | Program hangs forever | Always ensure the condition can become false; add a safety limit |
| Off-by-one in `std.random.int` | Wrong range | `std.random.int(1, 100)` includes both 1 and 100 |
| Using `=` instead of `==` | Assignment instead of comparison | Use `==` for comparison |
| Forgetting `break` after finding answer | Loop continues unnecessarily | Add `break` when the goal is achieved |

---

## Aura Syntax Quick Reference

Here's a handy cheat sheet of everything covered in these tutorials:

### Variables

```aura
let name = "Aura"           # Immutable — cannot change
let mut counter = 0          # Mutable — can be reassigned
counter = counter + 1        # OK!
```

### Functions

```aura
fn greet(name):
    std.io.println("Hello, {name}!")

fn add(a, b):
    return a + b
```

### Control Flow

```aura
# If / elif / else
if x > 0:
    std.io.println("positive")
elif x == 0:
    std.io.println("zero")
else:
    std.io.println("negative")

# For loop
for item in [1, 2, 3]:
    std.io.println(item)

# While loop
let mut n = 0
while n < 5:
    n = n + 1

# Match expression
let label = match score:
    100     -> "perfect"
    90      -> "excellent"
    _       -> "good"
```

### Data Types

```aura
let x = 42               # Int
let pi = 3.14            # Float
let name = "Aura"        # String
let active = true        # Bool
let nothing = none        # None
let items = [1, 2, 3]    # List
let pair = (1, "hello")  # Tuple
```

### String Interpolation

```aura
let name = "World"
let greeting = "Hello, {name}!"    # "Hello, World!"
let math = "2 + 2 = {2 + 2}"      # "2 + 2 = 4"
```

### Operators

| Operator | Description | Example |
|---|---|---|
| `+` `-` `*` `/` `%` | Arithmetic | `10 + 3` → `13` |
| `==` `!=` | Equality | `x == 5` |
| `<` `>` `<=` `>=` | Comparison | `x < 10` |
| `and` `or` `not` | Logical | `x > 0 and x < 100` |

### Printing

```aura
import std.io

std.io.println("Hello!")           # Print with newline
std.io.print("No newline")        # Print without newline
```

---

## Next Steps

Congratulations on completing these tutorials! 🎉 You now know the fundamentals of Aura programming. Here's where to go next:

### 📖 Read More Documentation

- **[Language Guide](../user_docs/language_guide.md)** — Deep dive into all Aura features: structs, enums, traits, specs, and more.
- **[Method Reference](../user_docs/method_reference.md)** — Complete reference for 108+ built-in methods on String, List, Map, Option, and Result.
- **[Examples](../user_docs/examples.md)** — Real-world code examples showing method chaining, functional programming, and error handling.
- **[Language Reference](../user_docs/language_reference.md)** — Formal reference for types, syntax, and the effect system.

### 🚀 Features to Explore Next

Now that you know the basics, try learning these powerful Aura features:

1. **Structs** — Create your own data types:
   ```aura
   pub struct Player:
       pub name: String
       pub score: Int = 0
   ```

2. **Enums** — Define types with variants:
   ```aura
   pub enum Shape:
       Circle(Float)
       Rectangle(Float, Float)
   ```

3. **Option & Result types** — Safe error handling:
   ```aura
   let maybe_value = list.first()   # Returns Option
   match maybe_value:
       Some(v) -> std.io.println("Found: {v}")
       None    -> std.io.println("List was empty")
   ```

4. **Pipeline operator** — Chain operations elegantly:
   ```aura
   let result = data |> transform |> format |> validate
   ```

5. **List comprehensions** — Create lists concisely:
   ```aura
   let evens = [x * 2 for x in range(10) if x % 2 == 0]
   ```

6. **The Effect System** — Aura's unique approach to managing side effects (file I/O, network, time) with full mockability for testing.

### 🤖 AI-First Development

Aura was designed for AI-human collaboration. Once you're comfortable with the basics, check out the [AI Mission Statement](../AI_MISSION.md) to learn about:
- **Spec blocks** — Structured contracts that AI uses to generate code
- **Effect annotations** — Tell AI exactly what side effects are allowed
- **`satisfies` clauses** — Automatic verification of AI-generated code

### 💬 Get Help

- Explore the **REPL** (`aura repl`) to experiment interactively
- Read the **test files** in the repository for working examples of every feature
- Check the [ROADMAP](../ROADMAP.md) to see what's coming next

---

*Happy coding with Aura! 🌟*
