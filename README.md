# monzo-account-status

## What

A Monzo account status outputter. For an authenticated user (API playground token, no OAuth refresh flow) with a current account, this outputs:

- The preferred name of the user who owns the account.
- The account ID.
- The balance of the current account.
- The total balance of the current account + savings + other pots.

## Usage

```shell
export MONZO_API_TOKEN=...
go run main.go
```

## Why

I wanted to brush up on Go for an actual useful home project.

## Future enhancements

I'm always wary of saying future things in READMEs, in case I never get time.  But this isn't going to get popular, so:

- Finding important transactions (eg rent, bills, groceries) and outputting a CSV for the shared house expenses spreadsheet?
- Command line switches for each of the current features?
- When I learn a bit more Go, refactoring the code to be better.
- ...
