#!/usr/bin/python3
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
  
    print("Processing...")
  
    sampleStudents = getSampleStudents("./sample_data/course_selection_data.csv", True)
    samplemockCourses = getSampleCourses("./sample_data/course_selection_data.csv", True)
    timetable = {}
    timetable["Version"] = 3
    timetable["timetable"] = generateScheduleV3(sampleStudents, samplemockCourses, 40, "./output/students.json", "./output/conflicts.json")
  
    f = open('./output/conflicts.json')
    conflicts = json.load(f)
 
    print(f"{len(conflicts)}/{len(sampleStudents)} have errors")
    errors = round(len(conflicts) / len(sampleStudents) * 100, 2)
    success = round(100 - errors, 2)

    print(f"{errors}% errors")
    print(f"{success}% success")
 
    f.close()

  else:
    print("Invalid argument")
    exit()

  with open("./output/timetable.json", "w") as outfile:
    json.dump(timetable, outfile, indent=2)

  print("Done")
