# DV (Dev CLI)

The goal of this project is to create a CLI that can control basic devtools that I am using in my development workflow.
Right now this includes a simple todo list and work log.
The distribution of this is a CLI so that it can really easily integrate into my development environment: VSCode + Terminal

This is a fun project I am working on to help learn Golang.
Feature ideas and feedback on the Golang code is welcome 🤝.

**Note: These docs have been updated for the upcoming version. If you would like to use 0.0.1, view the instructions [here](https://github.com/tmobaird/td/blob/58f6f37468b30ebfc8db83538deff3c2532b01b7/README.md#installation).**

### Installation

TODO

### Modules

This project consists of two modules for end usage:
- TD: Local todo list manager
- LG: Local developer logs tool

The usage for each is below.

### TD Usage

TD is a simple todo list manager that supports multiple lists, todo statuses, and AI-based scheduling.
This tool allows users to manage and optimize their time.

```
Usage:
  dv td [flags]
  dv td [command]

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

Use "dv td [command] --help" for more information about a command.
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
- [x] calendar
- [x] schedule
  - Type of task
  - Multiple contexts/domains (todo groups)
  - Estimated duration
  - Added at, finished at
  - ```
  td schedule              # shows the schedule or generates it if not saved
  td schedule --generate   # refreshes the last schedule
  ```
- [ ] LG Enhancements
  - [x] write
  - [x] show
  - [ ] log
  - [ ] template files

### Todo Metadata

To make scheduling easier td supports optional metadata associated with a todo.
The metadata is a key/value pair format, with an allowed list of keys.
The following metadata keys are allowed:
- duration (a number + interval indicator to express time. Example: `10m` == 10 minutes. Supported units: `m, b, h` - m = minutes, b = blocks (1 block = 30 min), and h = hours) (default = 1b = 30 minutes)
- type (Supported: `deep, shallow, quick`. Only deep has an impact on scheduling today.)

To add these pieces of metadata, simply encode them as a comma separated list where the key/value pair is separated by `=`.
For example:

```
- [ ] Do Homework (duration=10m,type=deep)
```

Will be parsed to have a duration of 10 minutes and a type of deep.
Both fields, and the metadata entirely can be left blank if desired.

### Scheduling

To run the scheduler use the following:

```
dv td schedule
```

The scheduler has 2 potential inputs: the todo list in your current context and your calendar.
Every todo has a priority (the order in the list) and a duration (see #todo-metadata for details about this).
The scheduler will follow a few rules when scheduling items:
1. The day of scheduling starts when the `schedule` command is run and ends at 6:30 PM (TODO: make this configurable).
2. Calendar events are always honored or seen as must schedule at the exact time.
  - This means that if you have an invite for 10-10:30 AM, the scheduler will never schedule anything else in the 10-10:30 AM window.
3. Deep work is prioritized during the first 1/3 of the day. During the beginning parts of our day we have more energy that can sometimes be required for deeper level tasks. Therefore in the first 1/3 of the day we will prioritize this deep work. We only use 1/3 of the day for this to leave time for the majority of potentially higher priority todos.
4. Outside of #3, todos are scheduled in priority order of where they appear in the todo list.
5. Scheduled todos that are broken up by meetings will have their duration be split as well.

**The only current supported calendar is Gmail. This will be improved soon.**