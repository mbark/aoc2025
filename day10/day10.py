# /// script
# dependencies = [
#   "z3-solver",
# ]
# ///

import json
import sys

from z3 import Int, Solver, sat, Optimize


def solve_system(equations, variables, domains, objective="minimize_sum"):
    """
    Solve linear system using Z3 SMT solver.
    """
    print(f"System: {len(variables)} vars, {len(equations)} eqs", file=sys.stderr)

    # Create Z3 integer variables
    z3_vars = {}
    for var in variables:
        z3_vars[var] = Int(var)

    # Use Optimizer to handle objective function
    solver = Optimize()

    # Add domain constraints
    for var in variables:
        min_val, max_val = domains.get(var, (0, 1000))
        solver.add(z3_vars[var] >= min_val)
        solver.add(z3_vars[var] <= max_val)
        print(f"  {var}: [{min_val}, {max_val}]", file=sys.stderr)

    # Add equation constraints
    for i, eq in enumerate(equations):
        expr = sum(coeff * z3_vars[var] for var, coeff in eq['variables'].items())
        solver.add(expr == eq['value'])
        print(f"  Equation {i + 1}: {' + '.join(f'{c}*{v}' for v, c in eq['variables'].items())} = {eq['value']}",
              file=sys.stderr)

    # Add objective function
    total_sum = sum(z3_vars[var] for var in variables)
    if objective == "minimize_sum":
        solver.minimize(total_sum)
    else:
        solver.maximize(total_sum)

    # Solve
    print("Solving...", file=sys.stderr)
    result = solver.check()

    if result == sat:
        model = solver.model()
        solution = {}
        for var in variables:
            solution[var] = model[z3_vars[var]].as_long()

        print(f"Solution found: sum={sum(solution.values())}", file=sys.stderr)
        return solution
    else:
        print("No solution found", file=sys.stderr)
        return None

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
