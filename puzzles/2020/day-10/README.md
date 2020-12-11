# Day 10 Notes

Day 10 was humbling. I danced near a correct solution for hours, and did
actually have a correct solution that was not efficient enough to run on the
non-sample input.

After struggling for a bit, I was told there was an iterative solution but it
didn't come to me, and I ultimately turned to Reddit, where I found one post
that had a very clean looking recursive approach. Cleaner than the one I had
implemented. This one focused on counting upward and checking whether a given
number was present in the list of adapters, rather than walking the list of
adapters and trying to count valid combinations by checking "connectivity" of
mutated lists of adapters.

I also experienced a number of issues trying to create slices from existing
slices in Go, where the results were not what I expected. Copying the slice
outright and modifying the copy worked, but seems clumsy. It makes sense if
I think about the way slices work in Go, how it's a header plus a pointer to
an array, so if you modify the array, you're modifying it in every slice that
points to it.

I'm putting everything is as it existed at the moment I got the right answer
to finish Part 2. This includes functions that represent aborted attempts. If
that doesn't appear to be the case, it means I've cleaned things up, but I
wanted to include the slightly misguided attempts I was trying--well after I
should have called it for the night--as historical context. If the source
doesn't look like that, then I've cleaned it up since then.

This all has me bewildered by the folks who do Advent of Code in so many
different languages, considering that I was just outsmarted by the language I
currently know best and am far most comfortable using. It's good to be
reminded to be humble, even if it's frustrating in the moment.
