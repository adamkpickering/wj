# wj

`wj` is a file format and command line tool for keeping a personal **w**ork **j**ournal.
The .wj file format is text-based and easy to read and write using any text editor. The
`wj` tool condenses useful metadata from .wj files.

`wj` follows a similar philosophy as [`clsr`](https://github.com/adamkpickering/clsr).
Like `clsr`, `wj` was inspired by [`ledger`](https://github.com/ledger/ledger), and
more generally by the [plain text accounting](https://plaintextaccounting.org/) ecosystem.


## Should I use `wj`?

`wj` will work well for you if:

- You do all of your work on a computer that you always have access to
- You want to store your journal in version control
- You want easy access to metadata about your journal for e.g. scripting
- You want your journal to be easily searchable using tools such as `grep`
- You work over SSH, or want the option to do so

You should not use `wj` if:

- You are not comfortable with the command line
- You need to include anything other than text (i.e. pictures, sounds)
  in your journal


## How does `wj` work?

You create a directory that holds .wj files (you can create a file
using `wj new`). Each file corresponds to a journal entry for a work
day, and is referred to as an **entry**. More than anything, `wj` is
about the format of these files. Here is an example of such a file:

```
November 24, 2023

To Do
- Write docs for Linux release
- Fill out annual review
- Get back to Justin regarding database migration

Done
- Look at Mark's feedback on graceful cancellation PR 6030
- Refactor WSL distro unregistering code
- Upgrade development environment on Windows machine
- Answer Isabela's questions
- Write script to check database for integrity

10:26 meta Gather info

10:49 work Look at Mark's feedback on graceful cancellation PR 6030

Is TaskRunner actually needed?
- We are going to need something like it in the future as part
  of #5555. Of course, that doesn't mean that we have to
  put it in now, but it's easier to build off of it if it's there.
- Running contextIsDone(ctx) in between each operation would
  work, but would also be more verbose.

12:00 meeting Standup

- Finished dealing with vacation update
- Compliance trainings
- Working on PR feedback
- Dentist this afternoon

12:20 break Lunch

12:40 personal Dentist appointment

14:05 work Look at Mark's feedback on graceful cancellation 6030

14:32 work Refactor WSL distro unregistering code

Is there an issue for this?
- No. Created issue #4321.

15:50 work Upgrade development environment on Windows machine

16:33 meeting Answer Isabela's questions

16:59 work Write script to check database for integrity

17:33 meta Recap

Have to remember to do X thing when Y happens. It is helpful when
I do Z thing.
```

You can think of an entry as having three parts:
- A To Do list. This is where you track tasks that need to be done.
- A Done list. This is where you move tasks that you did during the
  day. It isn't really necessary, but it can feel good to move tasks
  here, and feeling good is an important part of productivity!
- One or more tasks.

A **task** consists of a heading and a body. A task heading has three parts:
- The time at which the task was started. This, along with the start
  time of the subsequent task, is used for time tracking.
- One or more comma-separated tags, which can be used to filter tasks.
- A title that describes the task.

The body contains anything you want: stream-of-consciousness writing
when brainstorming/solving a problem, preparation for a meeting,
meeting notes, the reasoning behind a decision, et cetera.

As you can see, `wj` is three things:
- a to do list
- time tracking software
- a place to record thoughts and information

Once you have one or more entries, you can use the `wj` command to parse
those entries and do useful things. For example, you might use
`wj list tasks --tag work --last 7d` to quickly get a list of work-related
things you have done in the last 7 days.


## Installation

```
go install github.com/adamkpickering/wj@latest
```


## Credits

Thanks to SUSE for holding [Hack Week](https://hackweek.opensuse.org/) 23,
which helped to polish `wj`!
