import sys
import numpy

GRID_SIZE = int(sys.argv[1])
NUM_AGENTS = int(sys.argv[2])
NUM_OBSTACLES = int(sys.argv[3])

f = open("tests/agents_test_{}_{}_{}.txt".format(GRID_SIZE,
         NUM_AGENTS, NUM_OBSTACLES), "w")

for _ in range(NUM_AGENTS):
    coordinates = numpy.random.randint(0, GRID_SIZE - 1, 4)
    line = "{} {} {} {}\n".format(
        coordinates[0], coordinates[1], coordinates[2], coordinates[3])
    f.write(line)

f = open("tests/obstacles_test_{}_{}_{}.txt".format(GRID_SIZE,
         NUM_AGENTS, NUM_OBSTACLES), "w")

for _ in range(NUM_OBSTACLES):
    coordinates = numpy.random.randint(0, GRID_SIZE - 1, 2)
    line = "{} {}\n".format(
        coordinates[0], coordinates[1])
    f.write(line)
