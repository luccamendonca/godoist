# Godoist

Godoist provides a (really) simple way to add tasks to a project on Todoist.

## Installing

**Dependencies you _may_ need**
```
libgdk-pixbuf2.0-0
libgl-dev
libglib2.0-0
libglx-dev
libgtk-3-dev
libx11-dev
libxcursor-dev
libxi-dev
libxrandr-dev
```

For now, the only way to go is to build from source and place the binary within
your `$PATH`.

```sh
make build
```

You also need to create a config file in your home directory, named `.godoist`.
It must contain your Todoist API Key, in the following format:
```json
{"todoist": {"apiKey": "123ABC"}}
```

After that you should be able to run it :)

## Examples

**Adding a new task**
```
# Add to Inbox (default)
./main add "Some task"
Task created! Id: ******

# Add to a specific project
./main add "Some task" --project "Work"
Task created! Id: ******

# Using GUI mode with project specified (single dialog)
./main add -g --project "Work"
# Opens task dialog labeled "Task name (Project: Work)"

# Using GUI mode without project (prompts for project first)
./main add -g
# Opens "Project name" dialog with "Inbox" pre-filled, then task dialog
```

**Listing pending tasks**
```
./main list "Inbox"
[]godoist.Task{
  {
    Content: "Some task",
    Id: ******,
    ProjectId: ******,
    ParentId: ******,
    Url: "https://todoist.com/app/task/******",
  },
(...)
```

There's also the option to run in GUI-mode, by passing the `-g|--use-gui` flag.

This is useful to reduce the steps necessary to create a task. You can use
something like `sxhkd` to bind a hotkey to your `godoist -g add` command.

**GUI Workflow:**
- `./main add -g` - Single project input (pre-filled with "Inbox"), then task input
- `./main add -g --project "Work"` - Skip project selection, go straight to task input with "Project: Work" label
- Fast and keyboard-friendly - perfect for hotkey bindings
