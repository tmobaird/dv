# TD (Todo CLI)
This is a fun project I am working on to help learn Golang.
Feature ideas and feedback on the Golang code is welcome 🤝.

**Note: These docs have been updated for the upcoming version. If you would like to use 0.0.1, view the instructions [here](https://github.com/tmobaird/td/blob/58f6f37468b30ebfc8db83538deff3c2532b01b7/README.md#installation).

### Installation

TODO

### Usage

TD is a simple todo list manager that supports multiple lists, todo statuses, and AI-based scheduling.
This tool allows users to manage and optimize their time.

```
Usage:
  td [flags]
  td [command]

Available Commands:
  add         Add a new todo to list
  completion  Generate the autocompletion script for the specified shell
  config      Shows the current config.
  context     Get or set the current context
  help        Help about any command
  list        List todos for current context
  mark        Update a todo's current status
  open        Open current todo list for editing
  rank        Moves an existing todo from one position to another.
  remove      Remove a todo from list
  rename      Renames an existing item in the todo list

Flags:
  -h, --help   help for td

Use "td [command] --help" for more information about a command.
```

### Roadmap

- [x] add
- [x] completion
- [x] config
- [x] context
- [x] help
- [x] list
- [x] mark
- [x] open
- [x] rank
- [x] remove
- [x] rename
- [ ] calendar
- [ ] schedule
  - Type of task
  - Multiple contexts/domains (todo groups)
  - Estimated duration
  - Added at, finished at
  - ```
  td schedule              # shows the schedule or generates it if not saved
  td schedule --generate   # refreshes the last schedule
  ```

### LLM Prompt Ideas

Schedule should have an LLM and Calendar integration something like:

```
Here are my tasks and some notes about them for today (ordered based on importance):
- task 1
- task 2
- task 3

Here is my calendar for today (you can assume not accounted for time is free time):
- 9-9:30 AM Standup
- 11:45 - 12:45 PM Meeting
- 12:45 - 1:30 PM Lunch

Generate a schedule for my tasks, when I should do them and in what order.
The logic should try to factor in the following:
- Prioritization (ordering)
- Time necessary (if duration not estimated in work notes, assume anywhere between 30-90 minutes)

Return the results in the following structure:
interface ScheduledTask {
  name: string,
  id: string,
  start_time: datetime,
  end_time: datetime
}[]

A single task can have multiple entries, if the suggested approach is breaking them up over multiple blocks.
The list should be sorted based on time.
```