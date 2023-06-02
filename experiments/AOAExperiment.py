import subprocess
import changeDir
import experimentArguments
import experimentRun
import numpy as np
from multiprocessing import Pool

if __name__ == "__main__":

    # run the application for a specific seed and multiple
    # parameter configurations using AOA algorithm
    changeDir.changeDirectory()
    seed = experimentArguments.parseArguments("AOA experiment for airline crew rostering").seed
    
    for C2 in [2,4,6]:
        for C1 in [1,2]:
            for C3 in [1,2]:
                for C4 in [0.5, 1]:
                    C1str = str(C1)
                    C2str = str(C2)
                    C3str = str(C3)
                    C4str = str(C4)
                    outputFile  = "Output_seed_" + seed + "_C1_" + C1str + "_C2_" + C2str + "_C3_" + C3str + "_C4_" + C4str + ".xlsx"                    
                    subprocess.run(["go", "run", "airlineCrewRostering.go",
                                    "AOA", "-f", "Pairings.csv", "--startDate", "2011-11-1", "--endDate", "2012-3-4", 
                                    "--seed", seed, "--C1", C1str, "--C2", C2str, 
                                    "--C3", C3str, "--C4", C4str, "--generations", "200", "-p", "40", "--results", outputFile])