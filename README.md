# BunnyFinder

BunnyFinder is a framework for finding incentive flaws in Ethereum PoS with little manual effort. BunnyFinder exploits the idea of failure injection, a technique commonly used in
software testing for finding implementation vulnerabilities. Instead of finding implementation
vulnerabilities, BunnyFinder aim to find design flaws.

## Build the image

```shell
./build.sh
```
## Run experiments
There are few strategies already in library, you can choose one strategy to run.

```shell
./test.sh exante 3600
```

And you can choose multi strategies test together.

```shell
./test.sh exante,confuse,random 3600
```
