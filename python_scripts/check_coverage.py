
import argparse
import re
import subprocess
import sys
from pathlib import Path


COVERAGE_INPUT = Path("coverage.tmp")
COVERAGE_OUTPUT = Path("coverage.out")

COVERAGE_EXCLUDE = (
    r"main\.go|api\.go|redis_client_mock\.go|test|docs|_func_demo|config|"
    r"internal/service/mocks/|"
    r"internal/service/services_mocks/"
)


def filter_coverage():
    if not COVERAGE_INPUT.exists():
        raise FileNotFoundError("Coverage profile not found!")

    lines = COVERAGE_INPUT.read_text(encoding="utf-8").splitlines()

    if not lines:
        raise ValueError("Coverage profile is empty!")

    # Preserve coverage mode header
    filtered_lines = [lines[0]]

    # Exclude unwanted files using regex
    for line in lines[1:]:
        if not re.search(COVERAGE_EXCLUDE, line):
            filtered_lines.append(line)

    COVERAGE_OUTPUT.write_text(
        "\n".join(filtered_lines) + "\n",
        encoding="utf-8"
    )


def get_coverage():
    result = subprocess.run(
        [
            "go",
            "tool",
            "cover",
            f"-func={COVERAGE_OUTPUT}"
        ],
        capture_output=True,
        text=True,
        check=True
    )

    match = re.search(
        r"total:\s+\(statements\)\s+([\d.]+)%",
        result.stdout
    )

    if not match:
        raise ValueError("Cannot calculate total coverage!")

    return float(match.group(1))


def main():
    parser = argparse.ArgumentParser(
        description="Check Go test coverage threshold"
    )

    parser.add_argument(
        "--threshold",
        type=float,
        default=80,
        help="Minimum coverage percentage"
    )

    args = parser.parse_args()

    try:
        filter_coverage()

        total = get_coverage()

        print(f"\nTotal coverage: {total:.2f}%")
        print(f"Coverage threshold: {args.threshold:.2f}%")

        if total < args.threshold:
            print(
                f"\nERROR: Coverage ({total:.2f}%) "
                f"is below threshold ({args.threshold:.2f}%)"
            )
            sys.exit(1)

        print(
            f"\nSUCCESS: Coverage ({total:.2f}%) "
            f"meets threshold ({args.threshold:.2f}%)"
        )

    except (
        FileNotFoundError,
        ValueError,
        subprocess.CalledProcessError,
    ) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
