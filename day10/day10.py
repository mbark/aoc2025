# /// script
# dependencies = [
#   "numpy",
#   "scipy",
# ]
# ///

import json
import sys
import numpy as np
from itertools import product

def solve_system(equations, variables, domains, objective="minimize_sum"):
    """
    Solve linear system handling all cases.
    """
    n_vars = len(variables)
    n_eqs = len(equations)

    # Build coefficient matrix A and constants vector b
    A = np.zeros((n_eqs, n_vars))
    b = np.zeros(n_eqs)

    var_to_idx = {var: i for i, var in enumerate(variables)}

    for i, eq in enumerate(equations):
        b[i] = eq['value']
        for var, coeff in eq['variables'].items():
            A[i, var_to_idx[var]] = coeff

    print(f"Original system: {n_vars} vars, {n_eqs} eqs", file=sys.stderr)

    # Remove duplicate equations
    unique_rows = []
    unique_b = []
    seen = set()

    for i in range(len(A)):
        row_tuple = tuple(A[i])
        if row_tuple not in seen:
            seen.add(row_tuple)
            unique_rows.append(A[i])
            unique_b.append(b[i])
        else:
            for j in range(len(unique_rows)):
                if np.allclose(A[i], unique_rows[j]):
                    if not np.isclose(b[i], unique_b[j]):
                        print(f"Inconsistent duplicate equations", file=sys.stderr)
                        return None

    A = np.array(unique_rows)
    b = np.array(unique_b)
    n_eqs = len(A)

    rank = np.linalg.matrix_rank(A)
    dof = n_vars - rank

    print(f"After dedup: {n_eqs} eqs, rank={rank}, DOF={dof}", file=sys.stderr)

    # Check consistency
    Ab = np.column_stack([A, b])
    rank_Ab = np.linalg.matrix_rank(Ab)

    if rank_Ab > rank:
        print(f"Inconsistent: rank(A)={rank}, rank([A|b])={rank_Ab}", file=sys.stderr)
        return None

    if dof == 0:
        return solve_exact(A, b, variables, domains)
    else:
        return solve_underdetermined(A, b, variables, domains, objective, rank)

def solve_exact(A, b, variables, domains):
    """Solve exactly/over-determined system."""
    try:
        solution, residuals, rank, s = np.linalg.lstsq(A, b, rcond=None)

        if len(residuals) > 0 and residuals[0] > 1e-6:
            print(f"Large residuals: {residuals[0]}", file=sys.stderr)
            return None

        int_solution = np.rint(solution).astype(int)

        result_vec = A @ int_solution

        if not np.allclose(result_vec, b, atol=1e-6):
            print("Integer rounding broke solution", file=sys.stderr)
            return None

        result = {}
        for i, var in enumerate(variables):
            val = int(int_solution[i])
            min_val, max_val = domains.get(var, (0, 1000))
            if not (min_val <= val <= max_val):
                print(f"{var}={val} outside [{min_val}, {max_val}]", file=sys.stderr)
                return None
            result[var] = val

        return result
    except Exception as e:
        print(f"Error in solve_exact: {e}", file=sys.stderr)
        return None

def find_independent_columns(A, rank):
    """
    Find which columns of A are linearly independent.
    Returns basis_indices and free_indices.
    """
    n_vars = A.shape[1]

    # Use SVD to find independent columns
    U, s, Vt = np.linalg.svd(A, full_matrices=False)

    # Find columns with largest singular values contribution
    # Simpler approach: just try columns in order and check linear independence
    basis_indices = []
    remaining = list(range(n_vars))

    for i in range(n_vars):
        if len(basis_indices) >= rank:
            break

        # Try adding column i
        test_cols = basis_indices + [i]
        A_test = A[:, test_cols]

        if np.linalg.matrix_rank(A_test) == len(test_cols):
            # Column i is independent
            basis_indices.append(i)
            remaining.remove(i)

    free_indices = remaining

    return basis_indices, free_indices

def solve_underdetermined(A, b, variables, domains, objective, rank):
    """Solve underdetermined system."""
    n_vars = len(variables)
    n_eqs = len(A)
    dof = n_vars - rank

    print(f"Solving underdetermined system with DOF={dof}, rank={rank}", file=sys.stderr)

    # Find independent columns
    basis_indices, free_indices = find_independent_columns(A, rank)

    print(f"Basis indices: {basis_indices}", file=sys.stderr)
    print(f"Free indices: {free_indices}", file=sys.stderr)
    print(f"Basis vars: {[variables[i] for i in basis_indices]}", file=sys.stderr)
    print(f"Free vars: {[variables[i] for i in free_indices]}", file=sys.stderr)

    # Build search ranges for free variables
    free_ranges = []
    for idx in free_indices:
        var = variables[idx]
        min_val, max_val = domains.get(var, (0, 1000))
        free_ranges.append(range(min_val, max_val + 1))
        print(f"  {var}: [{min_val}, {max_val}] ({len(range(min_val, max_val + 1))} values)", file=sys.stderr)

    total = 1
    for r in free_ranges:
        total *= len(r)

    print(f"Total combinations to check: {total}", file=sys.stderr)

    if total > 5_000_000:
        print("Large search space, using sampling", file=sys.stderr)
        return solve_with_sampling(A, b, variables, domains, objective,
                                   basis_indices, free_indices, free_ranges)

    best_solution = None
    best_sum = float('inf') if objective == "minimize_sum" else float('-inf')

    A_basis = A[:, basis_indices]
    A_free = A[:, free_indices] if len(free_indices) > 0 else np.zeros((n_eqs, 0))

    print(f"A_basis shape: {A_basis.shape}", file=sys.stderr)
    print(f"A_free shape: {A_free.shape}", file=sys.stderr)

    checked = 0
    valid_found = 0

    for free_vals in product(*free_ranges):
        x = np.zeros(n_vars)

        # Set free variables
        for i, idx in enumerate(free_indices):
            x[idx] = free_vals[i]

        # Compute adjusted RHS
        if len(free_indices) > 0:
            b_adj = b - A_free @ np.array(free_vals)
        else:
            b_adj = b

        try:
            # Solve for basis variables
            x_basis = np.linalg.lstsq(A_basis, b_adj, rcond=None)[0]
            x_basis_int = np.rint(x_basis).astype(int)

            # Check if solution is close to integer
            if not np.allclose(x_basis, x_basis_int, atol=0.001):
                continue

            # Check domain constraints for basis variables
            valid = True
            for i, idx in enumerate(basis_indices):
                var = variables[idx]
                min_val, max_val = domains.get(var, (0, 1000))
                val = int(x_basis_int[i])
                if not (min_val <= val <= max_val):
                    valid = False
                    break
                x[idx] = val

            if not valid:
                continue

            # Verify the full solution satisfies all equations
            result_vec = A @ x
            if np.allclose(result_vec, b, atol=1e-6):
                valid_found += 1
                total_sum = int(np.sum(x))

                if valid_found <= 5:  # Print first few solutions
                    sol_dict = {variables[i]: int(x[i]) for i in range(n_vars)}
                    print(f"  Valid #{valid_found}: {sol_dict}, sum={total_sum}", file=sys.stderr)

                if (objective == "minimize_sum" and total_sum < best_sum) or \
                        (objective == "maximize_sum" and total_sum > best_sum):
                    best_sum = total_sum
                    best_solution = {var: int(x[i]) for i, var in enumerate(variables)}

        except Exception as e:
            if checked < 5:  # Only print first few errors
                print(f"Error at iteration {checked}: {e}", file=sys.stderr)
            continue

        checked += 1
        if checked % 1000 == 0:
            print(f"  Progress: {checked}/{total}, found {valid_found} valid", file=sys.stderr)

    print(f"Final: checked {checked}, found {valid_found} valid solutions", file=sys.stderr)
    if best_solution:
        print(f"Best solution: sum={best_sum}", file=sys.stderr)
    return best_solution

def solve_with_sampling(A, b, variables, domains, objective, basis_indices,
                        free_indices, free_ranges, n=100000):
    """Random sampling for large spaces."""
    import random

    n_vars = len(variables)
    n_eqs = len(A)
    A_basis = A[:, basis_indices]
    A_free = A[:, free_indices] if len(free_indices) > 0 else np.zeros((n_eqs, 0))

    best_solution = None
    best_sum = float('inf') if objective == "minimize_sum" else float('-inf')

    valid_found = 0

    for iteration in range(n):
        free_vals = [random.choice(list(r)) for r in free_ranges]

        x = np.zeros(n_vars)
        for i, idx in enumerate(free_indices):
            x[idx] = free_vals[i]

        if len(free_indices) > 0:
            b_adj = b - A_free @ np.array(free_vals)
        else:
            b_adj = b

        try:
            x_basis = np.linalg.lstsq(A_basis, b_adj, rcond=None)[0]
            x_basis_int = np.rint(x_basis).astype(int)

            if not np.allclose(x_basis, x_basis_int, atol=0.001):
                continue

            valid = True
            for i, idx in enumerate(basis_indices):
                var = variables[idx]
                min_val, max_val = domains.get(var, (0, 1000))
                if not (min_val <= x_basis_int[i] <= max_val):
                    valid = False
                    break
                x[idx] = x_basis_int[i]

            if valid and np.allclose(A @ x, b, atol=1e-6):
                valid_found += 1
                total_sum = int(np.sum(x))

                if (objective == "minimize_sum" and total_sum < best_sum) or \
                        (objective == "maximize_sum" and total_sum > best_sum):
                    best_sum = total_sum
                    best_solution = {var: int(x[i]) for i, var in enumerate(variables)}
        except:
            continue

        if iteration % 10000 == 0 and iteration > 0:
            print(f"  Sampling: {iteration}/{n}, found {valid_found} valid", file=sys.stderr)

    print(f"Sampling complete: found {valid_found} valid solutions", file=sys.stderr)
    return best_solution

def main():
    data = json.load(sys.stdin)

    equations = data['equations']
    variables = data['variables']
    domains = data.get('domains', {})
    objective = data.get('objective', 'minimize_sum')

    solution = solve_system(equations, variables, domains, objective)

    if solution:
        result = {
            'status': 'success',
            'solution': solution,
            'sum': sum(solution.values())
        }
    else:
        result = {
            'status': 'no_solution',
            'solution': None,
            'sum': None
        }

    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()