# Introduction

This repo is written for friend to manage remote machines with batched ssh. Just personal project. So, Think before forking.

## Usage
```bash
main.exe -h
Usage of main.exe:
  -f string
        指定csv文件名称 (default "demo.csv")


main.exe -f demo.csv
```

The `-f` flag is optional, the default value is `demo.csv` file within the some folder besides the `main.exe`.


## Notes
如果不指定-f选项，默认使用相同目录下的demo.csv文件，如果指定了，则使用指定的文件。