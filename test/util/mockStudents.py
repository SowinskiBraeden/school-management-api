#!/usr/bin/env python3
import random
import names
import json
import csv
import sys
from util.courses import mockCourses

flex= ["XAT--12A-S", "XAT--12B-S"]

mockStudents: list[dict] = []

# Generate n students for mock data
def generateMockStudents(n: int) -> list[dict]:
  for _ in range(n):
    newStudent = {
      "name": names.get_full_name(),
      "requests": [], # list of class codes
      "schedule": {
        "block1": "",
        "block2": "",
        "block3": "",
        "block4": "",
        "block5": "",
        "block6": "",
        "block7": "",
        "block8": ""
      }
    }
    # Get list of random class choices with no repeats
    # 8 primary choices, 2 secondary choices
    courseSelection = random.sample(range(0, len(mockCourses)), 10)
    for courseNum in courseSelection:
      newStudent["requests"].append(list(mockCourses)[courseNum])
    mockStudents.append(newStudent)

  return mockStudents


# sort real sample data into usable dictionary
def getSampleStudents(data_dir: str, log: bool = False) -> list[dict]:
  with open(data_dir, newline='') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
      exists = False
      for student in mockStudents:
        exists = True if student["Pupil #"] == row["Pupil #"] else False
        if exists: break
      alternate = True if row["Alternate?"] == 'TRUE' else False
      if exists:
        if len(mockStudents[student["studentIndex"]]["requests"]) >= 10 and not alternate and row["CrsNo"] not in flex: alternate = True
        mockStudents[student["studentIndex"]]["requests"].append({
          "CrsNo": row["CrsNo"],
          "Description": row["Description"],
          "alt": alternate
        })
        if row["CrsNo"] not in flex and not alternate and mockStudents[student["studentIndex"]]["expectedClasses"] < 10:
          mockStudents[student["studentIndex"]]["expectedClasses"] += 1
      else:
        newStudent = {
          "Pupil #": row["Pupil #"],
          "requests": [{
            "CrsNo": row["CrsNo"],
            "Description": row["Description"],
            "alt": alternate
          }],
          "schedule": {
            "block1": [],
            "block2": [],
            "block3": [],
            "block4": [],
            "block5": [],
            "block6": [],
            "block7": [],
            "block8": [],
            "block9": [],
            "block10": []
          },
          "expectedClasses": 1,
          "classes": 0,
          "remainingAlts": [],
          "studentIndex": len(mockStudents)
        }
        mockStudents.append(newStudent)

  if log:
    with open("./output/students.json", "w") as outfile:
      json.dump(mockStudents, outfile, indent=2)

  return mockStudents


if __name__ == '__main__':
  if len(sys.argv) == 1:
    print("Missing argument")
    exit()
  if sys.argv[1].lower() == 'sample':
    studentRequests: list[dict] = getSampleStudents("../sample_data/course_selection_data.csv")

    with open("../output/students.json", "w") as outfile:
      json.dump(studentRequests, outfile, indent=2)
  elif sys.argv[1].lower() == 'mock':
    studentRequests: list[dict] = generateMockStudents(400)

    with open("../output/students.json", "w") as outfile:
      json.dump(studentRequests, outfile, indent=2)
  else:
    print("Invalid argument")
    exit()

  print("done")
