# Task tracker

Hi,

Welcome to the task tracker cli, in this cli you could ADD, UPDATE, DELETE, MARK as IN-PROGRESS or DONE and list your tasks.

### Steps
First, you need to download the project via git clone.
Second, use the following commands below.

ADD
Adding a new task.
./task-cli add "Your task description"

UPDATE
Updating an existing task.
./task-cli update 1 "Updated description"

DELETE
Deleting an existing task.
./task-cli delete 1

MARK
Marking a task as in progress or done.
./task-cli mark-in-progress 1
./task-cli mark-done 1

LIST
Listing all tasks
./task-cli list

Listing tasks by status
./task-cli list done
./task-cli list todo
./task-cli list in-progress
