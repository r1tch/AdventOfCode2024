from functools import lru_cache

# Given a full A and a number of iterations, run the program and return outputs
def run_program(A):
    outputs = []
    current_A = A
    while current_A != 0:
        B = (current_A & 7)        # Step 1
        B ^= 1                     # Step 2
        C = current_A >> B         # Step 3
        B ^= 5                     # Step 4
        current_A = current_A >> 3 # Step 5
        B ^= C                     # Step 6
        out = B & 7                # Step 7
        outputs.append(out)
    return outputs

# Suppose we have a target output sequence:
target_outputs = [2, 4, 1, 1, 7, 5, 1, 5, 0, 3, 4, 4, 5, 5, 3, 0]

# We will build A incrementally, iteration by iteration.
# At iteration i:
#  - Bits needed: from 3*i to 3*i+5 (6 bits total).
#  - We try all 64 possibilities for these 6 bits.
#  - Combine them with existing partial A (which may already have bits set).
#  - Check if the outputs up to iteration i match target_outputs[:i+1].
#
# We'll keep a set of possible As after each iteration that match the prefix.

def get_bits(x, start, length):
    mask = (1 << length) - 1
    return (x >> start) & mask

def set_bits(x, start, length, value):
    # Set 'length' bits starting at 'start' in x to 'value'
    mask = (1 << length) - 1
    # Clear old bits
    x &= ~(mask << start)
    # Set new bits
    x |= (value & mask) << start
    return x

def check_compatibility(partial_A, start_bit, bits_to_set):
    # Check if partial_A already has bits set in the specified region
    # that conflict with bits_to_set.
    length = 6
    existing = get_bits(partial_A, start_bit, length)
    # We need to ensure that for all positions where partial_A is already defined,
    # it matches bits_to_set. Actually, since we are always setting bits,
    # partial_A might have been partially built. If we trust that partial_A
    # is always fully under our control, no explicit check is needed.
    # But if we want to be extra safe, we can check if partial_A & mask == bits_to_set.
    # For now, assume partial_A is always constructed by us, so no conflict.
    return True

def run_iterations(A, count):
    # Run the program for 'count' iterations only and return the first 'count' outputs.
    outputs = []
    current_A = A
    for _ in range(count):
        if current_A == 0:
            # If A is zero before we reach 'count' iterations, outputs end.
            break
        B = (current_A & 7)
        B ^= 1
        C = current_A >> B
        B ^= 5
        current_A = current_A >> 3
        B ^= C
        out = B & 7
        outputs.append(out)
    return outputs

# We'll do a forward incremental search:
partial_solutions = {0}  # Start with A=0 as a base (no bits set yet)
max_iterations = len(target_outputs)

for i in range(max_iterations):
    print("ITERATION", i, "partial_solutions", len(partial_solutions))
    new_solutions = set()
    # We need to set bits [3*i .. 3*i+5]
    start_bit = 3*i
    length = 11

    for partial_A in partial_solutions:
        #print("trying to expand", partial_A)
        # Try all 64 possibilities for these 6 bits
        for candidate in range(2**length):
            #print("trying", oct(candidate))
            # Check compatibility (if needed)
            if not check_compatibility(partial_A, start_bit, candidate):
                continue
            # Create a new A with these bits set
            A_candidate = set_bits(partial_A, start_bit, length, candidate)

            # Now check if the first i+1 outputs match:
            # We only need to run the program for i+1 iterations to verify the prefix.
            prefix_outputs = run_iterations(A_candidate, i+1)
            if prefix_outputs == target_outputs[:i+1]:
                #print('candidate:', oct(A_candidate), prefix_outputs)
                new_solutions.add(A_candidate)
            #else:
                #print(' NOT GOOD', oct(A_candidate), prefix_outputs)

    # Move to the next iteration with the filtered set of solutions
    partial_solutions = new_solutions

    # If no solutions remain at any point, we can stop early
    if not partial_solutions:
        break

# After all iterations, partial_solutions contains all As that produce the full sequence
print("Found solutions:", len(partial_solutions))
for sol in sorted(partial_solutions)[:10]:
    print("A candidate:", sol)
