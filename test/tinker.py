#!/usr/bin/python3
from prettytable import PrettyTable
from typing import Tuple
import json
import sys

# Import required utilities
from util.mockStudents import generateMockStudents, getSampleStudents
from util.generateCourses import getSampleCourses
from util.courses import mockCourses

# Import Algorithms
from scheduleGenerator.generator_v1 import generateScheduleV1
from scheduleGenerator.generator_v2 import generateScheduleV2
from scheduleGenerator.generator_v3 import generateScheduleV3

def errorOutput(students) -> Tuple[PrettyTable, dict, dict]:
  # Error Table calulation / output  
  f = open('./output/conflicts.json')
  conflicts = json.load(f)
  f.close()

  t = PrettyTable(['Type', 'Error %', 'Success %', 'Error Ratio'])
  
  errorsC = round(len(conflicts["Critical"]) / len(students) * 100, 2)
  successC = round(100 - errorsC, 2)
  errorsA = round(len(conflicts["Acceptable"]) / len(students) * 100, 2)
  successA = round(100 - errorsA, 2)
  
  t.add_row(['Critical', f"{errorsC} %", f"{successC} %", f"{len(conflicts['Critical'])}/{len(students)}"])
  t.add_row(['Acceptable', f"{errorsA} %", f"{successA} %", f"{len(conflicts['Acceptable'])}/{len(students)}"])
  
  return t, conflicts["Critical"], conflicts["Acceptable"]

if __name__ == '__main__':
  
  if len(sys.argv) == 1:
    print("Missing argument")
    exit()

  if sys.argv[1].upper() == 'V1':
    print("Processing...")

    mockStudents = generateMockStudents(400)
    timetable = {}
    timetable["Version"] = 1
    timetable["timetable"] = generateScheduleV1(mockStudents, mockCourses)
  
  elif sys.argv[1].upper() == 'V2':
    print("Processing...")
  
    mockStudents = generateMockStudents(400)
    timetable = {}
    timetable["Version"] = 2
    timetable["timetable"] = generateScheduleV2(mockStudents, mockCourses)
  

  elif sys.argv[1].upper() == 'V3':
  
    print("Processing...\n")
  
    sampleStudents = getSampleStudents("./sample_data/course_selection_data.csv", True)
    samplemockCourses = getSampleCourses("./sample_data/course_selection_data.csv", True)
    timetable = {}
    timetable["Version"] = 3
    timetable["timetable"] = generateScheduleV3(sampleStudents, samplemockCourses, 40, "./output/students.json", "./output/conflicts.json")

    errors, _, _ = errorOutput(sampleStudents)
    print(errors)

  elif sys.argv[1].upper() == "ERRORS":
    sampleStudents = getSampleStudents("./sample_data/course_selection_data.csv", True)
    errors, critical, acceptable = errorOutput(sampleStudents)
    print()
    print(errors)

    print(f"\n{len(critical)} critical errors")
    c1, c2, co = 0, 0, 0
    for error in critical:
      if error["Code"] == "C-MC": c1 +=1
      elif error["Code"] == "C-CSS": c2 += 1
      else: co += 1

    print(f"x{c1} C-MC Errors: Critical - Missing Classes")
    print(f"x{c2} C-CSS Errors: Critical - Couldn't Solve Schedule")
    print(f"x{co} Other/Undefined Critical Errors")

    print(f"\n{len(acceptable)} acceptable errors")
    a1, ao = 0, 0
    for error in acceptable:
      if error["Code"] == "A-MC": a1 +=1
      else: ao += 1

    print(f"x{a1} A-MC Errors: Acceptable - Missing 1-2 Classes")
    print(f"x{ao} Other/Undefined Acceptable Errors")

    exit()

  else:
    print("Invalid argument")
    exit()

  with open("./output/timetable.json", "w") as outfile:
    json.dump(timetable, outfile, indent=2)

  print("\nDone")
