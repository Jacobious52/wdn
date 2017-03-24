# wdn
Watch Do Notify (macOS)

Simple program to watch a set of files for changes, then execute a bash script and then push an macOS notification with the result

Usage of wdn:
``` 
  -batch
    	bool: should run save as batch job?
  -cmd string
    	string: shell script to run (default "echo hello")
  -log string
    	string: logging output file (default "/dev/null")
  -name string
    	string: name of command (default "Saved")
  -notify
    	bool: should push macOS notification?
  ```
  
 Example:
`wdn -name RVSIM -cmd "./test.sh | tail -n 1 | sed $'s,\[[0-9;]*[a-zA-Z],,g'" -notify *.h *.cpp`
  
