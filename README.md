# Robofinder

Robofinder is a Golang tool that gathers paths from robots.txt files of archived web pages. It uses the archive.org service to find the archived versions of a 
given domain and extracts the paths from the robots.txt files. It prints the paths without duplicates and with the host name.


## Installation
```bash
go install -v github.com/alirezac0/robofinder@latest
```

## Usage
To use Robofinder, you need to pass a domain name as a command line argument. For example:

```bash
robofinder https://example.com
```

This will print the paths from robots.txt files of archived versions of example.com.


## To Do
•  [ ] Adding other resources

•  [ ] Adding concurrency for faster processing

•  [ ] Adding flags for customizing output and parameters

•  [ ] Adding tests and documentation




