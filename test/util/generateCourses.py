#!/usr/bin/env python3
import json
import csv

realCourses = {}

# Get all courses from real sample data
def getSampleCourses(data_dir, log=False) -> dict:
  with open(data_dir, newline='') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
      exists = False
      for course in realCourses:
        exists = True if realCourses[course]["CrsNo"] == row["CrsNo"] else False
        if exists: break
      if not exists:
        realCourses[row["CrsNo"]] = {
          "CrsNo": row["CrsNo"],
          "Requests": 0,
          "Description": row["Description"],
          "Credits": 4,
          "students": []
        }

  if log:
    with open("./output/realCourses.json", "w") as outfile:
      json.dump(realCourses, outfile, indent=2)
        
  return realCourses


if __name__ == '__main__':
  courseSet: dict = getSampleCourses("../sample_data/course_selection_data.csv")

  with open("../output/realCourses.json", "w") as outfile:
    json.dump(courseSet, outfile, indent=2)
