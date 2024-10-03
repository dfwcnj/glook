
glook is like the unix/linux look command in that it performs a binary
search of a sorted, newline delimited text file. It differs in that it
can also search sorted non-delimited fixed length record. text files.

Usage of ./glook:<br/>
  -file string<br/>
    	name of sorted file to search (default "/usr/share/dict/words")<br/>
  -fold<br/>
    	fold case<br/>
  -key string<br/>
    	search key<br/>
  -rlen int<br/>
    	fixed length record length<br/>

