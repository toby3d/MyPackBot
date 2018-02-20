# Code standarts
Standards help to keep the code readable and understandable, although they may seem strange or uncomfortable. Code style described below is not strict, but I give priority to those contributors who follow it. :heart:

## Rules
- Format the changes via `go fmt` before committing.
- Indents with tabs (with width 8), no spaces.
- In `import` external packages are separated from native by empty line.
- Maximum line lenght is 120 characters.
- Do not forget to comment what do you do.
- Check what you are writing, tests before commiting.
- Double check what you are writing, remove all [gometalinter](https://github.com/alecthomas/gometalinter) warnings.

## Guidance
Keep in mind the following before, while and after writing the code:
- **Less is always more.**
Write the least amount of code possible to solve just the problem at hand.
- **Predicting the future is impossible.**
Try to distinguish between anticipating potential future problems and potential future features. The former is usually good, the latter is usually bad.
- **Functional programming is functional.**
Functions should be small and single-purpose. Large variable lists are a sign your function does too much.

# Git workflow
`master` contains a stable version of the project, when as `develop` it is constantly updated and contains the latest changes. When proposing changes, you must specify `develop` as the target branch of PR.

## Commits
- First line of commit message should be to 80 chars long as public description of what you have achieved with the commit.
- Leave a blank line after the first line.
- The 3rd line can reference issue with `issue #000` if you just want to mention an issue or `closes #000` if your commit closes an issue. If you don't have an issue to reference or close, think carefully about whether you need to raise one before opening a PR.
- Use bullet points on the following lines to explain what your changes achieve and in particular why you've used your approach.
- Add a contextual emoji before commit title [based on its content](https://gitmoji.carloscuesta.me) (or [use appropriate tool for commiting](https://github.com/carloscuesta/gitmoji-cli)). This **greatly** helps visually to distinguish commits among themselves.

If you need to update your existing commit message, you can do this by running `git commit --amend` on your branch.

## Pull requests
The easier it is for me to merge a PR, the faster we'll be able to do it. Please take steps to make merging easy and keep the history clean and useful.

- **Always work on a branch.**
It will make your life much easier, really. Not touching the `master` branch will also simplify keeping your fork up-to-date.
- **Use issues properly.**
Bugs, changes and features are all different and should be treated differently. Use your commit message to close or reference issues. The more information you provide, the more likely your PR will get merged.

## Issues
Feel free to pick up any issue which is not assigned. Please leave a comment on the issue to say you wish to pick it up, and it will get assigned to you.