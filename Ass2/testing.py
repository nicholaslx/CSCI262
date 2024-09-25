import matplotlib.pyplot as plt

def calculate_distribution(num_sub_puzzles, states_per_puzzle, k):
    counts = [0] * k
    total_cases = states_per_puzzle ** num_sub_puzzles
    # Iterate through each possible configuration
    for config in range(total_cases):
        if num_sub_puzzles == 1:
            hash_value = config % k
        else:
            # Convert config to a unique set of sub-puzzle values
            values = []
            temp = config
            for _ in range(num_sub_puzzles):
                values.append(temp % states_per_puzzle)
                temp //= states_per_puzzle
            hash_value = sum(values) % k
        counts[hash_value] += 1
    return counts

def plot_distribution(counts, k, title):
    plt.bar(range(k), counts)
    plt.xlabel('Hash Value')
    plt.ylabel('Number of Cases')
    plt.title(title)
    plt.show()

def average_hashes(counts):
    total_cases = sum(counts)
    weighted_sum = sum(count * hash_value for hash_value, count in enumerate(counts))
    return weighted_sum / total_cases

# Puzzle A parameters
k_A = 5
states_per_puzzle_A = 2**3  # 3 hidden bits
counts_A = calculate_distribution(1, states_per_puzzle_A, k_A)
plot_distribution(counts_A, k_A, 'Distribution of Hash Values for Puzzle A')
avg_hashes_A = average_hashes(counts_A)

# Puzzle B parameters
k_B = 3
states_per_puzzle_B = 2**3  # 3 hidden bits
counts_B = calculate_distribution(4, states_per_puzzle_B, k_B)
plot_distribution(counts_B, k_B, 'Distribution of Hash Values for Puzzle B')
avg_hashes_B = average_hashes(counts_B)

print(f"Average number of hashes needed for Puzzle A: {avg_hashes_A:.2f}")
print(f"Average number of hashes needed for Puzzle B: {avg_hashes_B:.2f}")
