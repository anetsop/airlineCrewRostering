import argparse
import numpy as np

def parseArguments(desc):
    # parse the arguments used by the experiments
    parser = argparse.ArgumentParser(description = desc)
    parser.add_argument('--seed', '-s', required=False, default="-1", \
        help='Seed for experiment', metavar='S')
    parser.add_argument('--times', '-t', required=False, default="1", \
        help='How many times to conduct the experiment', metavar='T')
    args = parser.parse_args()
    if args.seed == "-1":
        args.seed = generateSeed()
    return args

def generateSeed():
    # generate a large random number to be used as seed
    Seedpart1 = np.random.randint(100000, 200000)
    Seedpart2 = np.random.randint(100000, 999999)
    Seedpart3 = np.random.randint(1000000, 9999999)
    seed = str(Seedpart1) + str(Seedpart2) + str(Seedpart3)
    return seed