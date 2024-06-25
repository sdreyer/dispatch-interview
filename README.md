# Auction Algorithm
## Goal
The goal of this project is to implement an algorithm where people
can bid against for sale items. Bid entries with a starting bid, max bid, and bid 
are provided and the algorithm will determine the winner. 

## Structure
The project is structured into different packages for organization and extensibility. 
Summaries of the packages are provided below. 

### auction
The auction package contains common types that are used throughout the project. 

### bid_manager
The bid manager contains the logic and the algorithm to determine the winner. 
The package includes a DefaultBidManager that implements a BidManager interface to provide
the ability to swap the algorithm out with a different algorithm if desired. 

### currency
I was unsure if the use of the golang.org/x/text/currency package was allowed as it is hosted
by golang but not a standard library as specified by the requirements. I instead built
a custom, slimmed down, package as an alternative

### id_generator
This package contains an ID Generator to atomically create new EventIDs that are associated
with each new bid entry in order to properly break tied bids. I've created an in-memory ID Generator for testing, but have it
implementing an IDGenerator, so that it can be replaced with a remote backend in the case
that this project were to be distributed. It also contains a test to validate proper ID 
generation with concurrent calls

### storage
This package contains a storage layer to handle saving and fetching bid entries that are 
added. I've created an in-memory bid store that implements a BidStorer interface that allows
the memory implementation to be replaced with a database implementation in a distributed 
scenario.

## Design Choices
### Algorithm
I've designed the algorithm so that it calculates bids in rounds. It checks to see
if any bids are able to be increased to surpass the current winning bid each round, then
increases them to the minimum allowed value above the current winning bid if allowed
by the person max bid and increment value. 

### Code
I've provided additional comments in the code for more specific design decisions. 

## Running The Project
The project was built with Go 1.22.4 and tests can be run with `go test ./...` from the 
root directory. The code can be implemented in a larger project by importing `bid_manager.NewDefaultBidManager`.

## Tests
Test coverage for each package is around 90%-100% coverage with the missing coverage
being some error conditions.

### Structure
For the interface tests, I have used a pattern that has been effective before, where
each interface test file has a struct that can be initialized and called in the 
respective implementations test file. 

This allows for easier organization and separate for unit and integration tests if an
external database were to be used in an implementation. It also allows better separation
for implementation specific tests to be added