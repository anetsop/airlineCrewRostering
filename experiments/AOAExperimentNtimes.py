import subprocess
import changeDir
import experimentArguments
import numpy as np

if __name__ == "__main__":

    # run the application for multiple seeds and a single
    # parameter configuration using AOA algorithm
    changeDir.changeDirectory()
    times = int(experimentArguments.parseArguments("AOA experiment for airline crew rostering").times)
    C1 = 1.0
    C2 = 2.0
    C3 = 2.0
    C4 = 0.5
    C1str = str(C1)
    C2str = str(C2)
    C3str = str(C3)
    C4str = str(C4)
    for time in range(times):
        seed = experimentArguments.generateSeed()        
        outputFile  = "Output_AOA_" + time + ".xlsx"
        subprocess.run(["go", "run", "airlineCrewRostering.go",
                        "AOA", "-f", "Pairings.csv", "--startDate", "2011-11-1", "--endDate", "2012-3-4", 
                        "--seed", seed, "--C1", C1str, "--C2", C2str, 
                        "--C3", C3str, "--C4", C4str, "--generations", "200", "-p", "40", "--results", outputFile])