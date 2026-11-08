import re

COVERAGE_EXCLUDE = (
    r"main.go|api.go|test|docs|_func_demo|config|"
    r"internal/service/mocks/|"
    r"internal/service/services_mocks/"
)

with open("coverage.tmp", "r") as infile, open("coverage.out", "w") as outfile:
    for line in infile:
        if not re.search(COVERAGE_EXCLUDE, line):
            outfile.write(line)