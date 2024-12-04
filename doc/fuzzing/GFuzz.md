# GFuzz

[Paper](https://github.com/system-pclub/GFuzz/blob/main/gfuzz.pdf)

## Mutating Message Order

- Only select statements are modified
- A program state is defined as a list of tuples. Each tuple $(s_i, c_i, e_i)$ describes a select state in the code with
  - $s_i$: ID of the select
  - $c_i$: number of cases in the select
  - $e_i$: chosen case
- For the mutation, GFuzz goes through each tuple and changes $e_i$ to an random but valid case.
- The number of new mutations generated depends on runtime feedback when exercising the order (see Favoring Propitious Orders)

## Enforcing Message order

- The select is changed in such a way, that for a specified time period $T$, only the prefered case is waited on.
- It $T$ has passed, the operation returns to executing the original select operation.

## Favoring Propitious Orders

- GFuzz tracks the following two information
  - interleavings of channel operations
  - channel states
- Monitor operations for each individual channel
- Encode the order of two same-channel operations
- Uses global data structure to count how many pairs of channel operation have been executed (CountChOpPair)
- Additionally, the following values are counted/measured:
  - Number of distinct channels created (CreateCh)
  - Number of channels closed / remaining open (CloseCh / NotCloseCh)
  - Maximum fullness of each buffered channel
- After each run, it is checked, if the run was interesting. A run is interesting, if at least one of the following things occurred:
  - The run contains a new pair of channel operations (new meaning it has not been seen in any of the previous runs)
  - If an operation pair's execution counter changes significantly from previous order. Specifically, if the counter falls into a range $(2^{N-1}, 2^N]$ to which no previous counter belongs (what does that mean? what is N?)
  - If a new channel operation is triggered, such as creating, closing or not closing a channel for the first time
  - If a buffered channel gets a larger maximum fullness than in all previous executions (MaxChBufFull)
- If a run was interesting, it is added into a queue to be mutated.
- For each operation in the queue, a score is calculated
  - $score = \sum \log_2 CountChOpPair + 10 \cdot \#CreateCh +  10\cdot \#CloseCh + 10\cdot \sum MaxChBufFull$
- The score determines how often the run is mutated
  - $\# mutations = \left\lceil 5 \cdot \frac{NewScore}{MaxScore}\right\rceil$


