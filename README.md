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

# Schedule Brainstorming

- Command Ideas:
  - calendar and schedule

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

# Todo

```
type Todo {
  Name string
  Completed bool
  Metadata Metadata
}

type Metadata {
  durationValue: int
  durationUnit: string
  duration: 100s/10m/1h/1b (b = blocks, 1 block = 30 min)
  type: deep/shallow/quick
  tags: <anything>
}
```

When we schedule we use the duration or type to infer duration needed.

duration will be an integer representing minutes.

## How do I schedule tasks for my work?

I have a list of tasks to do
I have a basic understanding of:
- how long each task will take
- what level of focus is required for each
- how urgent each task is

When I schedule, I will use the following rules:
- Try to do deep focus work as early as possible
- Try to break up deep focus work as minimally as possible
- Non deep focus work can be broken up if necessary
- Ideally more urgent work is scheduled first
- I only spend as much time as necessary on a task
- Sometimes I need a buffer of 15-30 minutes upon completing a task

Ordering:
- How do we choose between deep work and ordered work?
  - Deep work in the first 1/3 gets prioritized in first part of day (first 10 blocks)
  - Order is rank based after that

Continuous work:
- Deep work is NEVER broken up
- Non-deep work can be broken up

### Flow

9
10 = meeting
11
12
1
2

windows = [
  9-10(free),
  10-11(meeting),
  11-3(free)
]

FUNCTION find_task(available [], possible_duration int, prioritize_deep_work bool) (int, int, error):
  var selected
  IF prioritize_deep_work:
    for i = 0; i < available.size / 3; i++:
      if available[i] == deep_work && (available[i].duration <= possible_duration || possible_duration == -1):
        SET selected = available[i]
  
  IF selected == nil:
    return available.first
  ELSE:
    return selected


SET tasks = []
SET task_bank = {id:minutesLeft}
OR
SET to_schedule = list of todos in order of rank // remove from list when scheduled
SET start_time = now
SET schedule = []
WHILE to_schedule.length > 0 && start_time > 5:30PM
  IF calendar.length > 0
    SET time_to_meeting = calendar.first.start_time - start_time
    find_task(possible_duration=time_to_meeting)
    // remove task from schedulable pool
  ELSE
    find_task(possible_duration=-1)
    // remove task from schedulable pool
  END

  SET start_time = scheduled_task.end_time
END



FOR window in windows:
  if event then:
    schedule.push event
  else
    var scheduled
    if window in first 10 blocks:
      for i = 0; i < to_schedule.size / 3; i++:
        set scheduled=to_schedule[i]
        if scheduled is deep work:
    if scheduled == nil:
      scheduled = to_schedule[0]
      to_schedule.pop
    if scheduled:
      schedule.push scheduled, start=window.start, end=window.start+scheduled.duration