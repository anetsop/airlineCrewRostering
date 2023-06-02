import subprocess
import changeDir
import experimentArguments
import numpy as np

if __name__ == "__main__":

    # run the application for a specific seed and multiple
    # parameter values using multi-step CSO algorithm

    changeDir.changeDirectory()
    seed = experimentArguments.parseArguments("CSO experiment for airline crew rostering").seed
    for fl in np.arange(0, 2.1, 0.1):
        fl = round(fl, 1)
        flStr = str(fl)
        outputFile  = "Output_FL_" + flStr + "-_seed_" + seed + "-.xlsx"
        subprocess.run(["go", "run", "airlineCrewRostering.go",
        "multiCSO", "-f", "Pairings.csv", "--startDate", "2011-11-1", "--endDate", "2012-3-4", 
        "--seed", seed, "--FL", flStr, "--generations", "200", "-p", "45", "--results", outputFile])