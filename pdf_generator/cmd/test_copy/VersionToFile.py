#!/usr/bin/python3
# Python script to create an empty file
# with current date as name.

# importing datetime module
from datetime import datetime
import os
  
# datetime.datetime.now() to get 
# current date as filename.
# TimeNow = datetime.datetime.now()
  
FILESUBVERSION="subversion.txt"
FILEDATE="date.txt"
  
# create empty file
def create_file():
    fmt = "%Y-%m-%d %H:%M:%S.%f"
    str1 = datetime.utcnow().strftime(fmt)[:-3]

    # Function creates an empty file
    # %d - date, %B - month, %Y - Year
    with open(FILEDATE, "w") as file:
        file.write(str1)
        file.close()

def set_vers():
    filename=FILESUBVERSION
    build=0
    mode = 'r' if os.path.exists(filename) else 'w+'
    with open(filename, encoding="utf8", mode=mode) as file_in:
        _str_build = file_in.read()
        file_in.close()
        try:
            build = int(_str_build)
        except ValueError as err:
            print("Build.__setVers(): при конвертировании строки в число, err=", err)
        finally:
            pass
    build += 1
    str_build = str(build)
    while len(str_build) < 5:
        str_build = "0" + str_build
    print("Build.__set_vers(): new build=", str_build)
    with open(filename, "w", encoding="utf8") as file_in:
        file_in.write(str_build)
        file_in.close()

  
# Driver Code
create_file()
set_vers()