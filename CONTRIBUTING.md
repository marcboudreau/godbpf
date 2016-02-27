# godbpf

This document outlines the process followed for contributing to this project.

## Issues

Defects and enhancement requests are tracked using GitHub issues. If you find an issue with the code, create an issue to let the project team know; or better yet, if you have a suggestion on how to fix it, create a pull request.

## Pull Requests

This is an open source project.  PR are always welcomed.

Here are some guidelines to follow when creating a PR against this project:
- Issue PR early
- Use tasks to track work-in-progress
- Squash multiple commits into one prior to requesting final review
- Don't use the Update Branch button

### Issue PR early

A pull request is a great collaboration tool between the contributor and the project team.  Starting a PR early on to
show the project team what you have in mind allows a discussion early on that will help build a change that has the
highest chance of being accepted.  We don't want to have to turn down work because it diverged too much from the project's
objectives, especially if it took a substantial amount of time and energy.

The convention used in this project is to prefix the PR summary with `[WIP]` until the author feels that
the PR is final and ready for merging.

### Use Tasks to Track Work-in-Progress

Using tasks lists, like the one below, help communicate the progress you've done on your PR and what you expect is 
left to be done.

Example of a Task List:
- [x] First completed item
- [x] Second completed item
- [ ] Uncompleted item

### Squash Multiple Commits

Squashing multiple commits into a single commit helps keep the git history clean and easy to follow.
It is perfectly fine to defer squashing commits until the PR is no longer a Work-in-Progress.

### Don't Use the Update Branch Button

In the event that a commit is made to the master branch after your branch, GitHub may offer to update the source branch by enabling the **Update Branch** button.  The result of clicking that button is that the commits that are not part of the source branch are added through a merge commit.  This is undesired because it muddies the git history and because it tends to complicate the process of squashing multiple commits.  Instead, whenever the branch is out-of-date, it should be rebased using the `git rebase` command.

## Code Quality

Code submitted should be as idiomatic as possible and should be tested by unit tests.
Code changes that improve the readability are always welcomed, even if it doesn't fix any
defects or change any behaviour.

