# optimg

`optimg` is a Go based CLI tool which resizes and compresses images for web development. It is intended to be run periodically on the static directory of a web project. 

## Usage
```
optimg 1.0.0
Usage: optimg [--input INPUT] [--steps STEPS] [--quality QUALITY] [--output OUTPUT] [--type TYPE] [--recurse] [--clear] [--strip] INPUT [INPUT ...]

Positional arguments:
  INPUT                  input directories or files

Options:
  --input INPUT, -i INPUT [default: [.webp .png .jpg .jpeg]]
  --steps STEPS, -s STEPS
                         resizing steps [default: 4]
  --quality QUALITY, -q QUALITY
                         output quality [default: 80]
  --output OUTPUT, -o OUTPUT
                         directory to output to relative to input [default: processed]
  --type TYPE, -t TYPE [default: .webp]
  --recurse, -r          recurse through directories
  --clear, -c            delete output directory before writing
  --strip, -s            strip file metadata
  --help, -h             display this help and exit
  --version              display version and exit
```

### Examples
Resize all `.png` files recursively in a given directory, clearing the previous output, and outputting them as quality 80 `.webp` files:

`optimg static -i .png -t .webp -q 80`

Create 16 resized `.png` files of quality 100 in the directory `output`:

`optimg static -t .webp -q 100 -o output`

## Installation
The recommended installation method is to run:

`go install github.com/emmalexandria/optimg@latest`

This requires that Go is installed.

## Long-term goals

- Implement more detailed and beautiful output
- Move to another argument parser to allow for chaining of flags
- Ensure paths function correctly on Windows


