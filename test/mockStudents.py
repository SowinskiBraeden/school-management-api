#!/usr/bin/python3
import random
import names
import json
import csv
from courses import mockCourses

mockStudents = []

# Generate n students for mock data
def generateMockStudents(n):
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
def getSampleStudents(log=False):
  with open("course_selection_data.csv", newline='') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
      exists = False
      for student in mockStudents:
        exists = True if student["Pupil #"] == row["Pupil #"] else False
        if exists: break
      alternate = True if row["Alternate?"] == 'TRUE' else False
      if exists:
        mockStudents[student["studentIndex"]]["requests"].append({
            "CrsNo": row["CrsNo"],
            "Description": row["Description"],
            "alt": alternate
          })
      else:
        newStudent = {
          "Pupil #": row["Pupil #"],
          "requests": [{
            "CrsNo": row["CrsNo"],
            "Description": row["Description"],
            "alt": alternate
          }],
          "schedule": {
            "block1": "",
            "block2": "",
            "block3": "",
            "block4": "",
            "block5": "",
            "block6": "",
            "block7": "",
            "block8": ""
          },
          "studentIndex": len(mockStudents)
        }
        mockStudents.append(newStudent)

  if log:
    with open("students.json", "w") as outfile:
      json.dump(mockStudents, outfile, indent=2)

  return mockStudents


if __name__ == '__main__':
  studentRequests = getSampleStudents()

  with open("students.json", "w") as outfile:
    json.dump(studentRequests, outfile, indent=2)
