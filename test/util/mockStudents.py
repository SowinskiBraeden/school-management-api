#!/usr/bin/env python3
import json
import csv

flex= ["XAT--12A-S", "XAT--12B-S"]

mockStudents: list[dict] = []

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
  studentRequests: list[dict] = getSampleStudents("../sample_data/course_selection_data.csv")

  with open("../output/students.json", "w") as outfile:
    json.dump(studentRequests, outfile, indent=2)

  print("done")
