import subprocess
import changeDir
import os
import ctypes

if __name__ == "__main__":
    changeDir.changeDirectory()
    subprocess.run(["go", "get", "-u", "-v", "github.com/akamensky/argparse"])
    subprocess.run(["go", "get", "gonum.org/v1/plot/..."])
    subprocess.run(["go", "get", "github.com/xuri/excelize/v2"])
 
    try: 
        os.mkdir("output", 0o777) 
    except OSError as error: 
        print("")
        
    ctypes.windll.user32.MessageBoxW(0, "Installation successful!", "Setup", 0)
