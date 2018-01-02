# roll
Command line dice rolling.  
Do you need this? Probably not! :-)  
Only written because I wanted to try making something in go.

Lets you roll any number of any dice (d20, d10, d6, d5, etc.), in the terminal.  
Also lets you add modifiers (+5, -10, etc.).


# Installation
Place roll.go somewhere in gopath/src, and run `go install`, and run using `roll [expression]` (assumes gopath/bin added to PATH)  
Or, just place roll.go in any folder, and run `go build`, then `./roll [expression]` there.

# Usage
* Input should consist of any number of terms.  
* Terms can be either a constant or dice throw, and they are separated by + or -.  
* If there is no sign before a term, it is assumed to be positive.  
* Dice are written as __adb__, where we throw __a__ dice with __b__ sides. (Use lowercase d!)  
* Constants are any number  
* Negative signs before a term can be used!  

Example inputs:
```bash
me@my-computer:~ roll 3d20 + 5 - 1d4 - 6
Rolling: 3d20 -1d4 5 -6
Throws:	 3d20 	-> 34
	 -1d4 	-> -3
Const:	 -1
Sum: 	 30
```

```bash
me@my-computer:~ roll 2d10 2d9 +2d8 -2d7 - 2d6 +5 - 4
Rolling: 2d10 2d9 2d8 -2d7 -2d6 5 -4
Throws:	 2d10 	-> 9
	 2d9 	-> 17
	 2d8 	-> 11
	 -2d7 	-> -9
	 -2d6 	-> -3
Const:	 1
Sum: 	 26
```





