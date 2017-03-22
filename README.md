# wdn
Watch Do Notify (macOS)

Simple program to watch a set of files for changes, then execute a bash script and then push an macOS notification with the result

usage:
  `wdn <NAME> <Bash script> {files}`

example usage:
  `wdn Tester "./test.sh | tail -n 1 | sed $'s,\[[0-9;]*[a-zA-Z],,g'" *.h *.cpp`
  
