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
./main add "Some task"
Task created! Id: ******
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
