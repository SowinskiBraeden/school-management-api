#!/usr/bin/python3
import sys
import json
import math
from prettytable import PrettyTable
from courses import courses, activeCourses
from mock import generateMockStudents

'''
  I will be using python to test and
  develop course selection / schedule
  generation

  this is honestly terrible code but 
  this is not an easy task for me to
  do

  Block 1-4 is first semester while
  block 5-8 is second semester
'''

'''
schedule example:
schedule: {
  "block1": "className",
  "block2": "className",
  "block3": "className",
  "block4": "className" 
  "block5": "className",
  "block6": "className",
  "block7": "className",
  "block8": "className"
}
'''

'''
running example:
running: {
  "block1": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block2": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block3": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block4": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block5": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block6": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block7": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}},
  "block8": {classCode:{"className":name,"students":[student Name]},classCode:{"className":name,"students":[student Name]}}
}
'''

# Error 1: No classes in schedule can fit this student
# Error 2: No more room in schedule for another class

# 400 by default
studentsNum = 400

err1, err2, = 0,0
minReq, classCap, blockClassLimit = 18, 30, 21
mockStudents = []
running = {
  "block1": {},
  "block2": {},
  "block3": {},
  "block4": {},
  "block5": {},
  "block6": {},
  "block7": {},
  "block8": {}
}


def generateScheduleV1():
  global err1, err2, err3
  # Collect data and calculate schedules
  for student in mockStudents:
    # Tally class request
    for request in student["requests"]:
      courses[request]["totalrequests"] += 1
      courses[request]["studentindexes"].append(mockStudents.index(student))
      # Add course to active list if enough requests
      if courses[request]["totalrequests"] > minReq and courses[request]["code"] not in activeCourses: activeCourses[courses[request]["code"]] = courses[request]

  for student in mockStudents:
    alternateOffset = len(student["requests"])-8
    alternateIndex = 8
    for i in range(len(student["requests"])-alternateOffset): # Subtract x classes as they are alternatives
      currentCourse = student["requests"][i]
      generate = True
      while generate:
        # If class is allowed to run
        if currentCourse in activeCourses:
          blockIndex = 1
          getFreeBlock = True
          while getFreeBlock:
            block = f"block{blockIndex}"
            if currentCourse in running[block]:
              if student["schedule"][block] == "": # Add student to class
                if len(running[block][currentCourse]["students"]) < classCap:
                  running[block][currentCourse]["students"].append(student["name"])
                  student["schedule"][block] = courses[currentCourse]["name"]
                  getFreeBlock = False
                else: # Find next available class or create new one
                  if blockIndex == 8:  # No available classes
                    if len(student["requests"]) == 8:
                      # This student never had any alternatives
                      # How to solve?

                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                    else:
                      if alternateIndex <= (len(student["requests"]) - 1):
                        currentCourse = student["requests"][alternateIndex]
                        alternateIndex += 1
                      else: # No more alterantive
                        # print(f"\nError 1: No more available classes for student {student['name']}")
                        err1 += 1
                        generate = False
                    getFreeBlock = False
                  else: blockIndex += 1
              else:
                if blockIndex == 8: # No available classes
                  if len(student["requests"]) == 8:
                    # This student never had any alternatives
                    # How to solve?

                    # print(f"\nError 1: No more available classes for student {student['name']}")
                    err1 += 1
                    generate = False
                  else:
                    if alternateIndex <= (len(student["requests"]) - 1):
                      currentCourse = student["requests"][alternateIndex]
                      alternateIndex += 1
                    else: # No more alternatives
                      # print(f"\nError 1: No more available classes for student {student['name']}")
                      err1 += 1
                      generate = False
                  getFreeBlock = False
                else: blockIndex += 1
            else:
              if blockIndex == 8:
                # Class does not exists
                # Create new class in first available slot 
                blockNum = 1
                while True:
                  newBlock = f"block{blockNum}"
                  if currentCourse not in running[newBlock] and len(running[newBlock]) < blockClassLimit:
                    if student["schedule"][newBlock] == "": # Add student to class
                      running[newBlock][currentCourse] = {
                        "name": courses[currentCourse]["name"],
                        "students": [student["name"]]
                      }
                      student["schedule"][newBlock] = courses[currentCourse]["name"]
                      break
                    else:
                      if blockNum == 8:
                        # All student classes have been filled
                        break
                      else: blockNum += 1
                  else:
                    if blockNum == 8:
                      # No room in school for more classes
                      # print(f"\nError 2: No more room in school for another class")
                      err2 += 1
                      break
                    else: blockNum += 1
                break
              else: blockIndex += 1
          break
        elif currentCourse not in activeCourses:
          if len(student["requests"]) == 8:
            # This student never had any alternatives
            # How to solve?

            # print(f"\nError 1: No more available classes for student {student['name']}")
            err1 += 1
            generate = False
          else:
            if alternateIndex <= (len(student["requests"]) - 1):
              currentCourse = student["requests"][alternateIndex]
              alternateIndex += 1
            else: # Out of alternatives
              # print(f"\nError 1: No more available classes for student {student['name']}")
              err1 += 1
              generate = False


def generateScheduleV2():
  # This is just for data testing/visualizations
  # calculate # of times to run class
  # t = PrettyTable(["Class Name", "Class Runcount"])
  for i in range(len(activeCourses)):
    classRunCount = math.floor(activeCourses[list(activeCourses)[i]]["totalrequests"] / classCap)
    # If there is minReq+ requests left, 1 more class could be run
    if (activeCourses[list(activeCourses)[i]]["totalrequests"] % classCap) > minReq: classRunCount += 1
    activeCourses[list(activeCourses)[i]]["classRunCount"] = classRunCount
    # t.add_row([activeCourses[list(activeCourses)[i]]["code"], classRunCount])
  # print(t)


if __name__ == '__main__':
  if len(sys.argv) == 1:
    print("Missing argument")
    exit()
  if sys.argv[1].upper() == 'V1':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV1()
  elif sys.argv[1].upper() == 'V2':
    if len(sys.argv) == 3:
      try:
        studentsNum = int(sys.argv[2])
      except:
        print("Error parsing number of students")
        exit()
    print("Processing...")
    mockStudents = generateMockStudents(studentsNum)
    generateScheduleV1()
  else: 
    print("Invalid argument")
    exit()

  print("\n")
  print(f"Error 1: x{err1}")
  print(f"Error 2: x{err2}")
  print("\n")

  ### Displays the number of students in each class
  # for block in running:
  #   print(f"\nBlock: {block}")
  #   print("=====================")
  #   for cl in running[block]:
  #     name = running[block][cl]["name"]
  #     students = len(running[block][cl]["students"])
  #     print(f"Class: {name} | Students: {students}")

  # Count errors in students schedules
  errors = 0
  for i in range(len(mockStudents)):
    count = 0
    for course in mockStudents[i]["schedule"]:
      if mockStudents[i]["schedule"][course]=="": count+=1
    if count > 0: errors += 1
  
  print(f"{errors}/{studentsNum} student(s) have a issue with their schedule")


  with open("schedule.json", "w") as outfile:
    json.dump(running, outfile, indent=2)

  with open("students.json", "w") as outfile:
    json.dump(mockStudents, outfile, indent=2)

  print("Done")
