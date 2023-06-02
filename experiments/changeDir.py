import os

def changeDirectory():
    # change working directory from experiments
    # to airline_crew_rostering
    directory = os.getcwd()
    par_dir = os.path.dirname(directory)
    directory = par_dir + "\\" + "airline_crew_rostering"
    os.chdir(directory)