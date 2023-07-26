# Todo CLI

### Task Ideas

- [x] Add tests for commands
- [x] Add reporter interface to print output to console
- Flags (https://pkg.go.dev/flag@go1.20.6)
    - [x] Support verbose list
    - [x] Help command
    - [ ] Test flags interface
- [x] Support marking as completed
- [x] Support marking as not completed
- [x] Support editing of todos
- [x] Support batch delete todos
- [x] Support batch add todos
- [x] Perform actions on todos by index or uuid (edit, delete, done, undo)
    - [x] delete
    - [x] edit
    - [x] done
    - [x] undo
- [x] Command shorthands (a for add, l for list, d for delete)
- [ ] Prioritize items in list
- [ ] Change persistence directory to ~/.td
- Persistent configurations
    - [ ] Dont show completed by default
    - [ ] Todo groups
- [ ] Subtodos
