# go-todo-cli
A to-do list tool for terminal, you can add tasks, remvoe, check as done/inprogress, filter list, sorting
--
## Installing
from releases you can download latest version "go-todo-cli.exe"

## Commands

### add new task
---
`add` command adds task content (head) and priority that optional with default value of routine

```bash
add priority <priority value> <task>
```
priority value can be one of the next:
- urgent
- important
- routine (default value)
- luxury
---
if your task priority value is routine so the command could be like follwing
```bash
add <task>
```

### show to-do list
---
`show` command presents your list in a table of columns ID, task, status (to-do, inprogress, or done), priority, age (time sence adding the task to the list), createdat (date and time for task adding), completedat (when task be chekced as done this value adding with checking date and time)

```bash
show
```
result could be like following table
| ID | TITLE | STATUS | PRIORITY | AGE | CREATED AT | COMPLTED AT |
|----| ----- | ------ | -------- | --- | ---------- | ----------- |
| 1 | do math homework | to-do | urgent | 2 hours ago | 2025-12-21 21:15:23 | |
| 2 | collect 10000 trophy in clash royale | in-progress | important | 20 days ago | 2025-12-01 23:35:46 | |

```bash
show <filter>
```
filter can be status or priority so you can use (to-do, in-progress, done) as status filter, (urgent, important, routine, luxury) as priority filter. You can use only one filter
### remove task
---
`remove` command removes task from the list by its ID (using show command to get IDs, it's auto increment)

```bash
remove <ID>
```

### start doing task
---
`start` command changes status from to-do to in-progress by its ID
```bash
start <ID>
```

### check task as done
---
`check` command changes status from to-do or in-progress to done by its ID
```bash
check <ID>
```

### change priority
---
`priority` command changes priority between (urgent, important, routine, luxury)by its ID
```bash
priority <ID> <new priority value>
```
### clear list
---
`clear` command remove all tasks from the list and be totally empty
```bash
clear
```

### help
---
`help` command print helping message
```bash
help
```