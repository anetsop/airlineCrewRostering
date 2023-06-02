import subprocess
import changeDir
import experimentArguments
import numpy as np

if __name__ == "__main__":

    # run the application for multiple seeds and a single
    # parameter value using multi-step CSO algorithm

    changeDir.changeDirectory()
    times = int(experimentArguments.parseArguments("CSO experiment for airline crew rostering").times)
    fl = 1.0
    flStr = str(fl)
    for time in range(times):
        seed = experimentArguments.generateSeed()        
        outputFile  = "Output_CSO_" + time + ".xlsx"
        subprocess.run(["go", "run", "airlineCrewRostering.go",
        "multiCSO", "-f", "Pairings.csv", "--startDate", "2011-11-1", "--endDate", "2012-3-4", 
        "--seed", seed, "--FL", flStr, "--generations", "200", "-p", "45", "--results", outputFile])