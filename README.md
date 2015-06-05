# wpl-to-m3u
A Go script to convert Windows Media Player Playlist files (.wpl) to M3U files.


Requirements:
This is a Go script, so you must have Go installed in order to run it. If you don't, visit https://golang.org/doc/install for more details.

Usage:
go run wpl-to-m3u.go <fully qualified path to .wpl file>
go run wpl-to-m3u.go <fully qualified path to a directory containing .wpl files>

If a single .wpl file is passed in as the argument, it will make a corresponding .m3u file alongside that file. 
If a directory is passed into the script as its argument, it will convert any .wpl files found in that directory and create corresponding .m3u files along side each one of them.
