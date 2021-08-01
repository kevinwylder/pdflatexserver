# pdflatexserver

This program starts a web server that compiles `.tex` files into PDF. There are 
tons of options for interactively using latex, I made this one since I like to use
my own text editor, and think its cool to just reload the page and have my document
recompile

## Usage

```
./pdflatexserver
USAGE: pdflatexserver <flags> <directory to serve>
  -port string
        port (default ":80")
  -template string
        path to template dir
```

After starting the server, visit the your computer's port in a web browser. The 
provided directory arg is indexed. Only directories and `.tex` files are shown.

This supports macos and linux. It uses whatever `pdflatex` you have in your path 
first, and if none is found it searches for a 
[texlive](https://www.tug.org/texlive/) installation (works with 
[MacTex](https://www.tug.org/mactex/) also). 


## Installation

precompiled binaries
* [x86 linux](https://kwylder.com/bin/x86_64-linux/pdflatexserver) 
* [x86 mac](https://kwylder.com/bin/x86_64-darwin/pdflatexserver) (Not compatible with macs M1 or newer)

If you have the golang 1.16 runtime already then you can install with 

```
git clone git@github.com:kevinwylder/pdflatexserver.git
go install ./cmd/pdflatexserver
```

Otherwise, there is a `Dockerfile` that can get you up and running too. I recommend
podman in general because 1. docker runs as root by default and 2. file permissions
are automatically mapped to the right user ID

```
git clone git@github.com:kevinwylder/pdflatexserver.git
podman build ./pdflatexserver -t pdflatexserver
cd /path/to/your/tex
podman run --rm -it -p 80:80 -v $(pwd):/data pdflatexserver /data
```

